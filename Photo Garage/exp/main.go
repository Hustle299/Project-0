package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "duytung299"
	dbname   = "photogarage_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	var id int
	for i := 1; i < 6; i++ {
		// Create some fake data
		userId := 1
		if i > 3 {
			userId = 2
		}
		amount := 1000 * i
		description := fmt.Sprintf("USB-C Adapter x%d", i)
		err = db.QueryRow(`
INSERT INTO orders (user_id, amount, description)
VALUES ($1, $2, $3)
RETURNING id`,
			userId, amount, description).Scan(&id)
		if err != nil {
			panic(err)
		}
		fmt.Println("Created an order with the ID:", id)
	}
	db.Close()
}
