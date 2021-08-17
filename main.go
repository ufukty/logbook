package main

import (
	"log"
	"logbook/main/controller/task"
	"net/http"
)

func main() {

	http.HandleFunc("/", task.Controller)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
