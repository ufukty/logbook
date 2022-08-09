/*
    FIXME: CHECK USER INPUT BEFORE START TO EXECUTION
    FIXME: INSERT A COLUMN TO TASKS TABLE FOR TASK ACTIVATION STATUS (BOOLEAN) 
*/

DROP DATABASE IF EXISTS logbook_dev;

CREATE DATABASE logbook_dev;

\c logbook_dev;

-- MARK: Enums

CREATE TYPE "ACCESS_EVENT_TYPE" AS ENUM (
    'SUCCESSFUL_LOGIN',
    'FAILED_LOGIN_ATTEMPT',
    'LOGOUT'
);

CREATE TYPE "CONTENT_TYPE" AS ENUM (
    'USER_GENERATED_TASK',
    'JOINED_BY_INVITATION',
    'DUPLICATED_FROM_BLUEPRINT'
);

CREATE TYPE "COLLABORATION_SETTINGS_UPDATE_TYPE" AS ENUM (
    'CREATION',
    'INVITATION',
    'JOIN',
    'REOPEN',
    'RESTART',  -- delete all sub-tasks and start over to resolution
    'COMPLETION'
);

-- MARK: Table definitions

CREATE TABLE "USER" (
    "user_id"                   UUID                        UNIQUE DEFAULT gen_random_UUID(),

    "username"                  TEXT                        NOT NULL,
    "email_address"             TEXT                        NOT NULL, 
    "email_address_truncated"   TEXT                        NOT NULL UNIQUE, -- dots and other non alpha-numerical characters wiped on the username. @ and domain saved
    "password_hash_encoded"     TEXT                        NOT NULL, -- should be start with algo name, parameters, salt and resulting hash
    
    "activated"                 BOOLEAN                     DEFAULT FALSE,
    
    "created_at"                TIMESTAMP                   DEFAULT CURRENT_TIMESTAMP
);

-- TODO: BOOKMARK TABLE

