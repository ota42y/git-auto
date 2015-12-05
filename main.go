package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "ota42y"
	app.Email = ""
	app.Usage = ""

	app.Flags = flags
	app.Action = start

	app.Run(os.Args)
}
