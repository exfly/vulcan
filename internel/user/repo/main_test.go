package user

import (
	"log"
	"os"
	"testing"

	"github.com/exfly/vulcan/internel/config"
	dbbuilder "github.com/exfly/vulcan/internel/database"
	"github.com/exfly/vulcan/pkg/database"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var testQueries *Queries
var testDB *database.DB

func TestMain(m *testing.M) {
	cfg, err := config.LoadConfig("../../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = dbbuilder.New(cfg)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
