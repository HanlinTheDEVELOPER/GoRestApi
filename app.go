package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	Db     *sql.DB
}

func (app *App) Initialize(Username string, Password string, DbName string) {
	connectionString := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", Username, Password, HOST, DBPORT, DbName)
	var err error
	app.Db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
}

func (app *App) Run() {
	app.handleRoutes()
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("localhost:%v", APPPORT), app.Router))
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/", getHomePage).Methods("GET")
	app.Router.HandleFunc("/product", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")
}
