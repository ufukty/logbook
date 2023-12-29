/*
FIXME: CHECK USER INPUT BEFORE START TO EXECUTION
FIXME: INSERT A COLUMN TO TASKS TABLE FOR TASK ACTIVATION STATUS (BOOLEAN) 
 */
DROP DATABASE IF EXISTS logbook_dev;

CREATE DATABASE logbook_dev;

CONNECT logbook_dev;

CREATE TABLE
    "TASK" (
        "task_id" UUID UNIQUE DEFAULT gen_random_UUID (),
        "document_id" UUID NOT NULL REFERENCES "DOCUMENT" ("document_id"),
        "parent_id" UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
        "content" TEXT NOT NULL,
        "degree" INT NOT NULL DEFAULT 1,
        "depth" INT NOT NULL DEFAULT 1,
        "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        "completed_at" DATE, -- FIXME: MAKE IT TIMESTAMP AND UPDATE document_overview FUNCTION
        "ready_to_pick_up" BOOLEAN NOT NULL DEFAULT TRUE
    );

-- RECURSIVE HELPER FUNCTION FOR:
--    * create_task
--    * reattach_task
CREATE FUNCTION update_task_degree (v_task_id UUID, v_increment INT) RETURNS UUID[] AS $$
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

CREATE FUNCTION update_task_readiness (v_task_id UUID) RETURNS UUID AS $$
    DECLARE
        v_undone_children "TASK";
        v_readiness BOOLEAN;
        v_task "TASK";
    BEGIN
        -- TODO: IT CAN RETURN UUID[] IF PARENTS OF PARENTS TAKEN INTO CALCULATION

        -- RAISE NOTICE 'update_task_readiness is running for %', v_task_id;

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
CREATE FUNCTION update_task_depth_helper (v_task_id UUID, v_new_depth INT) RETURNS UUID[] AS $$
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
CREATE FUNCTION update_task_depth (v_task_id UUID) RETURNS UUID[] AS $$
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

CREATE FUNCTION array_to_sorted_row_list (v_id_array UUID[], v_task_id UUID) RETURNS SETOF "TASK" AS $$
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
CREATE FUNCTION create_task (
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

CREATE FUNCTION reattach_task (v_task_id UUID, v_new_parent_id UUID) RETURNS SETOF "TASK" AS $$ -- RETURN "TASK"
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
            PERFORM update_task_readiness(v_task_old."parent_id");
        END IF;

        -- UPDATE NEW PARENT:
        --     * DEGREE (RECURSIVELY)
        --     * DEPTH (RECURSIVELY)
        --     * READINESS STATUS
        v_updated_task_list = array_cat(
            v_updated_task_list,
            update_task_degree(v_new_parent_id, v_task_old."degree")
        );
        PERFORM update_task_readiness(v_new_parent_id);
        v_updated_task_list = array_cat(
            v_updated_task_list,
            update_task_depth(v_task_id)
        );

        -- RETURN UPDATED TASKS AS ARRAY FOR UPDATING FRONTEND 
        RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task_id);
    END 
$$ LANGUAGE 'plpgsql';

CREATE FUNCTION mark_a_task_done (v_task_id UUID) RETURNS SETOF "TASK" AS $$
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
                update_task_readiness(v_task."parent_id") 
            );
        END IF;

        RETURN QUERY SELECT * FROM array_to_sorted_row_list(v_updated_task_list, v_task_id);
    END
$$ LANGUAGE 'plpgsql';
