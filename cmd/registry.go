package cmd

import (
	alias "bucketool/cmd/alias"
	bucket "bucketool/cmd/bucket"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var CommandRegistry = []cli.Command{}


func RegisterCommands(db *bolt.DB) {
    CommandRegistry = []cli.Command{
        alias.AliasCmd(db),
		bucket.BucketCmd(db),

    }
}