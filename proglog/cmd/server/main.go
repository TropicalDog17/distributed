package main

import (
	"log"

	"github.com/TropicalDog17/distributed/proglog/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
