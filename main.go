package main

import (
	"os"

	"github.com/rysmaadit/go-template/app"
	"github.com/rysmaadit/go-template/cli"
)

func main() {
	c := cli.NewCli(os.Args)
	c.Run(app.Init())
}
