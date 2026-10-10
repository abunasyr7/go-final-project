package main

import (
	"log"
	"os"

	"github.com/abunasyr7/go-final-project/pkg/db"
	"github.com/abunasyr7/go-final-project/pkg/server"
)

const defaultDBFile = "scheduler.db"

func main () {
	dbFile := defaultDBFile

	if env := os.Getenv("TODO_DBFILE"); env != "" {
		dbFile = env
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}