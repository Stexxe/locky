package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var commands map[string]struct{}

func init() {
	commands = make(map[string]struct{})
	commands["new"] = struct{}{}
}

func Usage() string {
	format := `
Locky 0.1.0
CLI tool for Ktor framework

USAGE: %s command

COMMANDS:
	new PROJECT_NAME [-open]	Generates new project
`
	return fmt.Sprintf(strings.Trim(format, "\n"), filepath.Base(os.Args[0]))
}

type General struct {
	Cmd string
	CmdArgs []string
}

func allCommands() []string {
	var result []string
	for name := range commands {
		result = append(result, name)
	}
	return result
}

func commandsMsg() string {
	return fmt.Sprintf("Available commands: %s", strings.Join(allCommands(), ", "))
}

func Parse(args []string) (General, error) {
	if len(args) < 2 {
		return General{}, errors.New(fmt.Sprintf("no command provided. %s", commandsMsg()))
	}

	cmd := args[1]

	_, ok := commands[cmd]

	if !ok {
		return General{}, errors.New(fmt.Sprintf("command %s does not exist. %s", cmd, commandsMsg()))
	}

	return General{Cmd: cmd, CmdArgs: args[2:]}, nil
}

type New struct {
	ProjectName string
	Open bool
}

func ParseNew(args []string) (New, error) {
	if len(args) < 1 {
		return New {}, errors.New("PROJECT_NAME argument is required for new command")
	}

	set := flag.NewFlagSet("new", flag.ContinueOnError)
	open := set.Bool("open", false, "")
	_ = set.Parse(args[1:])

	return New{ProjectName: args[0], Open: *open}, nil
}