package models

type BaseClient struct {
	Name     string `json:"name" db:"name"`
	Phone    string `json:"phone" db:"phone"`
	Password string `json:"password" db:"password"`
	CreateAt string `json:"create_at" db:"create_at"`
}

type Basket struct {
	Address  string          `json:"address"`
	Products []ProductToBask `json:"products"`
}

type ProductToBask struct {
	Count     int     `json:"count"`
	ProductId int     `json:"productId"`
	PriceOne  float32 `json:"priceOne"`
}
