package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-romajify"
	app.Usage = "Convert kana to romaji"
	app.Version = "0.1.0"
	app.Commands = Commands
	app.Run(os.Args)
}
