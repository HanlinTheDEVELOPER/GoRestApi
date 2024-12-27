package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]Product, error) {
	query := "SELECT * FROM inventory"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (p *Product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT * FROM inventory WHERE id=%v", p.Id)
	row := db.QueryRow(query)
	err := row.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO inventory(name, quantity, price) VALUES ('%v', '%v', '%v')", p.Name, p.Quantity, p.Price)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("UPDATE inventory SET name='%v', quantity='%v', price='%v' WHERE id='%v'", p.Name, p.Quantity, p.Price, p.Id)

	result, err := db.Exec(query)
	rowEffected, _ := result.RowsAffected()
	if rowEffected == 0 || err != nil {
		return errors.New("Product does not exist")
	}
	return nil
}

func (p *Product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("DELETE FROM inventory WHERE id='%v'", p.Id)
	_, err := db.Exec(query)
	if err != nil {
		log.Println("er", err)
		return err
	}
	return nil
}
