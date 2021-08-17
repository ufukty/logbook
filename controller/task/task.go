package task

import (
	"encoding/json"
	"fmt"
	"log"
	"logbook/main/controller/task/linearized_tasks"
	"net/http"
	"os"
)

func Controller(w http.ResponseWriter, r *http.Request) {

	fmt.Println(w)
	fmt.Println()

	file, _ := os.ReadFile("data/data.json")

	document := []linearized_tasks.TaskModel{}

	err := json.Unmarshal(file, &document)
	if err != nil {
		log.Println("Error: ", err)
	}

	linearized_tasks_ := linearized_tasks.DFS(document)

	w.Header().Set("Content-Type", "application/json")
	// pretty_print, _ := json.MarshalIndent(document, "", "\t")
	json.NewEncoder(w).Encode(linearized_tasks_)
	log.Println("Request processed.")

}