CREATE TABLE "ACCESS_EVENT_LOG"(
    "user_id"                   UUID                        NOT NULL REFERENCES "USER" ("user_id"),
    "event_type"                "ACCESS_EVENT_TYPE"         NOT NULL,
    "user-agent"                TEXT,
    "ip-address"                INET                        NOT NULL,
    "created_at"                TIMESTAMP                   DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE "OPERATION_SUMMARY" AS ENUM (
    'CREATE',
    'REORDER',
    'DELETE',
    'CONTENT_EDIT',
    'NOTE_CREATE',
    'NOTE_EDIT',
    'NOTE_DELETE',
    'MARK_COMPLETE',
    'MARK_UNCOMPLETE',
    'COLLABORATION_ASSIGN',
    'COLLABORATION_UNASSIGN',
    'COLLABORATION_RESTRICT',
    'COLLABORATION_DERESTRICT',
    'COLLABORATION_CHANGE_ROLE', -- explicit/implicit? 
    'HISTORY_ROLLBACK',
    'HISTORY_FASTFORWARD'
);

CREATE TYPE "OPERATION_STATUS" AS ENUM (
    'IN_REVIEW',
    'PRIV_ACCEPTED',
    'PRIV_REJECTED',
    'APPLIED_FASTFORWARD',
    'CONFLICT_DETECTED',
    'MANAGER_SELECTION_IN_REVIEW',
    'MANAGER_SELECTION_ACCEPTED',
    'MANAGER_SELECTION_APPLIED',
    'MANAGER_SELECTION_REJECTED'
);

CREATE TABLE "OPERATION" (
    "operation_id"              UUID                        UNIQUE DEFAULT gen_random_UUID(),

    "revision_id"               UUID                        UNIQUE DEFAULT gen_random_UUID(),
    "previous_revision_id"      UUID                        NOT NULL,

    "user_id"                   UUID                        NOT NULL,
    "operation_summary"         "OPERATION_SUMMARY"    NOT NULL,
    "operation_status"          "OPERATION_STATUS"     NOT NULL,

    "task_id"                   UUID                        NOT NULL, -- id of updated task
    "link_id"                   UUID                        NOT NULL, -- id of updated link

    "created_at"                TIMESTAMP                   DEFAULT CURRENT_TIMESTAMP,
    "archived_at"               TIMESTAMP                   -- user can delete task history items that they don't want to see again in history
);

CREATE TYPE "ROLE" AS ENUM (
    'MANAGER',
    'CREATOR',
    'COLLABORATOR'
);

CREATE TABLE "PRIVILEGE" (
    "revision_id"               UUID                        NOT NULL,
    "task_id"                   UUID                        NOT NULL,
    "role"                      "ROLE"                      NOT NULL,
    "created_at"                TIMESTAMP                   DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "TASK" (
    "revision_id"               UUID                        NOT NULL, -- created by TASK_OPERATION table entry
    "task_id"                   UUID                        DEFAULT gen_random_UUID(), -- inherit from previous version (if there is any). explicitly check, newly created uuid doesn't collide with previous ones.
    -- "sharding_id" UUID, -- might be the id of collaboration root

    "original_creator_user_id"  UUID                        NOT NULL REFERENCES "USER" ("user_id"), -- when creator unassigned from collaboration, don't change this
    "responsible_user_id"       UUID                        NOT NULL REFERENCES "USER" ("user_id"), -- handover resposible role when last one quits

    "content"                   TEXT                        NOT NULL,
    "notes"                     TEXT                        NOT NULL,

    "created_at"                TIMESTAMP                   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "completed_at"              TIMESTAMP                   ,  
    "archived_at"               TIMESTAMP                   ,
    "archived"                  BOOLEAN                     NOT NULL DEFAULT FALSE
);

CREATE TABLE "TASK_LINK" (
    "link_id"                   UUID                        NOT NULL DEFAULT gen_random_UUID(),
    "revision_id"               UUID                        NOT NULL,

    "task_id"                   UUID                        NOT NULL,
    "task_revision_id"          UUID                        NOT NULL,
    
    "super_task_id"             UUID                        NOT NULL,
    "super_task_revision_id"    UUID                        NOT NULL,

    "index"                     INT                         NOT NULL DEFAULT 0, -- each task can have only one index value per each supertask it belongs.
    "primary_link"              BOOLEAN                     NOT NULL DEFAULT FALSE, -- TRUE if the link is created with creation of task

    "created_at"                TIMESTAMP                   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE "TASK_LINK_USER_PREFERENCES" (
    "link_id"                   UUID                        NOT NULL,
    "user_id"                   UUID                        NOT NULL REFERENCES "USER" ("user_id"),
    "fold"                      BOOLEAN                     NOT NULL DEFAULT FALSE
);

CREATE TABLE "TASK_PROPS" (
    "revision_id"               UUID                        NOT NULL, -- revision of task which is props calculated for
    "task_id"                   UUID                        NOT NULL, -- unique identifier for any task
    "user_id"                   UUID                        NOT NULL, -- calculate depth, degree accordingly to the user seeing the task 
    
    "degree"                    INT                         NOT NULL DEFAULT 0,
    "depth"                     INT                         NOT NULL DEFAULT 0,
    "completion_percentage"     REAL                        DEFAULT 0
);

CREATE TABLE "BOOKMARK" (
    "user_id"                   UUID                        NOT NULL,
    "task_id"                   UUID                        NOT NULL,
    "display_name"              TEXT                        ,
    "root_bookmark"             BOOLEAN                     NOT NULL DEFAULT FALSE,
    "created_at"                TIMESTAMP                   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "deleted_at"                TIMESTAMP                   
);


-- CREATE TABLE "BLUEPRINT" (
--     "blueprint_id"              UUID        UNIQUE DEFAULT gen_random_UUID(),
--     "original_user_id"          UUID        NOT NULL REFERENCES "USER" ("user_id"),
--     "original_document_id"      UUID        NOT NULL REFERENCES "DOCUMENT" ("document_id"),
--     "original_task_id"          UUID        NOT NULL REFERENCES "TASK" ("task_id"),
--     -- "entry_point"               UUID        REFERENCES "BLUEPRINT_TASK" ("blueprint_id"),
--     "progress_tracking"         BOOLEAN     NOT NULL DEFAULT FALSE,
--     "sharing_uri"               TEXT        UNIQUE NOT NULL,
    
--     "created_at"                TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    
--     "archived"                  BOOLEAN     NOT NULL DEFAULT FALSE
-- );

-- -- Freeze 
-- CREATE TABLE "BLUEPRINT_TASK" (
--     "blueprint_task_id"         UUID        UNIQUE DEFAULT gen_random_UUID(),
--     "blueprint_id"              UUID        NOT NULL,

--     "content"                   TEXT        NOT NULL,

--     "parent_blueprint_task_id"  UUID        NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
--     "degree"                    INT         NOT NULL DEFAULT 0,
--     "depth"                     INT         NOT NULL DEFAULT 0,
--     "index"                     INT         NOT NULL DEFAULT 0,

--     "created_at"                TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
-- );

-- -- entry_point should be nullable to avoid circular dependency
-- ALTER TABLE "BLUEPRINT" ADD "entry_point" UUID REFERENCES "BLUEPRINT_TASK" ("blueprint_id"); 

-- CREATE TABLE "BLUEPRINT_RECEIVER_CONFIGURATION" (
--     "user_id"                   UUID        REFERENCES "USER" ("user_id"),
--     "blueprint_id"              UUID        REFERENCES "BLUEPRINT" ("blueprint_id"),
--     "progress_tracking"         BOOLEAN     DEFAULT FALSE,
--     "placed_document_id"        UUID        REFERENCES "DOCUMENT" ("document_id"),
--     "entry_point"               UUID        REFERENCES "TASK" ("task_id")
-- );

-- CREATE TABLE "TAKES" (
--     "take_id"                   UUID,
--     "user_id"                   UUID,
--     "document_id"               UUID,
--     "original_task_id"          UUID,
--     "entry_task_id"             UUID
-- );

-- MARK: Procedure definitions

-- CREATE VIEW tasks_linearized AS SELECT * FROM tasks;

-- DROP FUNCTION IF EXISTS create_document_with_task_groups;
-- CREATE FUNCTION create_document(
--     v_user_id UUID
-- ) 
-- RETURNS "DOCUMENT" 
-- AS $$
--     DECLARE
--         document "DOCUMENT";
--     BEGIN
--         INSERT INTO "DOCUMENT" ("user_id") VALUES (v_user_id) RETURNING * INTO document;
--         RETURN document;
--     END
-- $$ LANGUAGE 'plpgsql';


-- -- CREATE FUNCTION hierarchical_overview


-- CREATE FUNCTION hierarchical_placement(
-- 	v_user_id UUID,
-- 	v_document_id UUID,
-- 	v_parent_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
-- ) RETURNS TABLE ("task_id" UUID) AS $$
-- 	DECLARE
-- 		v_child_id UUID;
-- 	BEGIN

--         IF v_parent_id <> '00000000-0000-0000-0000-000000000000'
--         THEN
--             RETURN QUERY SELECT v_parent_id;
--         END IF;
	
-- 		FOR v_child_id IN 
		
-- 			SELECT t."task_id" 
-- 			FROM "TASK" as t
-- 			WHERE "user_id" = v_user_id
-- 				AND "document_id" = v_document_id
-- 				AND "parent_id" = v_parent_id
--                 AND "archived" = FALSE
--                 AND "fold" = FALSE
-- 			ORDER BY "index"

-- 		LOOP
		
-- 			-- TODO: cycle detection

--             -- RAISE NOTICE 'parent_id: % child_id: %', v_parent_id, v_child_id;
			
-- 			RETURN QUERY 
-- 				SELECT r."task_id"
--                 FROM hierarchical_placement(
-- 					v_user_id,
-- 					v_document_id,
-- 					v_child_id
-- 				) AS r;

-- 		END LOOP;
-- 	END
-- $$ LANGUAGE 'plpgsql';


-- -- -- TO GET THE LIST OF TASKS THAT ARE:
-- -- --   * COMPLETED WITHIN LAST 10 DAYS (<200 ITEM),
-- -- --   * ACTIVE,
-- -- --   * READY-TO-PICK-UP
-- -- --   * DRAWER (TO-DO)
-- -- CREATE FUNCTION document_overview(
-- --     v_user_id UUID, 
-- --     v_document_id UUID, 
-- --     v_limit INT DEFAULT 1000, 
-- --     v_offset INT DEFAULT 0
-- -- ) 
-- -- RETURNS SETOF "TASK" 
-- -- AS $$
-- --     DECLARE
-- --         last_row_selected       "DOCUMENT_CONTENT";
-- --         dfs_next_content_id     UUID;
-- --         limit_counter           INT DEFAULT 1;
-- --     BEGIN
       
-- --         SELECT *
-- --             INTO last_row_selected
-- --             FROM "DOCUMENT_CONTENT" dc
-- --             WHERE 
-- --                 dc."document_id" = v_document_id
-- --                 AND dc."user_id" = v_user_id
-- --             ORDER BY dc."created_at" ASC
-- --             LIMIT 1
-- --             OFFSET v_offset;

-- --         IF NOT FOUND THEN
-- --             RAISE EXCEPTION 'Document does not found.';
-- --         END IF;

-- --         RETURN NEXT last_row_selected;

-- --         LOOP

-- --             dfs_next_content_id := last_row_selected."next_content_id";
-- --             IF dfs_next_content_id IS NULL THEN
-- --                 RETURN;
-- --             END IF;

-- --             SELECT *
-- --                 INTO last_row_selected
-- --                 FROM "DOCUMENT_CONTENT" dc
-- --                 WHERE 
-- --                     dc."document_id" = v_document_id
-- --                     AND dc."user_id" = v_user_id
-- --                     AND dc."task_id" = dfs_next_content_id;

-- --             RETURN NEXT last_row_selected;

-- --             limit_counter := limit_counter + 1;
-- --             IF limit_counter = v_limit THEN
-- --                 RETURN;
-- --             END IF;

-- --         END LOOP;


-- --     END
-- -- $$ LANGUAGE 'plpgsql';

-- -- RECURSIVE HELPER FUNCTION FOR:
-- --    * create_task
-- --    * reattach_task
-- CREATE FUNCTION update_task_degree(v_task_id UUID, v_increment INT) RETURNS UUID[] AS $$
--     DECLARE
--         v_total_degrees_of_siblings INT;
--         v_task "TASK";
--         v_updated_task_list UUID[];
--     BEGIN
--         -- RAISE NOTICE 'update_task_degree, v_task_id = %', v_task_id;

--         UPDATE "TASK"
--         SET "degree" = "degree" + v_increment
--         WHERE "TASK"."task_id" = v_task_id;

--         SELECT *
--         INTO v_task
--         FROM "TASK"
--         WHERE "task_id" = v_task_id;

--         -- RAISE NOTICE 'v_task = %', v_task;
--         -- RAISE NOTICE 'v_task."parent_id" = %', v_task."parent_id";

--         -- ADD ITSELF TO v_updated_task_list BEFORE THE TASKS THAT WILL RETURNED BY PARENT
--         v_updated_task_list = array_append(v_updated_task_list, v_task."task_id");

--         IF v_task."parent_id" != '00000000-0000-0000-0000-000000000000' THEN
--             -- RAISE NOTICE 'recursing into parent';
--             v_updated_task_list = array_cat(
--                 v_updated_task_list, 
--                 update_task_degree(v_task."parent_id", v_increment)
--             );
--         END IF;
        
--         -- RAISE NOTICE 'no more parent to recurse further, returning to caller now';

--         RETURN v_updated_task_list;
--     END
-- $$ LANGUAGE 'plpgsql';

-- CREATE FUNCTION update_task_readiness(v_task_id UUID) RETURNS UUID AS $$
--     DECLARE
--         v_undone_children "TASK";
--         v_readiness BOOLEAN;
--         v_task "TASK";
--     BEGIN
--         -- TODO: IT CAN RETURN UUID[] IF PARENTS OF PARENTS TAKEN INTO CALCULATION

--         -- RAISE NOTICE 'update_task_readiness is running for %', v_task_id;

--         SELECT *
--         INTO v_undone_children
--         FROM "TASK"
--         WHERE "parent_id" = v_task_id
--             AND "completed_at" IS NULL;

--         IF v_undone_children IS NULL THEN
--             v_readiness = TRUE;
--         ELSE
--             v_readiness = FALSE;
--         END IF;

--         -- RAISE NOTICE 'v_readiness = %', v_readiness;

--         UPDATE "TASK"
--         SET "ready_to_pick_up" = v_readiness
--         WHERE "task_id" = v_task_id 
--             AND "ready_to_pick_up" <> v_readiness
--         RETURNING * INTO v_task;

--         RETURN v_task."task_id";
--     END
-- $$ LANGUAGE 'plpgsql';

-- -- RECURSIVE HELPER FUNCTION FOR:
-- --     * update_task_depth
-- -- IT UPDATES THE DEPTH OF TASK AND RECURSES INTO ITS CHILDREN
-- CREATE FUNCTION update_task_depth_helper(v_task_id UUID, v_new_depth INT) RETURNS UUID[] AS $$
--     DECLARE
--         v_child_id UUID;
--         v_updated_task_id UUID;
--         v_updated_task_list UUID[];
--     BEGIN
--         UPDATE "TASK"
--         SET "depth" = v_new_depth
--         WHERE "task_id" = v_task_id
--         RETURNING "task_id" INTO v_updated_task_id;

--         IF v_updated_task_id IS NULL THEN
--             RETURN v_updated_task_list;
--         END IF;

--         FOR v_child_id IN (SELECT * FROM "TASK" WHERE "parent_id" = v_task_id) 
--         LOOP
--             v_updated_task_list = array_cat(
--                 v_updated_task_list,
--                 update_task_depth_helper(v_child_id, v_new_depth + 1)
--             );
--         END LOOP;

--         RETURN v_updated_task_list;
--     END
-- $$ LANGUAGE 'plpgsql';

-- -- RECURSIVE HELPER FUNCTION FOR:
-- --     * create_task
-- --     * reattach_task
-- CREATE FUNCTION update_task_depth(v_task_id UUID) RETURNS UUID[] AS $$
--     DECLARE
--         v_task_depth INT;
--         v_task "TASK";
--         -- v_root_uuid UUID = '00000000-0000-0000-0000-000000000000';
--     BEGIN
--         SELECT *
--         INTO v_task
--         FROM "TASK"
--         WHERE "task_id" = v_task_id;

--         IF v_task."parent_id" = '00000000-0000-0000-0000-000000000000' THEN
--             v_task_depth = 1;
--         ELSE
--             SELECT "depth" + 1
--             INTO v_task_depth
--             FROM "TASK"
--             WHERE "task_id" = v_task_id;
--         END IF;

--         RETURN update_task_depth_helper(v_task_id, v_task_depth);
--     END
-- $$ LANGUAGE 'plpgsql';


-- CREATE FUNCTION array_to_sorted_row_list(v_id_array UUID[], v_task_id UUID) RETURNS SETOF "TASK" AS $$
--     BEGIN
--         RETURN QUERY (
--             SELECT * 
--             FROM "TASK" t 
--             WHERE t."task_id" = ANY(v_id_array) 
--             ORDER BY CASE WHEN t."task_id" = v_task_id THEN 0 ELSE 1 END
--         );
--     END
-- $$ LANGUAGE 'plpgsql';

-- -- Add new task to document with:
-- --   * updating the degree of parent task (and theirs, recursively).
-- --   * minding the depth of parent task.
-- --   * create references in "DOCUMENT_CONTENT" table and "TASK_HISTORY" table
-- -- Returns the list of updated tasks.
-- -- FIXME: Fill "index" column correctly for new task
-- CREATE FUNCTION create_task(
--     v_document_id UUID, 
--     v_user_id UUID, 
--     v_content VARCHAR, 
--     v_index INTEGER,
--     v_parent_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'
-- ) RETURNS SETOF "TASK" AS $$
--     DECLARE
--         v_depth INT;
--         v_degree INT;
--         v_task "TASK"%ROWTYPE;
--         v_updated_task_list UUID[];
--     BEGIN
--         -- RAISE NOTICE 'v_content = %, v_parent_id = %', v_content, v_parent_id;

--         -- DEGREE ALWAYS 1, WHEN THE TASK IS NEWLY ADDED
--         v_degree = 1;

--         IF v_parent_id IS NULL THEN
--             RAISE EXCEPTION 'Parent ID can not be NULL';
--             RETURN;
--         END IF;

--         -- DECIDE DEPTH
--         IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN

--             SELECT "TASK"."depth"+1 
--             INTO v_depth
--             FROM "TASK" 
--             WHERE "task_id" = v_parent_id;
--         ELSE
--             v_depth = 1;
--         END IF;

--         -- WRITE NEW TASK INTO DB
--         INSERT INTO "TASK"("document_id", "user_id", "parent_id", "content", "degree", "depth", "index")
--         VALUES (v_document_id, v_user_id, v_parent_id, v_content, v_degree, v_depth, v_index)
--         RETURNING * INTO v_task;

--         -- INITIALIZE THE RETURN LIST
--         v_updated_task_list = array_append(v_updated_task_list, v_task."task_id");

--         -- UPDATE PARENT'S READY-TO-PICK-UP STATUS TO FALSE
--         IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN
--             UPDATE "TASK"
--             SET "ready_to_pick_up" = FALSE
--             WHERE task_id = v_parent_id;
--         END IF;

--         -- UPDATE PARENTS' DEGREES
--         IF v_parent_id != '00000000-0000-0000-0000-000000000000' THEN
--             v_updated_task_list = array_cat(
--                 v_updated_task_list, 
--                 update_task_degree(v_parent_id, 1)
--             );
--         END IF;

--         RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task."task_id");
--     END
-- $$ LANGUAGE 'plpgsql';

-- CREATE FUNCTION reattach_task(v_task_id UUID, v_new_parent_id UUID) RETURNS SETOF "TASK" AS $$ -- RETURN "TASK"
--     DECLARE
--         v_updated_task_list UUID[];
--         v_task_old "TASK";
--     BEGIN
--         -- FIXME: CHECK CIRCULAR DEPENDENCY
     
--         -- TEMPORARILY STORE THE TASK WITH ITS CURRENT CONDITION
--         SELECT * 
--         INTO v_task_old 
--         FROM "TASK" 
--         WHERE "task_id" = v_task_id;

--         -- DON'T CONTINUE IF THE DESIRED NEW PARENT IS THE TASK'S ITSELF
--         IF v_task_id = v_new_parent_id THEN
--             RETURN;
--         END IF;

--         -- UPDATE TASK:
--         --     * PARENT
--         --     * DEPTH
--         UPDATE "TASK"
--         SET "parent_id" = v_new_parent_id
--         WHERE "task_id" = v_task_id;

--         -- INITILIAZE RETURN LIST WITH THE MODIFIED TASK AS ITS FIRST ITEM
--         v_updated_task_list = array_append(v_updated_task_list, v_task_old."task_id");

--         -- UPDATE OLD PARENT:
--         --     * DEGREE (RECURSIVELY)
--         --     * READINESS STATUS
--         IF v_task_old."parent_id" != '00000000-0000-0000-0000-000000000000' THEN
--             v_updated_task_list = array_cat(
--                 v_updated_task_list,
--                 update_task_degree(v_task_old."parent_id", -1 * v_task_old."degree")
--             );
--             PERFORM update_task_readiness(v_task_old."parent_id");
--         END IF;

--         -- UPDATE NEW PARENT:
--         --     * DEGREE (RECURSIVELY)
--         --     * DEPTH (RECURSIVELY)
--         --     * READINESS STATUS
--         v_updated_task_list = array_cat(
--             v_updated_task_list,
--             update_task_degree(v_new_parent_id, v_task_old."degree")
--         );
--         PERFORM update_task_readiness(v_new_parent_id);
--         v_updated_task_list = array_cat(
--             v_updated_task_list,
--             update_task_depth(v_task_id)
--         );

--         -- RETURN UPDATED TASKS AS ARRAY FOR UPDATING FRONTEND 
--         RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task_id);
--     END 
-- $$ LANGUAGE 'plpgsql';

-- CREATE FUNCTION mark_a_task_done(v_task_id UUID) RETURNS SETOF "TASK" AS $$
--     DECLARE
--         v_task "TASK";
--         v_updated_task_list UUID[];
--     BEGIN
--         -- RAISE NOTICE 'mark_a_task_done, v_task_id = %', v_task_id;

--         -- UPDATE TASK'S ITSELF
--         UPDATE "TASK"
--         SET completed_at = CURRENT_TIMESTAMP
--         WHERE "task_id" = v_task_id 
--             AND "completed_at" IS NULL
--         RETURNING * INTO v_task;
        
--         -- UPDATE TASK'
--         v_updated_task_list = array_append(v_updated_task_list, v_task."task_id");

--         -- UPDATE PARENT READINESS
--         IF v_task IS NOT NULL THEN
--             v_updated_task_list = array_append(
--                 v_updated_task_list,
--                 update_task_readiness(v_task."parent_id") 
--             );
--         END IF;

--         RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task_id);
--     END
-- $$ LANGUAGE 'plpgsql';

