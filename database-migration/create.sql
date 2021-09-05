DROP DATABASE IF EXISTS testdatabase;

CREATE DATABASE testdatabase;

\c testdatabase;

CREATE TABLE "DOCUMENT" (
    "document_id"       UUID UNIQUE DEFAULT gen_random_UUID(), 
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "ACCESS"(
    "document_id"       UUID NOT NULL REFERENCES "DOCUMENT" ("document_id"),
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "user-agent"        VARCHAR(256),
    "ip-address"        INET NOT NULL
);

CREATE TABLE "TASK" (
    "task_id"           UUID UNIQUE DEFAULT gen_random_UUID(),
    "document_id"       UUID NOT NULL REFERENCES "DOCUMENT" ("document_id"),
    "parent_id"         UUID DEFAULT '00000000-0000-0000-0000-000000000000',
    "content"           TEXT NOT NULL,
    "degree"            INT NOT NULL DEFAULT 1,
    "depth"             INT NOT NULL DEFAULT 1,
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "completed_at"      DATE,
    "ready_to_pick_up"  BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE "DOCUMENT" ADD 
    "active_task"       UUID REFERENCES "TASK" ("task_id");

-- CREATE VIEW tasks_linearized AS SELECT * FROM tasks;

-- DROP FUNCTION IF EXISTS create_document_with_task_groups;
CREATE FUNCTION create_document() RETURNS "DOCUMENT" AS $$
    DECLARE
        document "DOCUMENT"%ROWTYPE;
    BEGIN
        INSERT INTO "DOCUMENT" DEFAULT VALUES RETURNING * INTO document;
        RETURN document;
    END
$$ LANGUAGE 'plpgsql';

-- TO GET THE LIST OF TASKS THAT ARE:
--   * COMPLETED WITHIN LAST 10 DAYS (<200 ITEM),
--   * ACTIVE,
--   * READY-TO-PICK-UP
--   * DRAWER (TO-DO)
CREATE FUNCTION document_overview(v_document_id UUID) RETURNS SETOF "TASK" AS $$
    BEGIN
        RETURN QUERY (
            (
                -- ARCHIVED TASKS (LATEST 200 FROM LAST 10 DAYS)
                SELECT *
                FROM "TASK"
                WHERE "completed_at" IN (
                    SELECT DISTINCT ON ("completed_at") 
                        "completed_at"
                    FROM "TASK" 
                    WHERE "document_id" = v_document_id
                    ORDER BY "completed_at" DESC
                    LIMIT 10
                )
                LIMIT 200
            ) 
            UNION ALL
            (
                -- ACTIVE TASK
                SELECT *
                FROM "TASK"
                WHERE "task_id" = (
                    SELECT "active_task"
                    FROM "DOCUMENT"
                    WHERE "document_id" = v_document_id
                )
            )
            UNION ALL
            (
                -- READY-TO-PICK-UPS (100)
                SELECT *
                FROM "TASK"
                WHERE 
                    "completed_at" IS NULL
                    AND "ready_to_pick_up" = TRUE
                ORDER BY "created_at" ASC
                LIMIT 100
            )
            UNION ALL
            (
                -- DRAWER (20)
                SELECT *
                FROM "TASK"
                WHERE 
                    "completed_at" IS NULL
                    AND "ready_to_pick_up" = FALSE
                ORDER BY "degree" DESC
                LIMIT 20
            )
        );
    END
$$ LANGUAGE 'plpgsql';

-- RECURSIVE HELPER FUNCTION FOR:
--    * create_task
--    * reattach_task
CREATE FUNCTION update_task_degree(v_task_id UUID, v_increment INT) RETURNS UUID[] AS $$
    DECLARE
        v_total_degrees_of_siblings INT;
        v_task "TASK";
        v_updated_task_list UUID[];
    BEGIN
        -- RAISE NOTICE 'update_task_degree, v_task_id = %', v_task_id;

        UPDATE "TASK"
        SET "degree" = "degree" + v_increment
        WHERE "TASK"."task_id" = v_task_id;

        SELECT *
        INTO v_task
        FROM "TASK"
        WHERE "task_id" = v_task_id;

        -- RAISE NOTICE 'v_task = %', v_task;
        -- RAISE NOTICE 'v_task."parent_id" = %', v_task."parent_id";

        -- ADD ITSELF TO v_updated_task_list BEFORE THE TASKS THAT WILL RETURNED BY PARENT
        v_updated_task_list = array_append(v_updated_task_list, v_task."task_id");

        IF v_task."parent_id" != '00000000-0000-0000-0000-000000000000' THEN
            -- RAISE NOTICE 'recursing into parent';
            v_updated_task_list = array_cat(
                v_updated_task_list, 
                update_task_degree(v_task."parent_id", v_increment)
            );
        END IF;
        
        -- RAISE NOTICE 'no more parent to recurse further, returning to caller now';

        RETURN v_updated_task_list;
    END
$$ LANGUAGE 'plpgsql';

CREATE FUNCTION update_task_readineess(v_task_id UUID) RETURNS UUID AS $$
    DECLARE
        v_undone_children "TASK";
        v_readiness BOOLEAN;
        v_task "TASK";
    BEGIN
        -- TODO: IT CAN RETURN UUID[] IF PARENTS OF PARENTS TAKEN INTO CALCULATION

        -- RAISE NOTICE 'update_task_readineess is running for %', v_task_id;

        SELECT *
        INTO v_undone_children
        FROM "TASK"
        WHERE "parent_id" = v_task_id
            AND "completed_at" IS NULL;

        IF v_undone_children IS NULL THEN
            v_readiness = TRUE;
        ELSE
            v_readiness = FALSE;
        END IF;

        -- RAISE NOTICE 'v_readiness = %', v_readiness;

        UPDATE "TASK"
        SET "ready_to_pick_up" = v_readiness
        WHERE "task_id" = v_task_id 
            AND "ready_to_pick_up" <> v_readiness
        RETURNING * INTO v_task;

        RETURN v_task."task_id";
    END
$$ LANGUAGE 'plpgsql';

-- Add new task to document with:
--   * updating the degree of parent task (and theirs, recursively).
--   * minding the depth of parent task.
-- Returns the list of updated tasks.
CREATE FUNCTION create_task(
    v_document_id UUID, 
    v_content VARCHAR,
    v_parent_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
) RETURNS SETOF "TASK" AS $$
    DECLARE
        v_depth INT;
        v_degree INT;
        v_task "TASK"%ROWTYPE;
        v_updated_task_list UUID[];
    BEGIN
        -- RAISE NOTICE 'v_parent_id = %', v_parent_id;

        -- DEGREE ALWAYS 1, WHEN THE TASK IS NEWLY ADDED
        v_degree = 1;

        -- DECIDE DEPTH
        IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN
            SELECT "TASK"."depth"+1 
            INTO v_depth
            FROM "TASK" 
            WHERE "task_id" = v_parent_id;
        ELSE
            v_depth = 1;
        END IF;

        -- WRITE NEW TASK INTO DB
        INSERT INTO "TASK"("document_id", "parent_id", "content", "degree", "depth")
        VALUES (v_document_id, v_parent_id, v_content, v_degree, v_depth)
        RETURNING * INTO v_task;

        -- INITIALIZE THE RETURN LIST
        v_updated_task_list = array_append(v_updated_task_list, v_task."task_id");

        -- UPDATE PARENT'S READY-TO-PICK-UP STATUS TO FALSE
        IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN
            UPDATE "TASK"
            SET "ready_to_pick_up" = FALSE
            WHERE task_id = v_parent_id;
        END IF;

        -- UPDATE PARENTS' DEGREES
        IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN
            v_updated_task_list = array_cat(
                v_updated_task_list, 
                update_task_degree(v_parent_id, 1)
            );
        END IF;

        RETURN QUERY SELECT * FROM "TASK" WHERE "task_id" = ANY(v_updated_task_list);
    END
$$ LANGUAGE 'plpgsql';

CREATE FUNCTION reattach_task(v_task_id UUID, v_new_parent_id UUID) RETURNS SETOF "TASK" AS $$ -- RETURN "TASK"
    DECLARE
        v_updated_task_list UUID[];
        v_task_old "TASK";
    BEGIN
        -- FIXME: CHECK CIRCULAR DEPENDENCY
    
        -- TEMPORARILY STORE THE TASK WITH ITS CURRENT CONDITION
        SELECT * 
        INTO v_task_old 
        FROM "TASK" 
        WHERE "task_id" = v_task_id;

        -- DON'T CONTINUE IF THE DESIRED NEW PARENT IS THE TASK'S ITSELF
        IF v_task_id = v_new_parent_id THEN
            RETURN;
        END IF;

        -- UPDATE TASK:
        --     * PARENT
        --     * DEPTH
        UPDATE "TASK"
        SET "parent_id" = v_new_parent_id
        WHERE "task_id" = v_task_id;

        -- INITILIAZE RETURN LIST WITH THE MODIFIED TASK AS ITS FIRST ITEM
        v_updated_task_list = array_append(v_updated_task_list, v_task_old."task_id");

        -- UPDATE OLD PARENT:
        --     * DEGREE (RECURSIVELY)
        --     * READINESS STATUS
        IF v_task_old."parent_id" != '00000000-0000-0000-0000-000000000000' THEN
            v_updated_task_list = array_cat(
                v_updated_task_list,
                update_task_degree(v_task_old."parent_id", -1 * v_task_old."degree")
            );
            PERFORM update_task_readineess(v_task_old."parent_id");
        END IF;

        -- UPDATE NEW PARENT:
        --     * DEGREE (RECURSIVELY)
        --     * READINESS STATUS
        v_updated_task_list = array_cat(
            v_updated_task_list,
            update_task_degree(v_new_parent_id, v_task_old."degree")
        );
        PERFORM update_task_readineess(v_new_parent_id);

        -- RETURN UPDATED TASKS AS ARRAY FOR UPDATING FRONTEND 
        RETURN QUERY SELECT * FROM "TASK" WHERE "task_id" = ANY(v_updated_task_list);
    END 
$$ LANGUAGE 'plpgsql';

CREATE FUNCTION mark_a_task_done(v_task_id UUID) RETURNS SETOF "TASK" AS $$
    DECLARE
        v_task "TASK";
        v_updated_task_list UUID[];
    BEGIN
        RAISE NOTICE 'mark_a_task_done, v_task_id = %', v_task_id;

        -- UPDATE TASK'S ITSELF
        UPDATE "TASK"
        SET completed_at = CURRENT_TIMESTAMP
        WHERE "task_id" = v_task_id 
            AND "completed_at" IS NULL
        RETURNING * INTO v_task;
        
        -- UPDATE TASK'
        v_updated_task_list = array_append(v_updated_task_list, v_task."task_id");

        -- UPDATE PARENT READINESS
        IF v_task IS NOT NULL THEN
            v_updated_task_list = array_append(
                v_updated_task_list,
                update_task_readineess(v_task."parent_id") 
            );
        END IF;

        RETURN QUERY SELECT * FROM "TASK" WHERE "task_id" = ANY(v_updated_task_list);
    END
$$ LANGUAGE 'plpgsql';

CREATE PROCEDURE load_test_dataset() AS $$
    DECLARE
        document_id UUID;
        updated_tasks "TASK"[];
        v_tasks_1 "TASK"[];
        v_tasks_2 "TASK"[];
        v_tasks_3 "TASK"[];
        v_tasks_4 "TASK"[];
        v_tasks_5 "TASK"[];
        v_tasks_6 "TASK"[];
        v_tasks_7 "TASK"[];
        v_tasks_8 "TASK"[];
        v_tasks_9 "TASK"[];
        v_tasks_10 "TASK"[];
        v_tasks_11 "TASK"[];
        v_tasks_12 "TASK"[];
        v_tasks_13 "TASK"[];
        v_tasks_14 "TASK"[];
        v_tasks_15 "TASK"[];
        v_tasks_16 "TASK"[];
        v_tasks_17 "TASK"[];
        v_tasks_18 "TASK"[];
        v_tasks_19 "TASK"[];
        v_tasks_20 "TASK"[];
        v_tasks_21 "TASK"[];
        v_tasks_22 "TASK"[];
        v_tasks_23 "TASK"[];
        v_tasks_24 "TASK"[];
        v_tasks_25 "TASK"[];
        v_tasks_26 "TASK"[];
        v_tasks_27 "TASK"[];
        v_tasks_28 "TASK"[];
        v_tasks_29 "TASK"[];
        v_tasks_30 "TASK"[];
        v_tasks_31 "TASK"[];
        v_tasks_32 "TASK"[];
        v_tasks_33 "TASK"[];
        v_tasks_34 "TASK"[];
        v_tasks_35 "TASK"[];
        v_tasks_0 "TASK"[];
        v_task_35_mark_done "TASK"[];
    BEGIN
        INSERT INTO "DOCUMENT"("document_id") VALUES ('fe71c1e5-e6c8-587e-9647-e9ae9819eb8a') RETURNING "DOCUMENT"."document_id" INTO document_id;

        v_tasks_1 = create_task(v_document_id => document_id, v_content => 'deploy redis cluster on multi DC');
        v_tasks_2 = create_task(document_id, 'deploy redis cluster on 1 DC', v_tasks_1[1].task_id);
        v_tasks_3 = create_task(document_id, 'Revoke passwordless sudo rights after provision at cluster', v_tasks_1[1].task_id);
        v_tasks_4 = create_task(document_id, 'iptables for redis', v_tasks_2[1].task_id);
        v_tasks_5 = create_task(document_id, 'terraform for redis', v_tasks_3[1].task_id);
        v_tasks_6 = create_task(document_id, 'Update redis/tf file according to prod.tfvars file', v_tasks_4[1].task_id);
        v_tasks_7 = create_task(document_id, 'Remove: seperator from ovpn-auth', v_tasks_2[1].task_id);
        v_tasks_8 = create_task(document_id, 'Write tests for ovpn-auth', v_tasks_3[1].task_id);
        v_tasks_9 = create_task(document_id, 'Decrease timing gap of ovpn-auth under 1ms', v_tasks_3[1].task_id);
        v_tasks_10 = create_task(document_id, 'Prepare releases for ovpn-auth', v_tasks_4[1].task_id);
        v_tasks_11 = create_task(document_id, 'Provision golden-image for gitlab-runner', v_tasks_10[1].task_id);
        
        v_tasks_12 = create_task(v_document_id => document_id, v_content => 'gitlab-runner --(vpn)--> DNS ----> gitlab');
        v_tasks_13 = create_task(document_id, 'Firewall & unbound rules update from prov script (VPN)', v_tasks_12[1].task_id);
        v_tasks_14 = create_task(document_id, 'Script pic_gitlab_runner_post_creation', v_tasks_13[1].task_id);
        v_tasks_15 = create_task(document_id, 'Execute 1 CI/CD pipeline on gitlab-runner', v_tasks_14[1].task_id);
        v_tasks_16 = create_task(document_id, 'gitlab-runner provisioner with resolv.conf/docker/runner-register', v_tasks_12[1].task_id);
        v_tasks_17 = create_task(document_id, 'prepare gitlab-ci for ovpn-auth repo', v_tasks_13[1].task_id);
        v_tasks_18 = create_task(document_id, 'PAM for SSH', v_tasks_14[1].task_id);
        v_tasks_19 = create_task(document_id, 'ACL - Redis', v_tasks_13[1].task_id);
        v_tasks_20 = create_task(document_id, 'Redis security', v_tasks_13[1].task_id);
        v_tasks_21 = create_task(document_id, 'TOTP for SSH', v_tasks_14[1].task_id);
        v_tasks_22 = create_task(document_id, 'API gateway without redis', v_tasks_15[1].task_id);
        v_tasks_23 = create_task(document_id, 'Golden image interitance re-organize', v_tasks_16[1].task_id);
        v_tasks_24 = create_task(document_id, 'Postgres', v_tasks_12[1].task_id);
        v_tasks_25 = create_task(document_id, 'Auth service', v_tasks_13[1].task_id);
        v_tasks_26 = create_task(document_id, 'MQ', v_tasks_15[1].task_id);
        v_tasks_27 = create_task(document_id, 'Federated learning', v_tasks_16[1].task_id);
        v_tasks_28 = create_task(document_id, 'Bluetooth transmission test', v_tasks_12[1].task_id);
        v_tasks_29 = create_task(document_id, 'Intrusion detection system (centralised) (OSSEC', v_tasks_12[1].task_id);
        v_tasks_30 = create_task(document_id, 'Envoy - HAProxy - NGiNX', v_tasks_12[1].task_id);
        v_tasks_31 = create_task(document_id, 'web-front/Privacy against [friend/pubic/company/attackers]', v_tasks_13[1].task_id);
        v_tasks_32 = create_task(document_id, 'Redis/cluster script test for multi datacenter', v_tasks_13[1].task_id);
        v_tasks_33 = create_task(document_id, 'gitlab-runner firewall rules: close public internet', v_tasks_14[1].task_id);
        v_tasks_34 = create_task(document_id, 'static-challange for ovpn-auth', v_tasks_20[1].task_id);
        v_tasks_35 = create_task(document_id, 'Golden image for vpn server', v_tasks_21[1].task_id);

        UPDATE "TASK" SET "created_at" = '2021-01-11T18:19:27+03:00' WHERE "task_id" = v_tasks_27[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-01-23T21:37:55+03:00' WHERE "task_id" = v_tasks_4[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-02-14T01:03:33+03:00' WHERE "task_id" = v_tasks_1[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-02-18T02:23:12+03:00' WHERE "task_id" = v_tasks_35[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-02-23T00:42:48+03:00' WHERE "task_id" = v_tasks_19[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-02-26T20:43:38+03:00' WHERE "task_id" = v_tasks_3[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-02-28T12:58:22+03:00' WHERE "task_id" = v_tasks_22[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-03-21T04:48:09+03:00' WHERE "task_id" = v_tasks_14[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-03-28T05:00:14+03:00' WHERE "task_id" = v_tasks_28[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-04-22T01:22:57+03:00' WHERE "task_id" = v_tasks_31[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-04-26T18:27:59+03:00' WHERE "task_id" = v_tasks_8[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-04-29T04:37:40+03:00' WHERE "task_id" = v_tasks_2[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-04-30T20:48:02+03:00' WHERE "task_id" = v_tasks_16[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-05-02T05:46:29+03:00' WHERE "task_id" = v_tasks_33[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-05-09T05:25:48+03:00' WHERE "task_id" = v_tasks_12[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-05-14T10:54:34+03:00' WHERE "task_id" = v_tasks_24[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-05-24T04:01:24+03:00' WHERE "task_id" = v_tasks_21[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-07-03T05:47:47+03:00' WHERE "task_id" = v_tasks_17[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-07-08T03:07:40+03:00' WHERE "task_id" = v_tasks_18[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-07-25T22:54:19+03:00' WHERE "task_id" = v_tasks_13[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-08-22T11:06:13+03:00' WHERE "task_id" = v_tasks_20[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-08-24T20:57:14+03:00' WHERE "task_id" = v_tasks_29[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-08-28T14:20:12+03:00' WHERE "task_id" = v_tasks_5[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-09-03T09:48:06+03:00' WHERE "task_id" = v_tasks_6[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-09-04T04:33:43+03:00' WHERE "task_id" = v_tasks_23[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-09-15T01:14:20+03:00' WHERE "task_id" = v_tasks_10[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-09-16T04:13:49+03:00' WHERE "task_id" = v_tasks_26[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-09-23T21:12:23+03:00' WHERE "task_id" = v_tasks_9[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-10-04T08:37:10+03:00' WHERE "task_id" = v_tasks_11[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-10-10T21:55:15+03:00' WHERE "task_id" = v_tasks_30[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-10-15T02:51:16+03:00' WHERE "task_id" = v_tasks_32[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-10-15T11:37:09+03:00' WHERE "task_id" = v_tasks_25[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-11-20T09:08:29+03:00' WHERE "task_id" = v_tasks_15[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-12-20T09:51:34+03:00' WHERE "task_id" = v_tasks_7[1]."task_id";
        UPDATE "TASK" SET "created_at" = '2021-12-31T10:52:07+03:00' WHERE "task_id" = v_tasks_34[1]."task_id";

        PERFORM mark_a_task_done(v_tasks_18[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-02T05:46:29+03:00' WHERE "task_id" = v_tasks_18[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_19[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-03T05:47:47+03:00' WHERE "task_id" = v_tasks_19[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_20[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-03T09:48:06+03:00' WHERE "task_id" = v_tasks_20[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_21[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-08T03:07:40+03:00' WHERE "task_id" = v_tasks_21[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_22[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-09T05:25:48+03:00' WHERE "task_id" = v_tasks_22[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_23[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-11T18:19:27+03:00' WHERE "task_id" = v_tasks_23[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_24[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-14T01:03:33+03:00' WHERE "task_id" = v_tasks_24[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_25[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-14T10:54:34+03:00' WHERE "task_id" = v_tasks_25[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_26[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-18T02:23:12+03:00' WHERE "task_id" = v_tasks_26[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_27[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-21T04:48:09+03:00' WHERE "task_id" = v_tasks_27[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_28[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-22T01:22:57+03:00' WHERE "task_id" = v_tasks_28[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_29[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-22T11:06:13+03:00' WHERE "task_id" = v_tasks_29[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_30[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-23T00:42:48+03:00' WHERE "task_id" = v_tasks_30[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_31[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-23T21:37:55+03:00' WHERE "task_id" = v_tasks_31[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_32[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-24T04:01:24+03:00' WHERE "task_id" = v_tasks_32[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_33[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-24T20:57:14+03:00' WHERE "task_id" = v_tasks_33[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_34[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-25T22:54:19+03:00' WHERE "task_id" = v_tasks_34[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_35[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-26T18:27:59+03:00' WHERE "task_id" = v_tasks_35[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_6[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-26T20:43:38+03:00' WHERE "task_id" = v_tasks_6[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_7[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-28T05:00:14+03:00' WHERE "task_id" = v_tasks_7[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_8[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-28T12:58:22+03:00' WHERE "task_id" = v_tasks_8[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_9[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-28T14:20:12+03:00' WHERE "task_id" = v_tasks_9[1]."task_id";
        PERFORM mark_a_task_done(v_tasks_10[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-29T04:37:40+03:00' WHERE "task_id" = v_tasks_10[1]."task_id";
        v_task_35_mark_done = mark_a_task_done(v_tasks_11[1]."task_id");
        UPDATE "TASK" SET "completed_at" = '2021-10-30T20:48:02+03:00' WHERE "task_id" = v_tasks_11[1]."task_id";

        RAISE NOTICE 'v_task_35_mark_done = %', v_task_35_mark_done;

        PERFORM reattach_task(v_tasks_13[1]."task_id", v_tasks_1[1]."task_id");
        PERFORM reattach_task(v_tasks_27[1]."task_id", v_tasks_15[1]."task_id");
        v_tasks_0 = reattach_task(v_tasks_31[1]."task_id", v_tasks_5[1]."task_id");

        RAISE NOTICE 'v_tasks_0 = %', v_tasks_0;
    END
$$ LANGUAGE 'plpgsql';

CALL load_test_dataset();