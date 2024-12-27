package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getHomePage(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte("Hello World"))
}

func (app *App) getProducts(writer http.ResponseWriter, request *http.Request) {
	response, err := getProducts(app.Db)
	if err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}
	SuccessResponse(writer, http.StatusOK, response)
}

func (app *App) getProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	product := Product{Id: id}
	err = product.getProduct(app.Db)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ErrorResponse(writer, http.StatusNotFound, "Product not found")
		default:
			ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		}
		return
	}
	SuccessResponse(writer, http.StatusOK, product)
}

func (app *App) createProduct(writer http.ResponseWriter, request *http.Request) {
	var product Product
	err := json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	err = product.createProduct(app.Db)
	if err != nil {
		log.Println(err)
		ErrorResponse(writer, http.StatusBadRequest, "Invalid Request Payload")
		return
	}
	SuccessResponse(writer, http.StatusCreated, product)
}

func (app *App) updateProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	product := Product{Id: id}
	err = json.NewDecoder(request.Body).Decode(&product)
	if err != nil {
		ErrorResponse(writer, http.StatusInternalServerError, err.Error())
	}

	err = product.updateProduct(app.Db)
	if err != nil {
		log.Println(err)
		ErrorResponse(writer, http.StatusBadRequest, err.Error())
	}
	SuccessResponse(writer, http.StatusOK, product)
}

func (app *App) deleteProduct(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		ErrorResponse(writer, http.StatusBadRequest, err.Error())
		return
	}
	product := Product{Id: id}
	err = product.deleteProduct(app.Db)
	if err != nil {
		log.Println("erS", err.Error())
		ErrorResponse(writer, http.StatusNotAcceptable, err.Error())
		return
	}
	SuccessResponse(writer, http.StatusNoContent, map[string]string{"message": "Product successfully deleted"})
}
