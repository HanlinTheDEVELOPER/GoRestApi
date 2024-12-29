package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	app.Initialize(USERNAME, PASSWORD, "test")
	createTable()
	m.Run()
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
	insertQuery := fmt.Sprintf("INSERT INTO inventory(name, price, quantity) VALUES('%v', %v, %v);", name, price, quantity)
	_, err := app.Db.Exec(insertQuery)
	if err != nil {
		log.Fatal(err, "err")
	}
}

func TestGetAllProducts(t *testing.T) {
	clearTable()
	addProduct("keyboard", 50, 5)
	addProduct("mouse", 25, 10)
	addProduct("monitor", 100, 2)
	request, _ := http.NewRequest("GET", "/product", nil)
	response := sendRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

	var data []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &data)
	if len(data) != 3 {
		t.Errorf("Expected 3 products, got %d", len(data))
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 50, 5)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var jsonStr = []byte(`{"name":"keyboard", "price":50, "quantity":5}`)
	request, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response := sendRequest(request)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "keyboard" {
		t.Errorf("Expected product name to be 'keyboard'. Got %v", m["name"])
	}
	if m["price"] != 50.0 {
		t.Errorf("Expected product price to be '50'. Got %v", m["price"])
	}
	if m["quantity"] != 5.0 {
		t.Errorf("Expected product quantity to be '5'. Got %v", m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 50, 5)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(request)
	checkResponseCode(t, http.StatusNoContent, response.Code)

	request, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(request)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 5, 50)

	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var jsonStr = []byte(`{"name":"keyboard", "price":50, "quantity":50}`)

	request, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")
	response = sendRequest(request)
	checkResponseCode(t, http.StatusOK, response.Code)

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)
	log.Println(oldValue, newValue)
	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected the same id")
	}
	if oldValue["name"] != newValue["name"] {
		t.Errorf("Expected the same name")
	}
	if oldValue["price"] != newValue["price"] {
		t.Errorf("Expected the same price")
	}
	if oldValue["quantity"] == newValue["quantity"] {
		t.Errorf("Expected the quantity to change")
	}
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
