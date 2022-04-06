package main

import (
	"encoding/json"
	"github.com/nicourrrn/nix-shop-backend/internal/db"
	. "github.com/nicourrrn/nix-shop-backend/internal/models"
	"github.com/nicourrrn/nix-shop-backend/pkg/jwt_handler"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func PostSignUp(writer http.ResponseWriter, request *http.Request) {
	if !(request.Method == http.MethodPost) {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	req := BaseClient{}
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	clientId, err := db.Clients.NewClient(req.Name, req.Phone, req.Password)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	pair := jwt_handler.NewTokenPair(clientId, "client")
	ref, acc, err := pair.GetStrings()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	db.Clients.SetClientRefToken(clientId, ref)

	err = json.NewEncoder(writer).Encode(map[string]interface{}{
		"accessToken":  acc,
		"refreshToken": ref,
	})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostSignIn(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	req := BaseClient{}
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	id, name, err := db.Clients.GetClient(req.Phone, req.Password)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	ref, acc, err := jwt_handler.NewTokenPair(id, "client").GetStrings()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	db.Clients.SetClientRefToken(id, ref)

	err = json.NewEncoder(writer).Encode(map[string]interface{}{
		"name":         name,
		"accessToken":  acc,
		"refreshToken": ref,
	})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostRefresh(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	refRequest := Tokens{}
	err := json.NewDecoder(request.Body).Decode(&refRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	pair, err := jwt_handler.NewTokenPairFromStrings(refRequest.RefreshToken, refRequest.AccessToken)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	refToken := db.Clients.GetClientRefToken(pair.AccessToken.UserId)

	err = bcrypt.CompareHashAndPassword([]byte(refToken), []byte(refRequest.RefreshToken))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	ref, acc, err := jwt_handler.NewTokenPair(pair.AccessToken.UserId, "client").GetStrings()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	db.Clients.SetClientRefToken(pair.AccessToken.UserId, ref)

	err = json.NewEncoder(writer).Encode(map[string]string{
		"accessToken":  acc,
		"refreshToken": ref,
	})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostBasket(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	accessTokenString := request.Header.Get("Access-Token")
	accessClaim, err := jwt_handler.GetClaim(accessTokenString, jwt_handler.GetAccess())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	req := Basket{}
	err = json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	basketId := db.Clients.NewBacket(accessClaim.UserId, req)

	err = json.NewEncoder(writer).Encode(map[string]interface{}{
		"basketId": basketId,
	})

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetLogOut(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := request.Header.Get("Access-Token")
	claim, err := jwt_handler.GetClaim(token, jwt_handler.GetAccess())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	db.Clients.RemoveRefresh(claim.UserId)
}

func GetAllIngredients(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	response := make([]string, 0)
	for _, ingr := range db.Products.GetAllIngredients() {
		response = append(response, ingr.Name)
	}
	err := json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetSuppliers(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	suppliers := db.Suppliers.GetSuppliers()
	err := json.NewEncoder(writer).Encode(suppliers)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInsufficientStorage)
		return
	}
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := request.Header.Get("Access-Token")
	claim, err := jwt_handler.GetClaim(token, jwt_handler.GetAccess())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := db.Clients.ClientData(claim.UserId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetSupplierMenu(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	supplierId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(writer).Encode(db.Products.GetProducts(int64(supplierId)))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllBasket(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	token := request.Header.Get("Access-Token")
	claim, err := jwt_handler.GetClaim(token, jwt_handler.GetAccess())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	baskets, err := db.Clients.AllBasket(claim.UserId)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	if baskets == nil {
		baskets = make([]SavedBasket, 0)
	}
	err = json.NewEncoder(writer).Encode(baskets)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
