package cmd

import (
	"bucketool/model"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var aliasFlags = []cli.Flag{}


var subCommands = []cli.Command{
	SetAliasCMD,
	ListAliasCMD,
	CurrentAliasCMD,
	DeleteAliasCMD,
}

var Store *model.AliasStore

func AliasCmd(db *bolt.DB) cli.Command {
	Store = model.UseAliasStore(db)
	return cli.Command{
		Name:    "alias",
		Aliases: []string{"a"},
		Usage:   "alias management",
		Flags:   aliasFlags,
		Subcommands: subCommands,
	}
}