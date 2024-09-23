package bucket

import (
	"bucketool/model"
	"errors"
	"fmt"

	conn "bucketool/connexion"
	color "bucketool/utils/colorPrint"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var aliasFlags = []cli.Flag{
	cli.StringFlag{
		Name : "alias",
		Usage: "Alias name to use",
		Value : "",
	},

}

var subCommands = []cli.Command{
	CreateBucketCMD,
	ListBucketCMD,
	DeleteBucketCMD,
}

var Store *model.AliasStore
var Connexion *conn.Connexion

func BucketCmd(db *bolt.DB) cli.Command {
	Store = model.UseAliasStore(db)
	return cli.Command{
		Name:    "bucket",
		Aliases: []string{"b"},
		Usage:   "bucket management",
		Flags:   aliasFlags,
		Subcommands: subCommands,
	}
}

func useAliasDefaultName() string{
	currentAlias, err := Store.GetCurrentAlias()
	if err != nil {
		 fmt.Println(color.RedP("Can't get current alias"), err)
        return "default"
	}
	if currentAlias.Name == "" {
		return "default"
	}
	return currentAlias.Name
}

func BeforeUseAlias(c *cli.Context) error {
	ctxParent := c.Parent()
	aliasName := ctxParent.String("alias")
	if aliasName == "" {
		aliasName = useAliasDefaultName()
	}
	if !Store.IsAliasExist(aliasName){
		return errors.New(color.RedP("error while using alias : alias not found"))
	}
	currentAlias, err := Store.ReadAlias(aliasName)
	if err != nil {
		fmt.Println(color.RedP("Can't get alias"), err)
		return err
	}
	Connexion = conn.Use(currentAlias)
	
	return nil
}