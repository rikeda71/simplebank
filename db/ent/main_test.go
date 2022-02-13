package db

import (
	"testing"

	"github.com/s14t284/simplebank/ent"
	"github.com/s14t284/simplebank/ent/enttest"
	"github.com/s14t284/simplebank/util"

	_ "github.com/lib/pq"
	// _ "github.com/mattn/go-sqlite3"
)

// createEntClient create ent client and return client object
// enttest.Open が testing.T を引数に取るため、testing.M で初期化することが難しそう
func createEntClient(t *testing.T) *ent.Client {
	// sqlite3 の制約で並列に接続するとテーブルロックがかかってテストが実行できないのでpsqlを利用
	// client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	config, err := util.LoadConfig("../..")
	if err != nil {
		panic(err)
	}
	client := enttest.Open(t, config.DBDriver, config.DBSource)
	return client
}
