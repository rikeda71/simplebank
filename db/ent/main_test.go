package db

import (
	"fmt"
	"testing"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/enttest"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// createEntClient create ent client and return client object
// enttest.Open が testing.T を引数に取るため、testing.M で初期化することが難しそう
func createEntClient(t *testing.T) *ent.Client {
	// sqlite3 の制約で並列に接続するとテーブルロックがかかってテストが実行できないのでpsqlを利用
	// client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	client := enttest.Open(t, "postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		"localhost", "5432", "root", "simple_bank", "secret"))
	return client
}
