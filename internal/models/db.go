package models

type Type struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}

type Product struct {
	Id          int64    `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Price       float32  `json:"price" db:"price"`
	Image       string   `json:"image" db:"image"`
	Type        string   `json:"type" db:"type"`
	Ingredients []string `json:"ingredients" db:"ingredients"`
}

type Supplier struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Image   string `db:"image" json:"image"`
	Type    string `db:"type" json:"type"`
	OpenAt  string `db:"open_at" json:"openAt"`
	CloseAt string `db:"close_at" json:"closeAt"`
}

type SavedBasket struct {
	Id       int64     `db:"id" json:"id"`
	Date     string    `db:"create_at" json:"date"`
	Products []Product `json:"products"`
	Address  string    `db:"address" json:"address"`
	Price    float32   `db:"price" json:"price"`
}
