package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/seepre/go-lamp/cli/lamp/gen"
	"github.com/urfave/cli"
)

var (
	BuildVersion = "20210618"
	commands     = []cli.Command{
		{
			Name:  "new",
			Usage: "generate engineering framework",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "generate framework by this name",
				},
				cli.StringFlag{
					Name:  "d",
					Usage: "generate framework into this dir",
				},
			},
			Action: gen.GeneratorCommand,
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "a cli tool to generate code"
	app.Version = fmt.Sprintf("%s %s/%s", BuildVersion, runtime.GOOS, runtime.GOARCH)
	app.Commands = commands
	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
