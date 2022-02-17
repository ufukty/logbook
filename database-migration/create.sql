/*
    FIXME: CHECK USER INPUT BEFORE START TO EXECUTION
    FIXME: INSERT A COLUMN TO TASKS TABLE FOR TASK ACTIVATION STATUS (BOOLEAN) 
*/

DROP DATABASE IF EXISTS logbook_dev;

CREATE DATABASE logbook_dev;

\c logbook_dev;

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
    "parent_id"         UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    "content"           TEXT NOT NULL,
    "degree"            INT NOT NULL DEFAULT 1,
    "depth"             INT NOT NULL DEFAULT 1,
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "completed_at"      DATE, -- FIXME: MAKE IT TIMESTAMP AND UPDATE document_overview FUNCTION
    "ready_to_pick_up"  BOOLEAN NOT NULL DEFAULT TRUE
);

ALTER TABLE "DOCUMENT" ADD 
    "active_task"       UUID REFERENCES "TASK" ("task_id");

-- CREATE VIEW tasks_linearized AS SELECT * FROM tasks;

-- DROP FUNCTION IF EXISTS create_document_with_task_groups;
CREATE FUNCTION create_document() RETURNS "DOCUMENT" AS $$
    DECLARE
        document "DOCUMENT";
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
                    AND "document_id" = v_document_id
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
                    AND "document_id" = v_document_id
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

-- RECURSIVE HELPER FUNCTION FOR:
--     * update_task_depth
-- IT UPDATES THE DEPTH OF TASK AND RECURSES INTO ITS CHILDREN
CREATE FUNCTION update_task_depth_helper(v_task_id UUID, v_new_depth INT) RETURNS UUID[] AS $$
    DECLARE
        v_child_id UUID;
        v_updated_task_id UUID;
        v_updated_task_list UUID[];
    BEGIN
        UPDATE "TASK"
        SET "depth" = v_new_depth
        WHERE "task_id" = v_task_id
        RETURNING "task_id" INTO v_updated_task_id;

        IF v_updated_task_id IS NULL THEN
            RETURN v_updated_task_list;
        END IF;

        FOR v_child_id IN (SELECT * FROM "TASK" WHERE "parent_id" = v_task_id) 
        LOOP
            v_updated_task_list = array_cat(
                v_updated_task_list,
                update_task_depth_helper(v_child_id, v_new_depth + 1)
            );
        END LOOP;

        RETURN v_updated_task_list;
    END
$$ LANGUAGE 'plpgsql';

-- RECURSIVE HELPER FUNCTION FOR:
--     * create_task
--     * reattach_task
CREATE FUNCTION update_task_depth(v_task_id UUID) RETURNS UUID[] AS $$
    DECLARE
        v_task_depth INT;
        v_task "TASK";
        -- v_root_uuid UUID = '00000000-0000-0000-0000-000000000000';
    BEGIN
        SELECT *
        INTO v_task
        FROM "TASK"
        WHERE "task_id" = v_task_id;

        IF v_task."parent_id" = '00000000-0000-0000-0000-000000000000' THEN
            v_task_depth = 1;
        ELSE
            SELECT "depth" + 1
            INTO v_task_depth
            FROM "TASK"
            WHERE "task_id" = v_task_id;
        END IF;

        RETURN update_task_depth_helper(v_task_id, v_task_depth);
    END
$$ LANGUAGE 'plpgsql';


