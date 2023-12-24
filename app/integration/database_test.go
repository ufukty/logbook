package database

import (
	"fmt"
	"testing"
)

func TestDatabaseForGetAndCreateFunctions(t *testing.T) {

	// Initialize database connection with test database
	db, err := New("postgres://postgres:password@localhost:5432/testdatabase")
	if err != nil {
		t.Fatal(fmt.Errorf("prep, New: %w", err))
	}
	defer db.Close()

	doc, err = db.CreateDocument()
	if err != nil {
		t.Fatal(fmt.Errorf("Failed on create 'document' object.: %w", err))
	}

	test_modified_tasks, err = CreateTask(Task{
		DocumentId: doc.DocumentId,
		Content:    "One difficult task for testing go package.",
		ParentId:   "00000000-0000-0000-0000-000000000000",
	})
	if err != nil {
		t.Fatal(fmt.Errorf("Failed on create 'task' object.: %w", err))
	}
	test_task = test_modified_tasks[0]

	verify_task, err := GetTaskByTaskId(test_task.TaskId)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed on getTaskByTaskId(test_task.TaskId): %w", err))
	}
	if verify_task.Content != test_task.Content {
		t.Fatalf("Failed on comparing result of getTaskByTaskId(test_task.TaskId)")
	}

	verify_document, err := GetDocumentByDocumentId(doc.DocumentId)
	if err != nil {
		t.Fatal(fmt.Errorf("Failed on getDocumentByDocumentId(doc.DocumentId): %w", err))
	}
	if verify_document.CreatedAt != doc.CreatedAt {
		t.Fatalf("Failed on comparing result of getDocumentByDocumentId(doc.DocumentId)")
	}
}
