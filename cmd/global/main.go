package global

import (
	conn "bucketool/connexion"
	model "bucketool/model"
	color "bucketool/utils/colorPrint"
	"errors"
	"fmt"

	"github.com/urfave/cli"
)


var Store *model.AliasStore
var Connexion *conn.Connexion

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
	aliasName := c.String("alias")
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
