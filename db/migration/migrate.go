package main

import (
	"context"
	"fmt"
	"log"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/migrate"

	_ "github.com/lib/pq"
)

func main() {
	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		"localhost", "5432", "root", "simple_bank", "secret"))
	if err != nil {
		log.Fatalf("failed opening connection to sqllite: %v", err)
	}

	defer client.Close()

	if err := client.Debug().Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}
