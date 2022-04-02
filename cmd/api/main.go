package main

import (
	"github.com/justinas/alice"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	api := http.NewServeMux()
	headersMiddleware := alice.New(AddHeaders)
	api.Handle("/ingredients", headersMiddleware.ThenFunc(GetAllIngredients))
	api.Handle("/user/signin", headersMiddleware.ThenFunc(PostSignIn))
	api.Handle("/user/signup", headersMiddleware.ThenFunc(PostSignUp))
	api.Handle("/user/refresh", headersMiddleware.ThenFunc(PostRefresh))
	api.Handle("/user/logout", headersMiddleware.ThenFunc(GetLogOut))
	api.Handle("/suppliers", headersMiddleware.ThenFunc(GetSuppliers))
	api.Handle("/products", headersMiddleware.ThenFunc(GetSupplierMenu))
	api.Handle("/basket/new", headersMiddleware.ThenFunc(PostBasket))
	api.Handle("/basket/all", headersMiddleware.ThenFunc(GetAllBasket))
	log.Println(http.ListenAndServe(":"+port, api))

}

func AddHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, withCredentials, access-token, user-agent")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if request.Method != http.MethodOptions {
			next.ServeHTTP(writer, request)
		}
	})
}
