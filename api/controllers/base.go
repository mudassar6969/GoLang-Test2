package controllers

import (
	"assignment2/api/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize() {
	var err error
	connStr := "user=mudassarraza dbname=testdb sslmode=disabled"
	a.DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to database")
	}

	a.DB.Debug().AutoMigrate(&models.User{})
	a.DB.Debug().AutoMigrate(&models.Task{})

	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {

	a.Router.HandleFunc("/register", a.RegisterUser).Methods("POST")
	a.Router.HandleFunc("/login", a.LoginUser).Methods("POST")
	a.Router.HandleFunc("/task", a.AddTask).Methods("POST")
	a.Router.HandleFunc("/task", a.GetTasks).Methods("GET")
	a.Router.HandleFunc("/task/{id}", a.DeleteTask).Methods("DELETE")
	a.Router.HandleFunc("/task/{id}", a.UpdateTask).Methods("PUT")
}

func (a *App) RunServer() {
	http.ListenAndServe(":5000", a.Router)
}
