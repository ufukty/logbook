package database

import "testing"

func TestDatabaseForGetAndCreateFunctions(t *testing.T) {

	// Initialize database connection with test database
	Init("postgres://postgres:password@localhost:5432/testdatabase")
	defer Close()

	var (
		test_document       Document
		test_modified_tasks []Task
		test_task           Task
		err                 []error
	)

	test_document, err = CreateDocument()
	if err != nil {
		t.Fatalf("Failed on create 'document' object.\nError details: %s", err)
	}

	test_modified_tasks, err = CreateTask(Task{
		DocumentId: test_document.DocumentId,
		Content:    "One difficult task for testing go package.",
		ParentId:   "00000000-0000-0000-0000-000000000000",
	})
	if err != nil {
		t.Fatalf("Failed on create 'task' object.\nError details: %s", err)
	}
	test_task = test_modified_tasks[0]

	verify_task, err := GetTaskByTaskId(test_task.TaskId)
	if err != nil {
		t.Fatalf("Failed on getTaskByTaskId(test_task.TaskId)\nError details: %s", err)
	}
	if verify_task.Content != test_task.Content {
		t.Fatalf("Failed on comparing result of getTaskByTaskId(test_task.TaskId)")
	}

	verify_document, err := GetDocumentByDocumentId(test_document.DocumentId)
	if err != nil {
		t.Fatalf("Failed on getDocumentByDocumentId(test_document.DocumentId)\nError details: %s", err)
	}
	if verify_document.CreatedAt != test_document.CreatedAt {
		t.Fatalf("Failed on comparing result of getDocumentByDocumentId(test_document.DocumentId)")
	}
}
