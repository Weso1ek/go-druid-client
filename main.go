package main

import (
	"go-druid-client/httpcontroller"
	"log"
	"net/http"
)

func main() {
	druidController := httpcontroller.NewDruidController()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/stat/dau", druidController.StatDau)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux.ServeHTTP(w, r)
	})

	log.Fatal(http.ListenAndServe(":8099", handler))
}
