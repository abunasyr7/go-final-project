package main

import (
	"log"

	"github.com/abunasyr7/go-final-project/pkg/server"
)

func main () {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}