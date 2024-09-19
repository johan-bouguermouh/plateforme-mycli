package main

import (
	"log"
	"os"

	cmd "bucketool/cmd"

	"github.com/urfave/cli"
	bolt "go.etcd.io/bbolt"
)


func main() {
  app := cli.NewApp()

  // init bbolt
  db := UseBolt()
  defer db.Close()

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

