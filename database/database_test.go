package database

import "testing"

func TestDatabaseForGetAndCreateFunctions(t *testing.T) {

	// Initialize database connection with test database
	Init("postgres://postgres:password@localhost:5432/testdatabase")
	defer Close()

	var (
		test_document     Document
		test_drawer_group TaskGroup
		test_task         Task
		err               []error
	)

	test_document, err = CreateDocument(Document{})
	if err != nil {
		t.Fatalf("Failed on create 'document' object.\nError details: %s", err)
	}

	test_drawer_group, err = CreateTaskGroup(TaskGroup{
		DocumentId:    test_document.DocumentId,
		TaskGroupType: Drawer,
	})
	if err != nil {
		t.Fatalf("Failed on create 'task_group' object.\nError details: %s", err)
	}

	test_task, err = CreateTask(Task{
		Content:     "One difficult task for testing go package.",
		Degree:      1,
		Depth:       1,
		TaskGroupId: test_drawer_group.TaskGroupId,
		TaskStatus:  Active,
		ParentId:    "00000000-0000-0000-0000-000000000000",
	})
	if err != nil {
		t.Fatalf("Failed on create 'task' object.\nError details: %s", err)
	}

	verify_task, err := GetTaskByTaskId(test_task.TaskId)
	if err != nil {
		t.Fatalf("Failed on getTaskByTaskId(test_task.TaskId)\nError details: %s", err)
	}
	if verify_task.Content != test_task.Content {
		t.Fatalf("Failed on comparing result of getTaskByTaskId(test_task.TaskId)")
	}

	verify_drawer_group, err := GetTaskGroupByTaskGroupId(test_drawer_group.TaskGroupId)
	if err != nil {
		t.Fatalf("Failed on getTaskGroupByGroupId(test_drawer_group.TaskGroupId)\nError details: %s", err)
	}
	if verify_drawer_group.CreatedAt != test_drawer_group.CreatedAt {
		t.Fatalf("Failed on comparing result of getTaskGroupByGroupId(test_drawer_group.TaskGroupId)")
	}

	verify_document, err := GetDocumentByDocumentId(test_document.DocumentId)
	if err != nil {
		t.Fatalf("Failed on getDocumentByDocumentId(test_document.DocumentId)\nError details: %s", err)
	}
	if verify_document.CreatedAt != test_document.CreatedAt {
		t.Fatalf("Failed on comparing result of getDocumentByDocumentId(test_document.DocumentId)")
	}
}
