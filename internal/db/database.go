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
	db := Connect()
	if os.Getenv("FIRST") == "true" {
		Migration(db)
		defer WriteData()
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

func Migration(connection *sqlx.DB) {
	var queries [9]string
	queries[0] = "CREATE TABLE IF NOT EXISTS `supplier_types`( `id` int(11)  auto_increment NOT NULL , `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `update_at` timestamp NOT NULL DEFAULT current_timestamp(), `create_at` timestamp NOT NULL DEFAULT current_timestamp(), PRIMARY KEY (`id`), UNIQUE KEY `supplier_types_name_uindex` (`name`));"
	queries[1] = "CREATE TABLE IF NOT EXISTS `suppliers`(`id` int(11)  auto_increment NOT NULL, `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `type_id` int(11) NOT NULL, `image` text COLLATE utf8mb4_unicode_ci NOT NULL,`open_at` varchar(5) COLLATE utf8mb4_unicode_ci NOT NULL, `close_at`  varchar(5) COLLATE utf8mb4_unicode_ci NOT NULL, `update_at` timestamp NOT NULL DEFAULT current_timestamp(), `create_at` timestamp NOT NULL DEFAULT current_timestamp(), PRIMARY KEY (`id`), KEY `suppliers_supplier_types_id_fk` (`type_id`), CONSTRAINT `suppliers_supplier_types_id_fk` FOREIGN KEY (`type_id`) REFERENCES `supplier_types` (`id`) ON DELETE CASCADE);"
	queries[2] = "CREATE TABLE IF NOT EXISTS `clients`(`id` int(11)  auto_increment NOT NULL, `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `phone` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `refresh_token` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL, `update_at` timestamp NOT NULL DEFAULT current_timestamp(), `create_at` timestamp NOT NULL DEFAULT current_timestamp(), `password` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `clients_email_uindex` (`phone`));"
	queries[3] = "CREATE TABLE IF NOT EXISTS `baskets`(`id` int(11)  auto_increment NOT NULL, `client_id` int(11) NOT NULL, `address` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `create_at` timestamp NOT NULL DEFAULT current_timestamp(), `update_at` timestamp NOT NULL DEFAULT current_timestamp(), `price` float NOT NULL, PRIMARY KEY (`id`), KEY `baskets_clients_id_fk` (`client_id`), CONSTRAINT `baskets_clients_id_fk` FOREIGN KEY (`client_id`) REFERENCES `clients` (`id`) ON DELETE CASCADE);"
	queries[4] = "CREATE TABLE IF NOT EXISTS `ingredients`(`id` int(11) auto_increment NOT NULL,`name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `update_at` timestamp NOT NULL DEFAULT current_timestamp(), `create_at` timestamp NOT NULL DEFAULT current_timestamp(), PRIMARY KEY (`id`), UNIQUE KEY `ingredients_name_uindex` (`name`));"
	queries[5] = "CREATE TABLE IF NOT EXISTS `product_types`( `id` int(11)  auto_increment NOT NULL, `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `update_at` timestamp NOT NULL DEFAULT current_timestamp(), `create_at` timestamp NOT NULL DEFAULT current_timestamp(), PRIMARY KEY (`id`), UNIQUE KEY `product_types_name_uindex` (`name`));"
	queries[6] = "CREATE TABLE IF NOT EXISTS `product_ingredient` ( `product_id` int(11) NOT NULL, `ingredient_id` int(11) NOT NULL, `update_at` timestamp NOT NULL DEFAULT current_timestamp(), `create_at` timestamp NOT NULL DEFAULT current_timestamp());"
	queries[7] = "CREATE TABLE IF NOT EXISTS `products`(`id` int(11)  auto_increment NOT NULL, `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL, `price` float NOT NULL, `image` text COLLATE utf8mb4_unicode_ci NOT NULL, `type_id` int(11) NOT NULL, `supplier_id` int(11) NOT NULL, `create_at` timestamp NOT NULL DEFAULT current_timestamp(), `update_at` timestamp NOT NULL DEFAULT current_timestamp(), PRIMARY KEY (`id`), KEY `products_product_types_id_fk` (`type_id`), KEY `products_suppliers_id_fk` (`supplier_id`), CONSTRAINT `products_product_types_id_fk` FOREIGN KEY (`type_id`) REFERENCES `product_types` (`id`) ON DELETE CASCADE, CONSTRAINT `products_suppliers_id_fk` FOREIGN KEY (`supplier_id`) REFERENCES `suppliers` (`id`) ON DELETE CASCADE);"
	queries[8] = "CREATE TABLE IF NOT EXISTS `product_basket`( `product_id` int(11)  NOT NULL, `basket_id` int(11) NOT NULL, `count` int(11) NOT NULL, `create_at` timestamp NOT NULL DEFAULT current_timestamp(), `update_at` timestamp NOT NULL DEFAULT current_timestamp());"
	for _, i := range queries {
		_, err := connection.Exec(i)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func Connect() *sqlx.DB {
	jawsdb := os.Getenv("JAWSDB_URL")
	username := os.Getenv("sql_username")
	password := os.Getenv("sql_password")
	database_name := os.Getenv("sql_database")
	var tmp string
	if jawsdb != "" {
		tmp = jawsdb
	} else {
		tmp = fmt.Sprintf("%s:%s@/%s", username, password, database_name)
	}
	db, err := sqlx.Open("mysql", tmp)
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
