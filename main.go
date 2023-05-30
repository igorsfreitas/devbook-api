package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/igorsfreitas/devbook-api/src/config"
	"github.com/igorsfreitas/devbook-api/src/router"
)

func main() {
	config.Load()

	fmt.Printf("Rodando api na porta %d", config.Port)
	r := router.Gerar()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
