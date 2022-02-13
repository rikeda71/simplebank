package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/s14t284/simplebank/api"
	db "github.com/s14t284/simplebank/db/ent"
	"github.com/s14t284/simplebank/ent"
)

const (
	dbDriver      = "postgres"
	dbSource      = "host=localhost port=5432 user=root dbname=simple_bank password=secret sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	entClient, err := ent.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("failed opening connection to sqllite: %v", err)
	}
	defer entClient.Close()

	store := db.NewStore(entClient)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
}