-- -- CREATE FUNCTION update_task(
-- --     task_id UUID,
-- --     user_id UUID,
-- --     document_id UUID,
-- --     content TEXT,
-- --     parent_id
-- -- ) RETURNS TABLE ("task_id" UUID) AS $$
-- --     DECLARE
    
-- --     BEGIN

-- --     END
-- -- $$ LANGUAGE 'plpgsql';

-- -- get all takes of a 

-- -- CREATE FUNCTION get_all_takes_of_a_task("user_id", "document_id", "") RETURNS VOID AS $$ $$ LANGUAGE "plpgsql";



-- -- CREATE FUNCTION create_blueprint_from_task(v_user_id UUID, v_document_id UUID, v_task_id UUID, progress_tracking BOOLEAN, enabled_tasks UUID[]) RETURNS VOID AS $$
-- --     DECLARE
-- --         v_blueprint_id UUID,

-- --     BEGIN
-- --         -- check degree of root task is under allowed limit for user
-- --         SELECT degree
-- --             FROM "TASK"
-- --             WHERE 
-- --                 "user_id" = v_user_id
-- --                 AND "document_id" = v_document_id
-- --                 AND "task_id" = v_task_id;
-- --         IF degree > allowed_max_bluelimit_degree THEN
-- --             RAISE EXCEPTION 'Depth of task exceeds the limit.';
-- --         END IF;

