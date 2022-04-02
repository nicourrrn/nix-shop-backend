package db

import (
	"github.com/jmoiron/sqlx"
	. "github.com/nicourrrn/nix-shop-backend/internal/models"
)

type productRepo struct {
	connection *sqlx.DB
}

func (repo *productRepo) GetAllIngredients() (ingredients []Type) {
	repo.connection.Select(&ingredients, "SELECT id, name FROM ingredients")
	return
}

func (repo *productRepo) AddIngredients(newIngredients []string) {
	ingredients := make([]map[string]interface{}, 0)
	for _, i := range newIngredients {
		ingredients = append(ingredients, map[string]interface{}{"name": i})
	}
	_, err := repo.connection.NamedExec("INSERT INTO ingredients(name) VALUES (:name) ON DUPLICATE KEY UPDATE name=name", ingredients)
	if err != nil {
		panic(err)
	}
}

func (repo *productRepo) ConnProdWithIngr(productId int64, ingredients []string) {
	repo.AddIngredients(ingredients)
	allIngredients := repo.GetAllIngredients()
	toDB := make([]map[string]interface{}, 0)
	for _, ingr := range ingredients {
		ingrId := FindTypeId(allIngredients, ingr)
		toDB = append(toDB, map[string]interface{}{
			"product_id":    productId,
			"ingredient_id": ingrId,
		})
	}
	_, err := repo.connection.NamedExec(
		"INSERT INTO product_ingredient(product_id, ingredient_id) VALUES (:product_id, :ingredient_id) ON DUPLICATE KEY UPDATE product_id=product_id",
		toDB)
	if err != nil {
		panic(err)
	}
}

func (repo *productRepo) GetProducts(supplierId int64) (menu []Product) {
	err := repo.connection.Select(&menu,
		"SELECT products.id, products.name, products.price, products.image, pt.name as type FROM products JOIN product_types pt on pt.id = products.type_id WHERE supplier_id = ?", supplierId)
	if err != nil {
		panic(err)
	}
	for i, p := range menu {
		rows, err := repo.connection.Queryx(
			"SELECT name FROM product_ingredient JOIN ingredients i on i.id = product_ingredient.ingredient_id WHERE product_ingredient.product_id = ?",
			p.Id)
		if err != nil {
			panic(err)
		}
		var temp string
		for rows.Next() {
			err = rows.Scan(&temp)
			if err != nil {
				panic(err)
			}
			menu[i].Ingredients = append(menu[i].Ingredients, temp)
		}
	}
	return
}

func (repo *productRepo) AddProduct(product Product, sId int64) int64 {
	typeId := FindTypeId(repo.GetTypes(), product.Type)
	if typeId == -1 {
		typeId = repo.AddType(product.Type)
	}
	result, err := repo.connection.Exec(
		"INSERT INTO products(name, type_id, image, supplier_id, price) VALUE (?, ?, ?, ?, ?)",
		product.Name, typeId, product.Image, sId, product.Price)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	if product.Id != 0 {
		_, err = repo.connection.Exec("UPDATE products SET id = ? WHERE id = ?", product.Id, id)
		if err != nil {
			panic(err)
		}
		return product.Id
	}

	return id
}

func (repo *productRepo) GetTypes() (types []Type) {
	err := repo.connection.Select(&types, "SELECT id, name FROM product_types")
	if err != nil {
		panic(err)
	}
	return
}

func (repo *productRepo) AddType(name string) int64 {
	result, err := repo.connection.Exec("INSERT INTO product_types(name) VALUE (?)", name)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}
