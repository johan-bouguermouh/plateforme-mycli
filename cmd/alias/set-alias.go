package alias

import (
	conn "bucketool/connexion"
	env "bucketool/environment"
	"bucketool/model"
	utils "bucketool/utils"
	color "bucketool/utils/colorPrint"
	"context"

	"github.com/urfave/cli"
)

// Description how to use the command
var setAliasUsageText string = color.GreyP("Command Exemple : bucketool alias set <name> -p <port> -h <host> -k <keyname> -s <Secret> -c \n")+
    `This command will add a new alias to the list of alias.
	The first argument is the name of the alias.
	The flag -p or -port is required to set the port of the alias.
	The flag -H or -hostname is required to set the host of the alias.
	The flag -k or -keyname is required to set the keyname of the alias.
	The flag -s or -Secret is required to set the SecretKey of the alias.
    You can use the flag -c or -current to set the current alias directly.
    The Current Alias will be used by default when you use the command bucketool connect`

	// Description of the flags
var setAliasDesc string = color.ColorPrint("Black",
	" Flags Descriptions ", &color.Options{
		Underline: true,
		Bold: 	true,
		Background: "White",})+ "\n"+
		color.BlueP(" -p, --port ") + color.GreyP("(Required)")+ " : Port to connect\n"+
		color.BlueP(" -H, --hostname ") +  color.GreyP("(Required)")+" : Host to connect\n"+
		color.BlueP(" -k, --keyname ") + color.GreyP("(Required)") + " : KeyName to connect\n"+
		color.BlueP(" -s, --Secret ") + color.GreyP("(Required)") + " : SecretKey to connect\n"+
		color.BlueP(" -c, --current ") + color.GreyP("(Optionnal)") +" : Automatically set and registry the current alias\n"

var SettAliasFlags = []cli.Flag{
	cli.IntFlag{
		Name:  "port, p",
		Required : true,
		Usage: "Port to connect",
	},
	cli.StringFlag{
		Name:  "hostname, H",
		Required: true,
		Usage: "Host to connect",
	},
	cli.StringFlag{
		Name:  "keyname, k",
		Required: true,
		Usage: "KeyName to connect",
	},
	cli.StringFlag{
		Name:  "Secret, s",
		Required: true,
		Usage: "SecretKey to connect",
	},
	cli.BoolFlag{
		Name: "current, c",
		Usage: "Automatically set and registry the current alias",
	},
}

var SetAliasCMD = cli.Command{
        Name:    "set",
        Aliases: []string{"s"},
        Usage:   "add a task to the list",
		UsageText: setAliasUsageText,
		Description : setAliasDesc,
		Category : "Alias",
		Flags:   SettAliasFlags,
        Action: setAliasCmd,
		OnUsageError : utils.OnUsageError,
		CustomHelpTemplate : utils.HelpTemplate,
}


func setAliasCmd(c *cli.Context) error {
	PORT := c.Int("port")
	newAllias := model.Alias{
		ID:        0,
		Name:      c.Args().First(),
		Port:   PORT,
		HOST:   c.String("hostname"),
		KeyName:   c.String("keyname"),
		SecretKey:   c.String("Secret"),
		Current:   c.Bool("current"),
	}

	Store.SaveAlias(&newAllias)
	co := conn.Use(newAllias)

	//on fait un appel au ClientS3 pour s'assurer que l'alias est correct
	_, err :=  co.S3Client.ListBuckets(context.TODO(),nil)
	if err != nil {
		println(color.RedP("Error while connecting to the alias"))
		if env.IsDebugMode{
			println(err.Error())
		}
	}


    println("Registered Alias on Name: ",
	color.GreenP(newAllias.Name),
	color.ColorPrint("Grey", co.URL, 
	&color.Options{
		Italic: true,
	}))
	return nil
}

