package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go-druid-client/httpcontroller"
	"log"
	"net/http"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	druidController := httpcontroller.NewDruidController()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/stat/dau", druidController.StatDau)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})

	fmt.Println("Server Start")
	log.Fatal(http.ListenAndServe(":8099", handler))
}
