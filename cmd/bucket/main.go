package bucket

import (
	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var bucketFlags = []cli.Flag{

}

var subCommands = []cli.Command{
	CreateBucketCMD,
	ListBucketCMD,
	DeleteBucketCMD,
}

func BucketCmd(db *bolt.DB) cli.Command {
	return cli.Command{
		Name:    "bucket",
		Aliases: []string{"b"},
		Usage:   "bucket management",
		Flags:   bucketFlags,
		Subcommands: subCommands,
	}
}