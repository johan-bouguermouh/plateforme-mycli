package alias

import (
	conn "bucketool/connexion"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

var deleteUsageText string = color.GreyP("Command Exemple : bucketool alias delete <name> \n") +
	`This command will delete the alias from the list of alias.
	The first argument is the name of the alias.
	You can also use the flag -a or -all to delete all the alias, you don't need to specify the name of the alias in this case.
	When you use the flag -a or -all, the current alias will be set to nil. You can use flage -sc or -savecurrent to not delete the current alias
	You have second flag -clean or -c to delete all the alias can't be connected at server`
	
var deleteDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold: true,
		Background: "White",
	}) + "\n" +
	color.BlueP(" -a, --all ") + color.GreyP("(Optionnal)") + " : Delete all the alias\n" +
	color.BlueP(" -sc, --savecurrent ") + color.GreyP("(Optionnal)") + " : Save the current alias\n whene you use flag -a, -all" +
	color.BlueP(" -c, --clean ") + color.GreyP("(Optionnal)") + " : Delete all the alias can't be connected at server\n"

var deleteFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "all, a",
		Usage: "Delete all the alias",
	},
	cli.BoolFlag{
		Name:  "savecurrent, sc",
		Usage: "Save the current alias whene you use flag -a, -all",
	},
	cli.BoolFlag{
		Name:  "clean, c",
		Usage: "Delete all the alias can't be connected at server",
	},
}

var DeleteAliasCMD = cli.Command{
	Name:        "delete",
	Category:   "Alias",
	Usage:       "Delete an alias",
	UsageText:   deleteUsageText,
	Description: deleteDesc,
	Flags:       deleteFlags,
	Action:      deleteAliasCmd,
	After: ListAliasCmd,
	OnUsageError: utils.OnUsageError,
	CustomHelpTemplate: utils.HelpTemplate,
}


func deleteAliasCmd(c *cli.Context) error {
	currentAlias, err := Store.GetCurrentAlias()
	aliasArg := c.Args().First()
	var couldBeSavedCurrent bool = false
	if err != nil {
		return err
	}


	if (c.Bool("all") && !c.Bool("savecurrent")) ||
	(c.Bool("clean") && !c.Bool("savecurrent")) ||
	(aliasArg == currentAlias.Name && !c.Bool("savecurrent")) {
		userRes, err := confirmAction("The current Alias same as touched by this command will be deleted, do you want to delete ? (y/n) : ")
		if err != nil {
			return err
		}
		if userRes {
			couldBeSavedCurrent = true
		}
	}

	if c.Bool("all") {
		println(color.ColorPrint("Grey", "Deleting All Alias...", &color.Options{
			Italic: true,
		}))
			err = Store.DeleteAllAlias()
			if err != nil {
				return err
			}
	}
	if c.Bool("clean") && !c.Bool("all") {
		println(color.ColorPrint("Grey", "Cleaning Alias...", &color.Options{
			Italic: true,
			}))
		err := cleanAlias()
		if err != nil {
			return err
		}
	}
	if(aliasArg != "" && !Store.IsAliasExist(aliasArg)){
		return errors.New(color.RedP("Alias not found"))
	} else if(aliasArg != "" && !couldBeSavedCurrent) {
		err := Store.DeleteAliasByName(aliasArg)
		if err != nil {
			return err
		}
		println("Alias deleted :",color.GreenP(aliasArg))
	}

	if(currentAlias.Name == aliasArg && !c.Bool("savecurrent") && !couldBeSavedCurrent){
		println(color.YellowP("WARN | The current Alias has been deleted"))
	}

	if c.Bool("savecurrent") || couldBeSavedCurrent {
		err = Store.SaveAlias(&currentAlias)
		if err != nil {
			println(color.RedP("Error while saving the current Alias"))
			return err
		}
		}

	
	return nil
}

func confirmAction(prompt string) (bool, error) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(prompt)
    response, err := reader.ReadString('\n')
    if err != nil {
        return false, err
    }
    response = strings.TrimSpace(response)
    return strings.ToLower(response) != "y", nil
}

func cleanAlias() error {
	alias, err := Store.ListAliass()
	if err != nil {
		return err
	}
	for _, a := range alias {
		co := conn.Use(a)
		_, err := co.Connect()
		if err != nil {

			err = Store.DeleteAliasByName(a.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
