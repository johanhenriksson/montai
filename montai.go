package main

import (
	"github.com/johanhenriksson/montai/cli"
	"github.com/johanhenriksson/montai/server"
)

func main() {
	if !cli.IsContainer() {
		cli.Run()
	} else {
		server.Serve()
	}
}
