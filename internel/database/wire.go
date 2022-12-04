package database

import (
	"github.com/exfly/vulcan/pkg/database"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	New,
	wire.Bind(new(DBTX), new(*database.DB)),
)
