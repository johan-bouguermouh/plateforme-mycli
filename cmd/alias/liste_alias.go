package alias

import (
	co "bucketool/connexion"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"strings"

	"github.com/urfave/cli"
)

// Description how to use the command
var listUsageText string = color.GreyP("Command Exemple : bucketool list \n")+
    `This command will list all Alias Registred.
    You can use the flag -d or --detail to list all Alias registred with details url.
    Exemple : bucketool list -d
    Return : 
    *Alias Name : <AliasName> <hostanme>:<port>
    Also you can use the flag -f or --filter to filter Alias registred by name or if Hostname contain the filter.
    Exemple : bucketool list -f <filter>
    Return :
     - *Alias Name : <AliasName how contain the filter>
     - *Alias Name : <hostName how contain the filter>
     If the current Alias used, it will be marked with cursor ►
    `
    
// Description of the flags
var listAliasDesc string = color.ColorPrint("Black",
    " Flags Descriptions ", &color.Options{
        Underline: true,
        Bold: 	true,
        Background: "White",})+ "\n"+
        color.BlueP(" -d, --detail ") +  color.GreyP("(Optionnal)")+" : List all Alias registred with details url\n"+
        color.BlueP(" -f, --filter ") + color.GreyP("(Optionnal)") + " : Filter Alias registred by name or if Hostname contain the filter\n"



var ListFlags = []cli.Flag{
    cli.BoolFlag{
        Name: "detail, d",
        Usage: "List all Alias registred with details url",
    },
    cli.StringFlag{
        Name: "filter, f",
        Usage: "Filter Alias registred by name or if Hostname contain the filter",
    },
}


var ListAliasCMD = cli.Command{
        Name:    "list",
        Category: "Alias",
        Aliases: []string{"ls"},
        Usage:   "list all Alias Registred",
        UsageText: listUsageText,
        Description: listAliasDesc,
        Flags: ListFlags,
        OnUsageError: utils.OnUsageError,
        Action: ListAliasCmd,
        CustomHelpTemplate: utils.HelpTemplate,
}

func ListAliasCmd(c *cli.Context) error {
    aliass, err := Store.ListAliass()
    if err != nil {
        println("Erreur lors de la lecture des tâches:", err)
        return err
    }
    // Check if the flag of filter is used
    if c.String("filter") != "" {
        println("Liste was filtered by : ", color.BlueP(c.String("filter")))
    }


    for _, Alias := range aliass {
        prefix := "  -"
        if Alias.Current {
            prefix = "  ►"
        }
        if c.String("filter") != "" && !strings.Contains(Alias.Name, c.String("filter")) && !strings.Contains(Alias.HOST, c.String("filter")) {
            continue
        }
        line := prefix + " "+color.GreenP(Alias.Name) + " "
        if c.Bool("detail") {
            line += color.GreyP("("+co.CreateURL(Alias)+")")
        }
        println(line)
    }

    return nil
}