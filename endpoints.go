package main

import (
	"fmt"
	"logbook/main/controllers"

	"net/http"

	"github.com/gorilla/mux"
)

type Endpoint struct {
	Path    string
	Method  string
	Handler func(http.ResponseWriter, *http.Request)
}

func rootURIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<b>Forbidden</b><hr>Logbook API is not provided for 3rd party use.<br>
	Go to main website and create an account. 
	There is no guarantee for a future release won't break your code.`)
	w.Header().Set("Content-Type", "text/html")
}

func registerEndpoints(r *mux.Router) {
	// var (
	// 	documentOverviewHierarchical  = document.CDocumentOverviewHierarchical{}
	// 	documentOverviewChronological = document.CDocumentOverviewChronological{}
	// )

	r.HandleFunc("/", rootURIHandler).Methods("GET")

	// r.HandleFunc("/auth/user", authUserCreate).
	// 	Methods("POST").
	// 	Queries("email_address", "{email_address}")

	// r.HandleFunc("/auth/session", authSessionCreate).
	// 	Methods("POST").
	// 	Queries("email_address", "{email_address}").
	// 	Queries("password_kdf", "{password_kdf}")

	r.HandleFunc("/account", controllers.UserCreate).Methods("POST")

	r.HandleFunc("/task", controllers.TaskCreate).Methods("POST")
	// r.HandleFunc("/task/super", controllers.TaskPatchSuper.Methods("PATCH"))
	// r.HandleFunc("/task/content", controllers.TaskPatchContent).Methods("PATCH")

	r.HandleFunc("/placement-array/hiearchical", controllers.PlacementArrayHierarchical).Methods("POST")
	// r.HandleFunc("/placement-array/chronological", controllers.PlacementArrayChronological).Methods("POST")

}
