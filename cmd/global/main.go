package global

import (
	conn "bucketool/connexion"
	model "bucketool/model"
	color "bucketool/utils/colorPrint"
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)


var Store *model.AliasStore
var Connexion *conn.Connexion

func useAliasDefaultName() string{
	
	currentAlias, err := Store.GetCurrentAlias()
	if err != nil {
		if strings.Contains(err.Error(), "404 |") {
			return "default"
		}
		 fmt.Println(color.RedP("Can't get current alias"), err)
        return "default"
	}
	if currentAlias.Name == "" {
		return "default"
	}
	return currentAlias.Name
}

func BeforeUseAlias(c *cli.Context) error {
	if c.Args()[0] == "alias" {
		return nil
	}

	aliasName := c.String("alias")
	if aliasName == "" {
		aliasName = useAliasDefaultName()
	}
	if !Store.IsAliasExist(aliasName){
		if(aliasName == "default"){
			println(color.YellowP("Warning : No Alias registred, please use the command 'bucketool alias set' or 'bucketool alias current -s <alias>' to add an Alias"))
		}
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
