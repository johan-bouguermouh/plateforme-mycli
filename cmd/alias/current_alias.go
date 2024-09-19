package cmd

import (
	"errors"

	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"

	"github.com/urfave/cli"
)

// Description how to use the command
var currentAliasUageText string = color.GreyP("Command Exemple : bucketool current -switch <AliasName> \n")+
    `This command will switch the current Alias to the AliasName.
    You can use the flag -s or --switch to switch the current Alias to the AliasName.
    Exemple : bucketool current -switch <AliasName>
    Return :
    *AliasName is now the current Alias
    If you dont use the flag -s or --switch, the current Alias will be returned.
    `

// Description of the flags
var currentAliasDesc string = color.ColorPrint("Black",
    " Flags Descriptions ", &color.Options{
        Underline: true,
        Bold: 	true,
        Background: "White",})+ "\n"+
        color.BlueP(" -s, --switch ") +  color.GreyP("(Optionnal)")+" : Switch the current Alias to the AliasName\n"


    
var CurrentAliasFlags = []cli.Flag{
    cli.StringFlag{
        Name: "switch, s",
        Usage: "Switch the current Alias to the AliasName",
    },
}

var CurrentAliasCMD = cli.Command{
        Name:    "current",
        Category: "Alias",
        Aliases: []string{"c"},
        Usage:   "Gesture of the current Alias",
        UsageText: currentAliasUageText,
        Description: currentAliasDesc,
        OnUsageError: utils.OnUsageError,
        Flags: CurrentAliasFlags,
        Action: currentAliasCmd,
        CustomHelpTemplate: utils.HelpTemplate,
}

func currentAliasCmd(c *cli.Context) error {
    newCurrentAliasName := c.String("switch")
    if newCurrentAliasName == "" {
        currentAlias, err := Store.GetCurrentAlias()
        if err != nil {
            return err
        }
        println("Alias used : ", color.GreenP(currentAlias.Name))
        return nil
    }
    verifyAlias, err := Store.ReadAlias(newCurrentAliasName)
    if err != nil || Store.IsEmptyAlias(verifyAlias) {
        return errors.New(color.RedP("Alias not found"))
    }
    err = Store.SetCurrentAlias(newCurrentAliasName)
    if err != nil {
        return err
    }

    return nil
}