-- --         -- create a "BLUEPRINT" record, keep its ID
-- --         INSERT INTO "BLUEPRINT"(
-- --             "original_user_id", 
-- --             "original_document_id", 
-- --             "original_task_id",
-- --             "progress_tracking",

-- --         ) 
-- --         VALUES()
-- --         RETURNING "blueprint_id" INTO v_blueprint_id;

-- --         -- with blueprint_id, create a "BLUEPRINT_TASK" row for "TASK" keep its ID

-- --         -- dfs on "TASK" starting from root, children found by matches on "parent_id", repeat previous step

-- --         -- save the blueprint_id of root task to "BLUEPRINT" record
-- --         UPDATE TABLE "BLUEPRINT"
-- --             SET "entry_point" = v_blueprint_root_id
-- --             WHERE "blueprint_id" = v_blueprint_id;

-- --         RETURN QUERY (
-- --             SELECT sharing_uri 
-- --                 FROM "BLUEPRINT" 
-- --                 WHERE "blueprint_id" = v_blueprint_id;
-- --         );
-- --     END
-- -- $$ LANGUAGE "plpgsql";



-- -- CREATE FUNCTION () RETURNS VOID AS $$ $$ LANGUAGE "plpgsql";
-- -- CREATE FUNCTION () RETURNS VOID AS $$ $$ LANGUAGE "plpgsql";
-- -- CREATE FUNCTION () RETURNS VOID AS $$ $$ LANGUAGE "plpgsql";
-- -- CREATE FUNCTION () RETURNS VOID AS $$ $$ LANGUAGE "plpgsql";
-- -- CREATE FUNCTION () RETURNS VOID AS $$ $$ LANGUAGE "plpgsql";





