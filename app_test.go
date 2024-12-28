package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	app.Initialize(USERNAME, PASSWORD, "test")
	m.Run()
	createTable()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS inventory(
		id INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		price FLOAT(10,7) ,
		quantity INT,
		PRIMARY KEY (id)
	);`
	_, err := app.Db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.Db.Exec("DELETE FROM inventory")
	app.Db.Exec("ALTER TABLE inventory AUTO_INCREMENT = 1")
}

func addProduct(name string, price float64, quantity int) {
	insertQuery := fmt.Sprintf("INSERT INTO inventory(name, price, quantity) VALUES('%v', %v, %v)", name, price, quantity)
	_, err := app.Db.Exec(insertQuery, name, price, quantity)
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetProducts(t *testing.T) {
	clearTable()
	addProduct("keyboard", 50, 5)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, request)
	return recorder
}
