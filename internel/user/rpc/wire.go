package rpc

import (
	"github.com/google/wire"

	vulcanv1 "github.com/exfly/vulcan/pb"
)

var ProviderSet = wire.NewSet(
	New,
	wire.Bind(new(vulcanv1.SimpleBankServer), new(*Server)),
)
