package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/igorsfreitas/devbook-api/src/router"
)

func main() {
	fmt.Println("Rodando api")
	r := router.Gerar()

	log.Fatal(http.ListenAndServe(":5000", r))
}
