package main

import (
	"fmt"
	"net/http"
	"personal-web/connection"
	"personal-web/middleware"
	"personal-web/routes"

	"github.com/gorilla/mux"
)

// main function
func main() {

	//declaration new routeras r
	r := mux.NewRouter()

	//connect to db
	connection.DatabaseConnect()

	//creating static folder
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	//auth url handler
	r.HandleFunc("/register", routes.FormRegister).Methods("GET")
	r.HandleFunc("/register", routes.Register).Methods("POST")
	r.HandleFunc("/login", routes.FormLogin).Methods("GET")
	r.HandleFunc("/login", routes.Login).Methods("POST")

	//url handler
	r.HandleFunc("/", routes.Home).Methods("GET")
	r.HandleFunc("/add-project", routes.AddProject).Methods("GET")
	r.HandleFunc("/add-project", middleware.UploadFile(routes.StoreProject)).Methods("POST")
	r.HandleFunc("/edit-project/{id}", routes.FormEditProject).Methods("GET")
	r.HandleFunc("/edit-project/{id}", middleware.UploadFile(routes.StoreEdit)).Methods("POST")
	r.HandleFunc("/detail-project/{id}", routes.ProjectDetail).Methods("GET")
	r.HandleFunc("/delete-project/{id}", routes.DeleteProject).Methods("GET")
	r.HandleFunc("/contact-me", routes.ContactMe).Methods("GET")

	//creat local server on port 3000
	fmt.Println("Server running on port 3000")
	http.ListenAndServe(":3000", r)
}
