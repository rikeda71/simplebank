package db

import (
	"testing"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
)

// createEntClient create ent client and return client object
// enttest.Open が testing.T を引数に取るため、testing.M で初期化することが難しそう
func createEntClient(t *testing.T) *ent.Client {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	return client
}