-- -- MARK: Test dataset

-- CREATE PROCEDURE load_test_dataset() AS $$
--     DECLARE
--         v_user_id                 UUID DEFAULT '13600fd8-2c4a-50df-80e2-5cc0d0e711df';
--         v_username                TEXT DEFAULT 'Name Surname';
--         v_user_email_address      TEXT DEFAULT 'admin@logbook';
--         v_email_address_truncated TEXT DEFAULT 'admin@logbook';
--         v_password_hash_encoded   TEXT DEFAULT 'Ft7isEJgIrKdWgA496C9GnPHvAhlo2x4';
--         v_document_id             UUID DEFAULT '61bbc44a-c61c-4d49-8804-486181081fa7';
--         v_task_1                  UUID;
--         v_task_2                  UUID;
--         v_task_3                  UUID;
--         v_task_4                  UUID;
--         v_task_5                  UUID;
--         v_task_6                  UUID;
--         v_task_7                  UUID;
--         v_task_8                  UUID;
--         v_task_9                  UUID;
--         v_task_10                 UUID;
--         v_task_11                 UUID;
--         v_task_12                 UUID;
--         v_task_13                 UUID;
--         v_task_14                 UUID;
--         v_task_15                 UUID;
--         v_task_16                 UUID;
--         v_task_17                 UUID;
--         v_task_18                 UUID;
--         v_task_19                 UUID;
--         v_task_20                 UUID;
--         v_task_21                 UUID;
--         v_task_22                 UUID;
--         v_task_23                 UUID;
--         v_task_24                 UUID;
--         v_task_25                 UUID;
--         v_task_26                 UUID;
--         v_task_27                 UUID;
--         v_task_28                 UUID;
--         v_task_29                 UUID;
--         v_task_30                 UUID;
--         v_task_31                 UUID;
--         v_task_32                 UUID;
--         v_task_33                 UUID;
--         v_task_34                 UUID;
--         v_task_35                 UUID;
--     BEGIN
--         -- SELECT "document_id" INTO v_document_id FROM create_document();
--         INSERT INTO "USER"("user_id", "username", "email_address", "email_address_truncated", "password_hash_encoded") 
--             VALUES (v_user_id, v_username, v_user_email_address, v_email_address_truncated, v_password_hash_encoded);

