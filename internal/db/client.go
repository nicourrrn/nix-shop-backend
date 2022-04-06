package db

import (
	"errors"
	"github.com/jmoiron/sqlx"
	. "github.com/nicourrrn/nix-shop-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type clientRepo struct {
	connection *sqlx.DB
}

func (repo *clientRepo) NewClient(name, phone, password string) (int64, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, errors.New(BcryptError)
	}
	scannedClient := BaseClient{
		Name:     name,
		Phone:    phone,
		Password: string(encryptedPassword),
	}
	result, err := repo.connection.NamedExec(
		"INSERT INTO clients(name, phone, password) VALUE (:name, :phone, :password)",
		scannedClient)
	if err != nil {
		return -1, errors.New(TakenClient)
	}
	return result.LastInsertId()
}

func (repo *clientRepo) GetClient(phone string, password string) (id int64, name string, err error) {
	row := repo.connection.QueryRow("SELECT id, name, password FROM clients WHERE phone = ?", phone)
	var savedPassword string
	err = row.Scan(&id, &name, &savedPassword)
	if err != nil {
		err = errors.New(InvalidClient)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password))
	if err != nil {
		err = errors.New(InvalidClient)
		return
	}
	return
}

func (repo *clientRepo) SetClientRefToken(id int64, refreshToken string) {
	encryptedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = repo.connection.Exec(
		"UPDATE clients SET refresh_token = ? WHERE id = ?",
		string(encryptedToken), id)
	if err != nil {
		log.Fatalln(err)
	}
}

func (repo *clientRepo) GetClientRefToken(id int64) (refToken string) {
	err := repo.connection.Get(&refToken, "SELECT refresh_token FROM clients WHERE id = ?", id)
	if err != nil {
		panic(err)
	}
	return
}

func (repo *clientRepo) NewBacket(clientId int64, data Basket) int64 {
	var finalPrice float32
	for _, p := range data.Products {
		finalPrice += p.PriceOne * float32(p.Count)
	}
	result, err := repo.connection.Exec("INSERT INTO baskets(client_id, price, address) VALUE (?, ?, ?)", clientId, finalPrice, data.Address)
	if err != nil {
		panic(err)
	}

	basketId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	for _, p := range data.Products {
		_, err = repo.connection.Exec("INSERT INTO product_basket(product_id, basket_id, count) VALUE (?, ?, ?)", p.ProductId, basketId, p.Count)
		if err != nil {
			panic(err)
		}
	}
	return basketId
}

func (repo *clientRepo) RemoveRefresh(clientId int64) {
	_, err := repo.connection.Exec("UPDATE clients SET refresh_token = '' WHERE id = ?", clientId)
	if err != nil {
		panic(err)
	}
}

func (repo *clientRepo) AllBasket(clientId int64) (baskets []SavedBasket, err error) {
	err = repo.connection.Select(&baskets, "SELECT id, address, price, create_at FROM baskets WHERE client_id = ?", clientId)
	if err != nil {
		return
	}
	for i := range baskets {
		err = repo.connection.Select(&baskets[i].Products, "SELECT p.id, p.name FROM product_basket as pb JOIN products p on p.id = pb.product_id WHERE pb.basket_id = ?", baskets[i].Id)
		if err != nil {
			return
		}
	}
	return
}
