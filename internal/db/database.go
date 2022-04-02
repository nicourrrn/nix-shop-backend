package db

import (
	. "backend/internal/models"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
)

var Clients clientRepo
var Products productRepo
var Suppliers supplierRepo

func init() {
    jawsdb := os.Getenv("JAWSDB_URL")
	username := os.Getenv("sql_username")
	password := os.Getenv("sql_password")
	database_name := os.Getenv("sql_database")
    if jawsdb != ""{
		tmp = jawsdb

	} else {
		tmp := fmt.Sprintf("%s:%s@/%s", username, password, database_name)
	}
	db, err := sqlx.Open("mysql", tmp)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	Clients = clientRepo{connection: db}
	Products = productRepo{connection: db}
	Suppliers = supplierRepo{connection: db}
}

func FindTypeId(types []Type, name string) int64 {
	for _, t := range types {
		if t.Name == name {
			return t.Id
		}
	}
	return -1
}
