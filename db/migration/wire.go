package migration

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/pkg/errors"
)

func NewMigration(dbSource string) (*migrate.Migrate, error) {
	dr, err := iofs.New(MigrateFS, ".")
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	m, err := migrate.NewWithSourceInstance("iofs", dr, dbSource)
	if err != nil {
		return nil, errors.Wrap(err, "NewWithSourceInstance")
	}

	return m, nil
}