--         INSERT INTO "DOCUMENT"("document_id", "user_id") VALUES (v_document_id, v_user_id);
--         RAISE NOTICE 'user_id: %', v_user_id;
--         RAISE NOTICE 'document_id: %', v_document_id; 

--         -- FIRST ROOT TASK

--         SELECT "task_id" INTO v_task_1 FROM create_task(v_document_id => v_document_id, v_user_id => v_user_id, v_content => 'deploy redis cluster on multi DC', v_index => 0);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '40 DAYS' WHERE "task_id" = v_task_1;
        
--         SELECT "task_id" INTO v_task_2 FROM create_task(v_document_id, v_user_id, 'deploy redis cluster on 1 DC', 0, v_task_1);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '39 DAYS' - INTERVAL '14 HOURS' WHERE "task_id" = v_task_2;

--         SELECT "task_id" INTO v_task_3 FROM create_task(v_document_id, v_user_id, 'Revoke passwordless sudo rights after provision at cluster', 1, v_task_1);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '39 DAYS' - INTERVAL '13 HOURS' WHERE "task_id" = v_task_18;
        
--         SELECT "task_id" INTO v_task_4 FROM create_task(v_document_id, v_user_id, 'iptables for redis', 0, v_task_3);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '39 DAYS' - INTERVAL '12 HOURS' WHERE "task_id" = v_task_29;
        
