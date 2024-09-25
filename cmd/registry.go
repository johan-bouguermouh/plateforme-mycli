package cmd

import (
	alias "bucketool/cmd/alias"
	bucket "bucketool/cmd/bucket"
	global "bucketool/cmd/global"

	model "bucketool/model"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var CommandRegistry = []cli.Command{}



func RegisterCommands(db *bolt.DB) {
	global.Store = model.UseAliasStore(db)


    CommandRegistry = []cli.Command{
        alias.AliasCmd(db),
		bucket.BucketCmd(db),
		CopyObjectCMD,
		ListBucketObjectsCMD,
		DownloadObjectCMD,
    }
}