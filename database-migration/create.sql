DROP DATABASE IF EXISTS testdatabase;

CREATE DATABASE testdatabase;

\c testdatabase;

CREATE TABLE "DOCUMENT" (
    "document_id"       UUID UNIQUE DEFAULT gen_random_UUID(), 
    "created_at"        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "total_task_groups" INTEGER DEFAULT 0
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
    "completion_at"     DATE
);

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

-- TO GET THE LIST OF TASKS THAT ARE COMPLETED WITHIN LAST 10 DAYS
CREATE FUNCTION document_overview(UUID) RETURNS SETOF "TASK" AS $$
    BEGIN
        RETURN QUERY (
            (
                SELECT *
                FROM "TASK"
                WHERE "completion_at" IN (
                    SELECT DISTINCT ON ("completion_at") 
                        "completion_at"
                    FROM "TASK" 
                    WHERE "document_id" = $1
                    ORDER BY "completion_at" 
                    LIMIT 10
                )
                LIMIT 500
            ) UNION (
                SELECT *
                FROM "TASK"
                WHERE "completion_at" IS NULL
                ORDER BY "created_at" ASC
                LIMIT 50
            )
        );
    END
$$ LANGUAGE 'plpgsql';

-- RECURSIVE HELPER FUNCTION FOR create_task
CREATE PROCEDURE update_parent_task(v_task_id UUID) AS $$
    DECLARE
        v_total_degrees_of_siblings INT;
        v_task "TASK"%ROWTYPE;
    BEGIN
        RAISE NOTICE 'update_parent_task, v_task_id = %', v_task_id;

        SELECT sum("degree")
        INTO v_total_degrees_of_siblings
        FROM "TASK"
        WHERE "parent_id" = v_task_id;

        RAISE NOTICE 'v_total_degrees_of_siblings = %', v_total_degrees_of_siblings;

        UPDATE "TASK"
        SET degree = v_total_degrees_of_siblings + 1
        WHERE "TASK"."task_id" = v_task_id;

        SELECT *
        INTO v_task
        FROM "TASK"
        WHERE "task_id" = v_task_id;

        RAISE NOTICE 'v_task."parent_id" = %', v_task."parent_id";

        IF v_task."parent_id" != '00000000-0000-0000-0000-000000000000' THEN
            RAISE NOTICE 'recursing to v_task."parent_id" = %', v_task."parent_id";
            CALL update_parent_task(v_task."parent_id");
        END IF;
        
        RAISE NOTICE 'no more parent to recurse further, returning to caller now';
    END
$$ LANGUAGE 'plpgsql';

-- Add new task to document with:
--   * updating the degree of parent task (and theirs, recursively).
--   * minding the depth of parent task.
CREATE FUNCTION create_task(
    v_document_id UUID, 
    v_content VARCHAR,
    v_parent_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
) RETURNS "TASK" AS $$
    DECLARE
        v_depth INT;
        v_degree INT;
        v_task "TASK"%ROWTYPE;
    BEGIN
        v_degree = 1;

        IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN
            SELECT "TASK"."depth"+1 
            INTO v_depth
            FROM "TASK" 
            WHERE "task_id" = v_parent_id;
        ELSE
            v_depth = 1;
        END IF;

        INSERT INTO "TASK"("document_id", "parent_id", "content", "degree", "depth")
        VALUES (v_document_id, v_parent_id, v_content, v_degree, v_depth)
        RETURNING * INTO v_task;

        IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN
            CALL update_parent_task(v_parent_id);
        END IF;

        RETURN v_task;
    END
$$ LANGUAGE 'plpgsql';

