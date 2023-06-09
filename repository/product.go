package repository

import (
	"database/sql"
	"go-api-meli/model"
)

type products struct {
	db *sql.DB
}

func RepositoryProduct(db *sql.DB) *products {
	return &products{db}
}

func (products products) CreateProduct(product model.Product) (uint64, error) {
	statement, err := products.db.Prepare(
		"insert into tb_product (title, price, quantity) values (?,?,?)",
	)
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	result, err := statement.Exec(product.Title, product.Price, product.Quantity)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ID), nil

}
func (product products) GetAll() ([]model.Product, error) {
	rows, err := product.db.Query("select * from tb_product")
	if err != nil {
		return nil, err

	}

	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var product model.Product

		if err = rows.Scan(&product.ID, &product.Title, &product.Price, &product.Quantity); err != nil {
			return nil, err
		}

		products = append(products, product)
	}
	return products, nil

}
func (product products) GetById(ID uint64) (model.Product, error) {
	row, err := product.db.Query("select idtb_product, title, price, quantity from tb_product where idtb_product = ?", ID)
	if err != nil {
		return model.Product{}, err
	}
	defer row.Close()

	var prd model.Product

	if row.Next() {
		if err = row.Scan(
			&prd.ID,
			&prd.Title,
			&prd.Price,
			&prd.Quantity,
		); err != nil {
			return model.Product{}, err

		}
	}
	return prd, nil
}
