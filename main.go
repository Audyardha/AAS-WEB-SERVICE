package main

import (
	"audy/controller/auth"
	"audy/controller/karya"
	"audy/controller/pelukis"
	"audy/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	fmt.Println("Hello World")

	router := mux.NewRouter()

	//Auth routes
	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	//Pelukis routes
	router.HandleFunc("/pelukis", pelukis.GetPelukis).Methods("GET")
	router.HandleFunc("/pelukis", pelukis.PostPelukis).Methods("POST")
	router.HandleFunc("/pelukis/{id}", auth.JWTAuth(pelukis.PutPelukis)).Methods("PUT")
	router.HandleFunc("/pelukis/{id}", auth.JWTAuth(pelukis.DeletePelukis)).Methods("DELETE")

	//Karya routes
	router.HandleFunc("/karya", karya.GetKarya).Methods("GET")
	router.HandleFunc("/karya", auth.JWTAuth(karya.PostKarya)).Methods("POST")
	router.HandleFunc("/karya/{id}", auth.JWTAuth(karya.PutKarya)).Methods("PUT")
	router.HandleFunc("/karya/{id}", auth.JWTAuth(karya.DeleteKarya)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug:          true,
	})

	handler := c.Handler(router)

	fmt.Println("Server is running on http://localhost:8020")
	log.Fatal(http.ListenAndServe(":8020", handler))
}