--         SELECT "task_id" INTO v_task_5 FROM create_task(v_document_id, v_user_id, 'terraform for redis', 0, v_task_4);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '39 DAYS' - INTERVAL '11 HOURS' WHERE "task_id" = v_task_3;
        
--         SELECT "task_id" INTO v_task_6 FROM create_task(v_document_id, v_user_id, 'Update redis/tf file according to prod.tfvars file', 1, v_task_4);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '39 DAYS' - INTERVAL '10 HOURS' WHERE "task_id" = v_task_6;

--         SELECT "task_id" INTO v_task_7 FROM create_task(v_document_id, v_user_id, 'Remove: seperator from ovpn-auth', 0, v_task_2);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '34 DAYS' WHERE "task_id" = v_task_7;
        
--         SELECT "task_id" INTO v_task_8 FROM create_task(v_document_id, v_user_id, 'Write tests for ovpn-auth', 1, v_task_3);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '33 DAYS' - INTERVAL '12 HOURS' WHERE "task_id" = v_task_8;
        
--         SELECT "task_id" INTO v_task_9 FROM create_task(v_document_id, v_user_id, 'Decrease timing gap of ovpn-auth under 1ms', 2, v_task_3);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '33 DAYS' - INTERVAL '11 HOURS' WHERE "task_id" = v_task_9;
        
--         SELECT "task_id" INTO v_task_10 FROM create_task(v_document_id, v_user_id, 'Prepare releases for ovpn-auth', 2, v_task_4);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '33 DAYS' - INTERVAL '10 HOURS' WHERE "task_id" = v_task_10;
        
--         SELECT "task_id" INTO v_task_11 FROM create_task(v_document_id, v_user_id, 'Provision golden-image for gitlab-runner', 0, v_task_10);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '30 DAYS' WHERE "task_id" = v_task_11;

--         -- SECOND ROOT TASK

--         SELECT "task_id" INTO v_task_12 FROM create_task(v_document_id => v_document_id, v_user_id => v_user_id, v_content => 'gitlab-runner --(vpn)--> DNS ----> gitlab', v_index => 1);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '23 HOURS' WHERE "task_id" = v_task_12;
       
--         SELECT "task_id" INTO v_task_13 FROM create_task(v_document_id, v_user_id, 'Firewall & unbound rules update from prov script (VPN)', 0, v_task_12);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '22 HOURS' WHERE "task_id" = v_task_13;
        
--         SELECT "task_id" INTO v_task_14 FROM create_task(v_document_id, v_user_id, 'Script pic_gitlab_runner_post_creation', 0, v_task_13);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '21 HOURS' WHERE "task_id" = v_task_14;
        
--         SELECT "task_id" INTO v_task_15 FROM create_task(v_document_id, v_user_id, 'Execute 1 CI/CD pipeline on gitlab-runner', 0, v_task_14);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '20 HOURS' WHERE "task_id" = v_task_15;
        
--         SELECT "task_id" INTO v_task_16 FROM create_task(v_document_id, v_user_id, 'gitlab-runner provisioner with resolv.conf/docker/runner-register', 1, v_task_12);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '19 HOURS' WHERE "task_id" = v_task_16;
        
--         SELECT "task_id" INTO v_task_17 FROM create_task(v_document_id, v_user_id, 'prepare gitlab-ci for ovpn-auth repo', 1, v_task_13);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '17 HOURS' WHERE "task_id" = v_task_17;
        
--         SELECT "task_id" INTO v_task_18 FROM create_task(v_document_id, v_user_id, 'PAM for SSH', 1, v_task_14);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '16 HOURS' WHERE "task_id" = v_task_18;
        
--         SELECT "task_id" INTO v_task_19 FROM create_task(v_document_id, v_user_id, 'ACL - Redis', 2, v_task_13);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '28 DAYS' - INTERVAL '15 HOURS' WHERE "task_id" = v_task_19;
        
--         SELECT "task_id" INTO v_task_20 FROM create_task(v_document_id, v_user_id, 'Redis security', 3, v_task_13);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '27 DAYS' WHERE "task_id" = v_task_20;
        
--         SELECT "task_id" INTO v_task_21 FROM create_task(v_document_id, v_user_id, 'TOTP for SSH', 2, v_task_14);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '26 DAYS' WHERE "task_id" = v_task_21;
        
--         SELECT "task_id" INTO v_task_22 FROM create_task(v_document_id, v_user_id, 'API gateway without redis', 0, v_task_15);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '25 DAYS' - INTERVAL '10 HOURS' WHERE "task_id" = v_task_22;
        
--         SELECT "task_id" INTO v_task_23 FROM create_task(v_document_id, v_user_id, 'Golden image interitance re-organize', 0, v_task_16);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '25 DAYS' - INTERVAL '9 HOURS' WHERE "task_id" = v_task_23;
        
--         SELECT "task_id" INTO v_task_24 FROM create_task(v_document_id, v_user_id, 'Postgres', 2, v_task_12);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '25 DAYS' - INTERVAL '5 HOURS' WHERE "task_id" = v_task_24;
        
--         SELECT "task_id" INTO v_task_25 FROM create_task(v_document_id, v_user_id, 'Auth service', 4, v_task_13);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '24 DAYS' WHERE "task_id" = v_task_25;
        
