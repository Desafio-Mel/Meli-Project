package controller

import (
	"encoding/json"
	"fmt"
	"go-api-meli/database"
	"go-api-meli/model"
	"go-api-meli/repository"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var product model.Product

	if err = json.Unmarshal(request, &product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	repository := repository.RepositoryProduct(db)
	ID, err := repository.CreateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Product id %d", ID)))

}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	repository := repository.RepositoryProduct(db)
	product, err := repository.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)

}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, err := strconv.ParseUint(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	repository := repository.RepositoryProduct(db)
	product, err := repository.GetById(ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, _ := strconv.ParseInt(paramters["productID"], 10, 32)
	request, _ := ioutil.ReadAll(r.Body)

	var product model.Product

	json.Unmarshal(request, &product)

	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update tb_product set title = ?, price = ?, quantity = ? where idtb_product = ? ")

	defer statement.Close()
	statement.Exec(&product.Title, &product.Price, &product.Quantity, ID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Product updated successfully.")))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, _ := strconv.ParseInt(paramters["productID"], 10, 32)

	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()
	statement, err := db.Prepare("delete from tb_product where idtb_product = ? ")

	defer statement.Close()
	statement.Exec(ID)

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Product deleted successfully.")))

}