CREATE FUNCTION array_to_sorted_row_list(v_id_array UUID[], v_task_id UUID) RETURNS SETOF "TASK" AS $$
    BEGIN
        RETURN QUERY (
            SELECT * FROM "TASK" t WHERE t."task_id" = ANY(v_id_array) 
            ORDER BY CASE WHEN t."task_id" = v_task_id THEN 0 ELSE 1 END
        );
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
        -- RAISE NOTICE 'v_content = %, v_parent_id = %', v_content, v_parent_id;

        -- DEGREE ALWAYS 1, WHEN THE TASK IS NEWLY ADDED
        v_degree = 1;

        IF v_parent_id IS NULL THEN
            RAISE EXCEPTION 'Parent ID can not be NULL';
            RETURN;
        END IF;

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

        RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task."task_id");
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
        --     * DEPTH (RECURSIVELY)
        --     * READINESS STATUS
        v_updated_task_list = array_cat(
            v_updated_task_list,
            update_task_degree(v_new_parent_id, v_task_old."degree")
        );
        PERFORM update_task_readineess(v_new_parent_id);
        v_updated_task_list = array_cat(
            v_updated_task_list,
            update_task_depth(v_task_id)
        );

        -- RETURN UPDATED TASKS AS ARRAY FOR UPDATING FRONTEND 
        RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task_id);
    END 
$$ LANGUAGE 'plpgsql';

CREATE FUNCTION mark_a_task_done(v_task_id UUID) RETURNS SETOF "TASK" AS $$
    DECLARE
        v_task "TASK";
        v_updated_task_list UUID[];
    BEGIN
        -- RAISE NOTICE 'mark_a_task_done, v_task_id = %', v_task_id;

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

        RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task_id);
    END
$$ LANGUAGE 'plpgsql';