--         SELECT "task_id" INTO v_task_26 FROM create_task(v_document_id, v_user_id, 'MQ', 1, v_task_15);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '23 DAYS' - INTERVAL '10 HOURS' WHERE "task_id" = v_task_26;
        
--         SELECT "task_id" INTO v_task_27 FROM create_task(v_document_id, v_user_id, 'Federated learning', 1, v_task_16);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '23 DAYS' - INTERVAL '9 HOURS' WHERE "task_id" = v_task_27;
        
--         SELECT "task_id" INTO v_task_28 FROM create_task(v_document_id, v_user_id, 'Bluetooth transmission test', 3, v_task_12);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '23 DAYS' - INTERVAL '8 HOURS' WHERE "task_id" = v_task_28;
        
--         SELECT "task_id" INTO v_task_29 FROM create_task(v_document_id, v_user_id, 'Intrusion detection system (centralised) (OSSEC', 4, v_task_12);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '22 DAYS' - INTERVAL '15 HOURS' WHERE "task_id" = v_task_29;
        
--         SELECT "task_id" INTO v_task_30 FROM create_task(v_document_id, v_user_id, 'Envoy - HAProxy - NGiNX', 5, v_task_12);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '22 DAYS' - INTERVAL '14 HOURS' WHERE "task_id" = v_task_30;
        
--         SELECT "task_id" INTO v_task_31 FROM create_task(v_document_id, v_user_id, 'web-front/Privacy against [friend/pubic/company/attackers]', 5, v_task_13);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '22 DAYS' - INTERVAL '13 HOURS' WHERE "task_id" = v_task_31;
        
--         SELECT "task_id" INTO v_task_32 FROM create_task(v_document_id, v_user_id, 'Redis/cluster script test for multi datacenter', 6, v_task_13);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '22 DAYS' - INTERVAL '12 HOURS' WHERE "task_id" = v_task_32;
        
--         SELECT "task_id" INTO v_task_33 FROM create_task(v_document_id, v_user_id, 'gitlab-runner firewall rules: close public internet', 3, v_task_14);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '22 DAYS' - INTERVAL '11 HOURS' WHERE "task_id" = v_task_33;
        
--         SELECT "task_id" INTO v_task_34 FROM create_task(v_document_id, v_user_id, 'static-challange for ovpn-auth', 0, v_task_20);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '22 DAYS' - INTERVAL '10 HOURS' WHERE "task_id" = v_task_34;
        
--         SELECT "task_id" INTO v_task_35 FROM create_task(v_document_id, v_user_id, 'Golden image for vpn server', 0, v_task_21);
--         UPDATE "TASK" SET "created_at" = CURRENT_TIMESTAMP - INTERVAL '21 DAYS' WHERE "task_id" = v_task_35;

--         -- COMPLETE SOME TASKS

--         PERFORM mark_a_task_done(v_task_18);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '4 DAYS' WHERE "task_id" = v_task_18;
--         PERFORM mark_a_task_done(v_task_19);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_19;
--         PERFORM mark_a_task_done(v_task_20);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_20;
--         PERFORM mark_a_task_done(v_task_21);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_21;
--         PERFORM mark_a_task_done(v_task_22);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_22;
--         PERFORM mark_a_task_done(v_task_23);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_23;
--         PERFORM mark_a_task_done(v_task_24);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_24;
--         PERFORM mark_a_task_done(v_task_25);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_25;
--         PERFORM mark_a_task_done(v_task_26);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '4 DAYS' WHERE "task_id" = v_task_26;
--         PERFORM mark_a_task_done(v_task_27);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_27;
--         PERFORM mark_a_task_done(v_task_28);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_28;
--         PERFORM mark_a_task_done(v_task_29);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_29;
--         PERFORM mark_a_task_done(v_task_30);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_30;
--         PERFORM mark_a_task_done(v_task_31);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_31;
--         PERFORM mark_a_task_done(v_task_32);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_32;
--         PERFORM mark_a_task_done(v_task_33);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_33;
--         PERFORM mark_a_task_done(v_task_34);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '2 DAYS' WHERE "task_id" = v_task_34;
--         PERFORM mark_a_task_done(v_task_35);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_35;
--         PERFORM mark_a_task_done(v_task_6);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_6;
--         PERFORM mark_a_task_done(v_task_7);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_7;
--         PERFORM mark_a_task_done(v_task_8);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '3 DAYS' WHERE "task_id" = v_task_8;
--         PERFORM mark_a_task_done(v_task_9);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_9;
--         PERFORM mark_a_task_done(v_task_10);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '1 DAYS' WHERE "task_id" = v_task_10;
--         PERFORM mark_a_task_done(v_task_11);
--         UPDATE "TASK" SET "completed_at" = CURRENT_TIMESTAMP + INTERVAL '5 DAYS' WHERE "task_id" = v_task_11;

--         PERFORM reattach_task(v_task_13, v_task_1);
--         PERFORM reattach_task(v_task_27, v_task_15);
--         PERFORM reattach_task(v_task_31, v_task_5); 

--     END
-- $$ LANGUAGE 'plpgsql';

-- CALL load_test_dataset();

-- SELECT * FROM hierarchical_placement(
--     v_user_id => '13600fd8-2c4a-50df-80e2-5cc0d0e711df',
--     v_document_id => '61bbc44a-c61c-4d49-8804-486181081fa7'
-- ); -- OFFSET 2 LIMIT 10

