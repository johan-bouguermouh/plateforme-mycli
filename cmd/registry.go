package cmd

import (
	alias "bucketool/cmd/alias"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var CommandRegistry = []cli.Command{}


func RegisterCommands(db *bolt.DB) {
    CommandRegistry = []cli.Command{
        alias.AliasCmd(db),
    }
}