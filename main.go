package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fmt.Println(mux)
}
