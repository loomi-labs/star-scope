package database

import (
	"fmt"
	"github.com/loomi-labs/star-scope/ent"
	"github.com/loomi-labs/star-scope/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func testClient(t *testing.T) *ent.Client {
	opts := []enttest.Option{
		enttest.WithOptions(ent.Log(t.Log)),
	}
	filename := fmt.Sprintf("file:ent%v?mode=memory&cache=shared&_fk=1", t.Name())
	client := enttest.Open(t, "sqlite3", filename, opts...)
	return client
}

func closeTestClient(client *ent.Client) {
	//goland:noinspection GoUnhandledErrorResult
	defer client.Close()
}

func newTestDbManagers(t *testing.T) *DbManagers {
	userManager := newTestUserManager(t)
	return &DbManagers{
		UserManager: userManager,
	}
}
