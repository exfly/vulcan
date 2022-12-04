//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/exfly/vulcan/internel/config"
	"github.com/exfly/vulcan/internel/database"
	userrepo "github.com/exfly/vulcan/internel/user/repo"
	userrpc "github.com/exfly/vulcan/internel/user/rpc"
)

func NewApp(config.Path) (*App, error) {
	panic(wire.Build(
		database.ProviderSet,
		config.ProviderSet,
		userrepo.ProviderSet,
		userrpc.ProviderSet,
		newApp,
	))
}
