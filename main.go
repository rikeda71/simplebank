package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/s14t284/simplebank/api"
	db "github.com/s14t284/simplebank/db/ent"
	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	entClient, err := ent.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("failed opening connection to sqllite: %v", err)
	}
	defer entClient.Close()

	store := db.NewStore(entClient)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
}