CREATE PROCEDURE load_test_dataset() AS $$
    DECLARE
        document_id UUID;
        v_task_1 "TASK";
        v_task_2 "TASK";
        v_task_3 "TASK";
        v_task_4 "TASK";
        v_task_5 "TASK";
        v_task_6 "TASK";
        v_task_7 "TASK";
        v_task_8 "TASK";
        v_task_9 "TASK";
        v_task_10 "TASK";
        v_task_11 "TASK";
        v_task_12 "TASK";
        v_task_13 "TASK";
        v_task_14 "TASK";
        v_task_15 "TASK";
        v_task_16 "TASK";
        v_task_17 "TASK";
        v_task_18 "TASK";
        v_task_19 "TASK";
        v_task_20 "TASK";
        v_task_21 "TASK";
        v_task_22 "TASK";
        v_task_23 "TASK";
        v_task_24 "TASK";
        v_task_25 "TASK";
        v_task_26 "TASK";
        v_task_27 "TASK";
        v_task_28 "TASK";
        v_task_29 "TASK";
        v_task_30 "TASK";
        v_task_31 "TASK";
        v_task_32 "TASK";
        v_task_33 "TASK";
        v_task_34 "TASK";
        v_task_35 "TASK";
    BEGIN
        INSERT INTO "DOCUMENT"("document_id") VALUES ('fe71c1e5-e6c8-587e-9647-e9ae9819eb8a') RETURNING "DOCUMENT"."document_id" INTO document_id;

        v_task_1 = create_task(v_document_id => document_id, v_content => 'deploy redis cluster on multi DC');
        v_task_2 = create_task(document_id, 'deploy redis cluster on 1 DC', v_task_1.task_id);
        v_task_3 = create_task(document_id, 'Revoke passwordless sudo rights after provision at cluster', v_task_1.task_id);
        v_task_4 = create_task(document_id, 'iptables for redis', v_task_2.task_id);
        v_task_5 = create_task(document_id, 'terraform for redis', v_task_3.task_id);
        v_task_6 = create_task(document_id, 'Update redis/tf file according to prod.tfvars file', v_task_4.task_id);
        v_task_7 = create_task(document_id, 'Remove: seperator from ovpn-auth', v_task_2.task_id);
        v_task_8 = create_task(document_id, 'Write tests for ovpn-auth', v_task_3.task_id);
        v_task_9 = create_task(document_id, 'Decrease timing gap of ovpn-auth under 1ms', v_task_3.task_id);
        v_task_10 = create_task(document_id, 'Prepare releases for ovpn-auth', v_task_4.task_id);
        v_task_11 = create_task(document_id, 'Provision golden-image for gitlab-runner', v_task_10.task_id);
        
        v_task_12 = create_task(v_document_id => document_id, v_content => 'gitlab-runner --(vpn)--> DNS ----> gitlab');
        v_task_13 = create_task(document_id, 'Firewall & unbound rules update from prov script (VPN)', v_task_12.task_id);
        v_task_14 = create_task(document_id, 'Script pic_gitlab_runner_post_creation', v_task_13.task_id);
        v_task_15 = create_task(document_id, 'Execute 1 CI/CD pipeline on gitlab-runner', v_task_14.task_id);
        v_task_16 = create_task(document_id, 'gitlab-runner provisioner with resolv.conf/docker/runner-register', v_task_12.task_id);
        v_task_17 = create_task(document_id, 'prepare gitlab-ci for ovpn-auth repo', v_task_13.task_id);
        v_task_18 = create_task(document_id, 'PAM for SSH', v_task_14.task_id);
        v_task_19 = create_task(document_id, 'ACL - Redis', v_task_13.task_id);
        v_task_20 = create_task(document_id, 'Redis security', v_task_13.task_id);
        v_task_21 = create_task(document_id, 'TOTP for SSH', v_task_14.task_id);
        v_task_22 = create_task(document_id, 'API gateway without redis', v_task_15.task_id);
        v_task_23 = create_task(document_id, 'Golden image interitance re-organize', v_task_16.task_id);
        v_task_24 = create_task(document_id, 'Postgres', v_task_12.task_id);
        v_task_25 = create_task(document_id, 'Auth service', v_task_13.task_id);
        v_task_26 = create_task(document_id, 'MQ', v_task_15.task_id);
        v_task_27 = create_task(document_id, 'Federated learning', v_task_16.task_id);
        v_task_28 = create_task(document_id, 'Bluetooth transmission test', v_task_12.task_id);
        v_task_29 = create_task(document_id, 'Intrusion detection system (centralised) (OSSEC', v_task_12.task_id);
        v_task_30 = create_task(document_id, 'Envoy - HAProxy - NGiNX', v_task_12.task_id);
        v_task_31 = create_task(document_id, 'web-front/Privacy against [friend/pubic/company/attackers]', v_task_13.task_id);
        v_task_32 = create_task(document_id, 'Redis/cluster script test for multi datacenter', v_task_13.task_id);
        v_task_33 = create_task(document_id, 'gitlab-runner firewall rules: close public internet', v_task_14.task_id);
        v_task_34 = create_task(document_id, 'static-challange for ovpn-auth', v_task_20.task_id);
        v_task_35 = create_task(document_id, 'Golden image for vpn server', v_task_21.task_id);
    END
$$ LANGUAGE 'plpgsql';

-- CALL load_test_dataset();