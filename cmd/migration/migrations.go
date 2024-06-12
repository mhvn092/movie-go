package main

import (
	"errors"
	"github.com/mhvn092/movie-go/internal/migration"
	"github.com/mhvn092/movie-go/internal/util"
	"os"
	"strings"
)

var validCommand = [3]string{"up", "down", "create"}

func main() {
	if len(os.Args) < 2 {
		util.ErrorExit(errors.New("no args provided"), "you should provide the arguments")
	}

	command := os.Args[1]
	if !isValidCommand(command) {
		util.ErrorExit(errors.New("args not valid"), "you should provide valid arguments")
	}

	switch command {
	case "create":
		handleCreateCommand(os.Args[2:])
	case "up":
		handleUpCommand()
	case "down":
		handleDownCommand()
	default:
		util.ErrorExit(errors.New("unknown command"), "unknown command")
	}
}

func isValidCommand(command string) bool {
	for _, valid := range validCommand {
		if valid == command {
			return true
		}
	}
	return false
}

func handleCreateCommand(args []string) {
	if len(args) < 2 {
		util.ErrorExit(errors.New("no name provided"), "you should provide the name")
	}

	name := strings.Join(args[1:], " ")

	if name == "" {
		util.ErrorExit(errors.New("no name provided"), "you should provide the name")
	}

	migration.CreateMigrationFile(name)
}

func handleUpCommand() {
	conn := util.InitDb()
	defer conn.Close()
	migration.RunMigrations(conn)
}

func handleDownCommand() {
	conn := util.InitDb()
	defer conn.Close()
	migration.RevertTheLastMigration(conn)
}
