package main

import (
	"log"
	"os"

	cmd "bucketool/cmd"
	global "bucketool/cmd/global"
	env "bucketool/environment"
	model "bucketool/model"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)

var GeneralFlags = []cli.Flag{
  cli.BoolFlag{
    Name:  "debug",
    Usage: "Enable debug mode",
    EnvVar: "DEBUG",
    Destination: &env.IsDebugMode,
  },
  cli.StringFlag{
		Name : "alias",
		Usage: "Alias name to use",
		Value : "",
	},
  }

func main() {
  app := cli.NewApp()

  // init bbolt
  db := UseBolt()
  defer db.Close()

  //integration of flags in the app
  app.Flags = GeneralFlags
  global.Store = model.UseAliasStore(db)
  //execute the command before all the command
  app.Before = global.BeforeUseAlias

  cmd.RegisterCommands(db)

  app.Commands = cmd.CommandRegistry

  err := app.Run(os.Args)
  if err != nil {
  log.Fatal(err)
  }
}

func UseBolt() *bolt.DB {
  db, err := bolt.Open("data.db", 0600, nil)
  if err != nil {
      log.Fatal(err)
  }
  return db
}

