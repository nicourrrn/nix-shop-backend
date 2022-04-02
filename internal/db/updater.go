package db

import (
	"encoding/json"
	. "github.com/nicourrrn/nix-shop-backend/internal/models"
	"log"
	"os"
)

func NewFromJson(filename string) (Shop, error) {
	open, err := os.Open(filename)
	if err != nil {
		return Shop{}, err
	}
	var shop Shop
	err = json.NewDecoder(open).Decode(&shop)
	if err != nil {
		return Shop{}, err
	}
	return shop, nil
}

func WriteShop(shop Shop) {
	Suppliers.AddSupplier(Supplier{
		Id:      int64(shop.Id),
		Name:    shop.Name,
		Image:   shop.Image,
		Type:    shop.Type,
		OpenAt:  shop.WorkingHours.Opening,
		CloseAt: shop.WorkingHours.Closing,
	})

	for _, prod := range shop.Menu {
		Products.AddProduct(Product{
			Id:    prod.Id,
			Name:  prod.Name,
			Price: float32(prod.Price),
			Image: prod.Image,
			Type:  prod.Type,
		}, int64(shop.Id))

		Products.ConnProdWithIngr(prod.Id, prod.Ingredients)
	}
}

func WriteData() {
	path := "internal/db/shops"
	dir, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatalln(err)
	}
	for _, i := range files {
		shop, err := NewFromJson(path + "/" + i.Name())
		if err != nil {
			log.Println(err)
		}
		WriteShop(shop)
	}
}
