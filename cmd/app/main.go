package main

import (
	"errors"
	"fmt"
	"locky/cmd/cli"
	"locky/internal/project"
	"log"
	"os"
)

func main() {
	config, err := cli.Parse(os.Args)

	if err != nil {
		fatalError(err)
	}

	switch config.Cmd {
	case "new":
		data, err := cli.ParseNew(config.CmdArgs)

		if err != nil {
			fatalError(err)
		}

		err = project.GenServer(project.Config{Name: data.ProjectName, Open: data.Open}, os.Stderr)

		if err != nil {
			log.Fatalln(err)
		}
	default:
		fatalError(errors.New("command is unspecified"))
	}
}

func fatalError(err error) {
	fmt.Fprintln(os.Stderr, err)
	log.Fatalln(cli.Usage())
}
