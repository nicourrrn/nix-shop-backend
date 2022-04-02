package db

import (
	"github.com/jmoiron/sqlx"
	. "github.com/nicourrrn/nix-shop-backend/internal/models"
)

type supplierRepo struct {
	connection *sqlx.DB
}

func (repo *supplierRepo) GetSuppliers() (suppliers []Supplier) {
	err := repo.connection.Select(&suppliers, "SELECT suppliers.id, suppliers.name, suppliers.image, st.name as type, suppliers.open_at, suppliers.close_at FROM suppliers JOIN supplier_types st on st.id = suppliers.type_id")
	if err != nil {
		panic(err)
	}
	return
}

func (repo *supplierRepo) AddSupplier(supplier Supplier) int64 {
	typeId := FindTypeId(repo.GetTypes(), supplier.Type)
	if typeId == -1 {
		typeId = repo.AddType(supplier.Type)
	}
	result, err := repo.connection.Exec(
		"INSERT INTO suppliers(name, type_id, image, open_at, close_at) VALUE (?, ?, ?, ?, ?)",
		supplier.Name, typeId, supplier.Image, supplier.OpenAt, supplier.CloseAt)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	if supplier.Id != 0 {
		_, err = repo.connection.Exec("UPDATE suppliers SET id = ? WHERE id = ?", supplier.Id, id)
		if err != nil {
			panic(err)
		}
		return supplier.Id
	}

	return id
}

func (repo *supplierRepo) GetTypes() (types []Type) {
	err := repo.connection.Select(&types, "SELECT id, name FROM supplier_types")
	if err != nil {
		panic(err)
	}
	return
}

func (repo *supplierRepo) AddType(name string) int64 {
	result, err := repo.connection.Exec("INSERT INTO supplier_types(name) VALUE (?)", name)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}
