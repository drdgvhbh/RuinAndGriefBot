// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	cli2 "drdgvhbh/discordbot/internal/cli"
	"drdgvhbh/discordbot/internal/db/pg"
	"drdgvhbh/discordbot/internal/user/api"
	"drdgvhbh/discordbot/internal/user/mapper"
	"github.com/urfave/cli"
)

// Injectors from wire.go:

func InitializeUserRepository() (*api.UserRepository, error) {
	config := pg.ProvideConfig()
	db := pg.ProvideConnector(config)
	userMapper := mapper.CreateUserMapper()
	userRepository := api.CreateUserRepository(db, userMapper)
	return userRepository, nil
}

func InitializeCLI() *cli.App {
	config := cli2.ProvideConfig()
	app := cli2.ProvideCLI(config)
	return app
}