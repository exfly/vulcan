package user

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	New,
	NewStore,
	wire.Bind(new(Querier), new(*Queries)),
)