CREATE PROCEDURE load_test_dataset() AS $$
    DECLARE
        v_document_id UUID DEFAULT '61bbc44a-c61c-4d49-8804-486181081fa7';
        v_task_1 UUID;
        v_task_2 UUID;
        v_task_3 UUID;
        v_task_4 UUID;
        v_task_5 UUID;
        v_task_6 UUID;
        v_task_7 UUID;
        v_task_8 UUID;
        v_task_9 UUID;
        v_task_10 UUID;
        v_task_11 UUID;
        v_task_12 UUID;
        v_task_13 UUID;
        v_task_14 UUID;
        v_task_15 UUID;
        v_task_16 UUID;
        v_task_17 UUID;
        v_task_18 UUID;
        v_task_19 UUID;
        v_task_20 UUID;
        v_task_21 UUID;
        v_task_22 UUID;
        v_task_23 UUID;
        v_task_24 UUID;
        v_task_25 UUID;
        v_task_26 UUID;
        v_task_27 UUID;
        v_task_28 UUID;
        v_task_29 UUID;
        v_task_30 UUID;
        v_task_31 UUID;
        v_task_32 UUID;
        v_task_33 UUID;
        v_task_34 UUID;
        v_task_35 UUID;
    BEGIN
        -- SELECT "document_id" INTO v_document_id FROM create_document();
        INSERT INTO "DOCUMENT"("document_id") VALUES (v_document_id);

        RAISE NOTICE 'document_id: %', v_document_id; 

        -- FIRST ROOT TASK

        SELECT "task_id" INTO v_task_1 FROM create_task(v_document_id => v_document_id, v_content => 'deploy redis cluster on multi DC');
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '20 DAYS' WHERE "task_id" = v_task_1;
        
        SELECT "task_id" INTO v_task_2 FROM create_task(v_document_id, 'deploy redis cluster on 1 DC', v_task_1);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '20 DAYS' WHERE "task_id" = v_task_2;

        SELECT "task_id" INTO v_task_3 FROM create_task(v_document_id, 'Revoke passwordless sudo rights after provision at cluster', v_task_1);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '20 DAYS' WHERE "task_id" = v_task_18;
        
        SELECT "task_id" INTO v_task_4 FROM create_task(v_document_id, 'iptables for redis', v_task_3);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '20 DAYS' WHERE "task_id" = v_task_29;
        
        SELECT "task_id" INTO v_task_5 FROM create_task(v_document_id, 'terraform for redis', v_task_4);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '20 DAYS' WHERE "task_id" = v_task_3;
        
        SELECT "task_id" INTO v_task_6 FROM create_task(v_document_id, 'Update redis/tf file according to prod.tfvars file', v_task_4);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '20 DAYS' WHERE "task_id" = v_task_6;

        SELECT "task_id" INTO v_task_7 FROM create_task(v_document_id, 'Remove: seperator from ovpn-auth', v_task_2);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '19 DAYS' WHERE "task_id" = v_task_7;
        
        SELECT "task_id" INTO v_task_8 FROM create_task(v_document_id, 'Write tests for ovpn-auth', v_task_3);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '18 DAYS' WHERE "task_id" = v_task_8;
        
        SELECT "task_id" INTO v_task_9 FROM create_task(v_document_id, 'Decrease timing gap of ovpn-auth under 1ms', v_task_3);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '18 DAYS' WHERE "task_id" = v_task_9;
        
        SELECT "task_id" INTO v_task_10 FROM create_task(v_document_id, 'Prepare releases for ovpn-auth', v_task_4);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '18 DAYS' WHERE "task_id" = v_task_10;
        
        SELECT "task_id" INTO v_task_11 FROM create_task(v_document_id, 'Provision golden-image for gitlab-runner', v_task_10);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '18 DAYS' WHERE "task_id" = v_task_11;

        -- SECOND ROOT TASK

        SELECT "task_id" INTO v_task_12 FROM create_task(v_document_id => v_document_id, v_content => 'gitlab-runner --(vpn)--> DNS ----> gitlab');
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '17 DAYS' WHERE "task_id" = v_task_12;
       
        SELECT "task_id" INTO v_task_13 FROM create_task(v_document_id, 'Firewall & unbound rules update from prov script (VPN)', v_task_12);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '17 DAYS' WHERE "task_id" = v_task_13;
        
        SELECT "task_id" INTO v_task_14 FROM create_task(v_document_id, 'Script pic_gitlab_runner_post_creation', v_task_13);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '17 DAYS' WHERE "task_id" = v_task_14;
        
        SELECT "task_id" INTO v_task_15 FROM create_task(v_document_id, 'Execute 1 CI/CD pipeline on gitlab-runner', v_task_14);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '17 DAYS' WHERE "task_id" = v_task_15;
        
        SELECT "task_id" INTO v_task_16 FROM create_task(v_document_id, 'gitlab-runner provisioner with resolv.conf/docker/runner-register', v_task_12);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '17 DAYS' WHERE "task_id" = v_task_16;
        
        SELECT "task_id" INTO v_task_17 FROM create_task(v_document_id, 'prepare gitlab-ci for ovpn-auth repo', v_task_13);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '17 DAYS' WHERE "task_id" = v_task_17;
        
        SELECT "task_id" INTO v_task_18 FROM create_task(v_document_id, 'PAM for SSH', v_task_14);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '17 DAYS' WHERE "task_id" = v_task_18;
        
        SELECT "task_id" INTO v_task_19 FROM create_task(v_document_id, 'ACL - Redis', v_task_13);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '16 DAYS' WHERE "task_id" = v_task_19;
        
        SELECT "task_id" INTO v_task_20 FROM create_task(v_document_id, 'Redis security', v_task_13);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '16 DAYS' WHERE "task_id" = v_task_20;
        
        SELECT "task_id" INTO v_task_21 FROM create_task(v_document_id, 'TOTP for SSH', v_task_14);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '16 DAYS' WHERE "task_id" = v_task_21;
        
        SELECT "task_id" INTO v_task_22 FROM create_task(v_document_id, 'API gateway without redis', v_task_15);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '15 DAYS' WHERE "task_id" = v_task_22;
        
        SELECT "task_id" INTO v_task_23 FROM create_task(v_document_id, 'Golden image interitance re-organize', v_task_16);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '15 DAYS' WHERE "task_id" = v_task_23;
        
        SELECT "task_id" INTO v_task_24 FROM create_task(v_document_id, 'Postgres', v_task_12);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '14 DAYS' WHERE "task_id" = v_task_24;
        
        SELECT "task_id" INTO v_task_25 FROM create_task(v_document_id, 'Auth service', v_task_13);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '14 DAYS' WHERE "task_id" = v_task_25;
        
        SELECT "task_id" INTO v_task_26 FROM create_task(v_document_id, 'MQ', v_task_15);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '14 DAYS' WHERE "task_id" = v_task_26;
        
        SELECT "task_id" INTO v_task_27 FROM create_task(v_document_id, 'Federated learning', v_task_16);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '13 DAYS' WHERE "task_id" = v_task_27;
        
        SELECT "task_id" INTO v_task_28 FROM create_task(v_document_id, 'Bluetooth transmission test', v_task_12);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '13 DAYS' WHERE "task_id" = v_task_28;
        
        SELECT "task_id" INTO v_task_29 FROM create_task(v_document_id, 'Intrusion detection system (centralised) (OSSEC', v_task_12);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '12 DAYS' WHERE "task_id" = v_task_29;
        
        SELECT "task_id" INTO v_task_30 FROM create_task(v_document_id, 'Envoy - HAProxy - NGiNX', v_task_12);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '12 DAYS' WHERE "task_id" = v_task_30;
        
        SELECT "task_id" INTO v_task_31 FROM create_task(v_document_id, 'web-front/Privacy against [friend/pubic/company/attackers]', v_task_13);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '12 DAYS' WHERE "task_id" = v_task_31;
        
        SELECT "task_id" INTO v_task_32 FROM create_task(v_document_id, 'Redis/cluster script test for multi datacenter', v_task_13);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '12 DAYS' WHERE "task_id" = v_task_32;
        
        SELECT "task_id" INTO v_task_33 FROM create_task(v_document_id, 'gitlab-runner firewall rules: close public internet', v_task_14);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '12 DAYS' WHERE "task_id" = v_task_33;
        
        SELECT "task_id" INTO v_task_34 FROM create_task(v_document_id, 'static-challange for ovpn-auth', v_task_20);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '11 DAYS' WHERE "task_id" = v_task_34;
        
        SELECT "task_id" INTO v_task_35 FROM create_task(v_document_id, 'Golden image for vpn server', v_task_21);
        UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '11 DAYS' WHERE "task_id" = v_task_35;

        -- COMPLETE SOME TASKS

        PERFORM mark_a_task_done(v_task_18);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '4 DAYS' WHERE "task_id" = v_task_18;
        PERFORM mark_a_task_done(v_task_19);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_19;
        PERFORM mark_a_task_done(v_task_20);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_20;
        PERFORM mark_a_task_done(v_task_21);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_21;
        PERFORM mark_a_task_done(v_task_22);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_22;
        PERFORM mark_a_task_done(v_task_23);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_23;
        PERFORM mark_a_task_done(v_task_24);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_24;
        PERFORM mark_a_task_done(v_task_25);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_25;
        PERFORM mark_a_task_done(v_task_26);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '4 DAYS' WHERE "task_id" = v_task_26;
        PERFORM mark_a_task_done(v_task_27);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_27;
        PERFORM mark_a_task_done(v_task_28);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_28;
        PERFORM mark_a_task_done(v_task_29);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_29;
        PERFORM mark_a_task_done(v_task_30);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_30;
        PERFORM mark_a_task_done(v_task_31);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_31;
        PERFORM mark_a_task_done(v_task_32);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_32;
        PERFORM mark_a_task_done(v_task_33);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_33;
        PERFORM mark_a_task_done(v_task_34);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_34;
        PERFORM mark_a_task_done(v_task_35);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_35;
        PERFORM mark_a_task_done(v_task_6);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_6;
        PERFORM mark_a_task_done(v_task_7);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_7;
        PERFORM mark_a_task_done(v_task_8);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_8;
        PERFORM mark_a_task_done(v_task_9);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_9;
        PERFORM mark_a_task_done(v_task_10);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_10;
        PERFORM mark_a_task_done(v_task_11);
        UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_11;

        PERFORM reattach_task(v_task_13, v_task_1);
        PERFORM reattach_task(v_task_27, v_task_15);
        PERFORM reattach_task(v_task_31, v_task_5);
    END
$$ LANGUAGE 'plpgsql';

CALL load_test_dataset();