package database

import (
	"fmt"
	"testing"
)

func TestPostgresErrorsStripSQLState(t *testing.T) {
	initRegex()
	sqlError := `ERROR: duplicate key value violates unique constraint "USER_email_address_key" (SQLSTATE 23505)`
	expectedResult := "23505"
	if expectedResult != StripSQLState(sqlError) {
		t.Fail()
	}
}

func TestCreateUser(t *testing.T) {

	// Initialize database connection with test database
	Init("postgres://ufuktan:password@localhost:5432/logbook_dev")
	defer CloseDatabaseConnections()

	// Create USER
	myUser := User{
		NameSurname:    string("Name Surname"),
		EmailAddress:   string("testUserCreate@golang.example.com"),
		Salt:           string("loremipsum"),
		HashedPassword: string("$argon2id$v=19$m=4096,t=3,p=1$bG9yZW1pcHN1bQ$6ASAMXM/1Czod3ixSuEe6x+nb96mFkTWjlruH+fAGtY"),
	}
	result := Db.Create(&myUser)
	if result.Error != nil {
		fmt.Println()
	}

	// var (
	// 	myUser *User
	// 	myDoc *Document
	// )

	// myUser, err := UserCreate(
	// 	"Name Surname",
	// 	"test.test@test.tld",
	// 	"testtest@test.tld",
	// 	"$argon2id$v=19$m=4096,t=3,p=1$bG9yZW1pcHN1bQ$6ASAMXM/1Czod3ixSuEe6x+nb96mFkTWjlruH+fAGtY",
	// )
	// if err != nil {
	// 	t.Fatalf("Failed on UserCreate() : %s", err)
	// }

	// // Create DOCUMENT

	// myDoc, err := DocumentCreate(myUser.UserId)
	// if err != nil {
	// 	t.Fatalf("Failed on DocumentCreate() : %s", err)
	// }

	// // Create USER

	// test_modified_tasks, err := CreateTask(Task{
	// 	DocumentId: myDoc.DocumentId,
	// 	Content:    "One difficult task for testing go package.",
	// 	ParentId:   "00000000-0000-0000-0000-000000000000",
	// })
	// if err != nil {
	// 	t.Fatalf("Failed on CreateTask() : %s", err)
	// }
	// test_task := test_modified_tasks[0]

	// verify_task, err := GetTaskByTaskId(test_task.TaskId)
	// if err != nil {
	// 	t.Fatalf("Failed on GetTaskByTaskId() : %s", err)
	// }
	// if verify_task.Content != test_task.Content {
	// 	t.Fatalf("Failed on comparing result of GetTaskByTaskId(test_task.TaskId)")
	// }

	// myDocGot, err := DocumentGet(myUser.UserId, myDoc.DocumentId)
	// if err != nil {
	// 	t.Fatalf("Failed on DocumentGet(test_document.DocumentId)\nError details: %s", err)
	// }
	// if myDocGot.CreatedAt != myDoc.CreatedAt {
	// 	t.Fatalf("Failed on comparing result of DocumentGet(test_document.DocumentId)")
	// }
}
