package main

import (
	"errors"
	"flag"
	"github.com/mhvn092/movie-go/internal/migration"
	"github.com/mhvn092/movie-go/internal/util"
	"strings"
)

var validCommand = [3]string{"up", "down", "create"}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		util.ErrorExit(errors.New("no args Provided"), "you should provide the arguments")
	}

	command := flag.Arg(0)
	if !isValidCommand(command) {
		util.ErrorExit(errors.New("args Not Valid"), "you should provide the valid arguments")
	}

	if command == "create" {
		name := strings.Join(flag.Args()[1:], "_")
		migration.CreateMigrationFile(name)
		return
	}

	conn := util.InitDb()
	defer conn.Close()

	if command == "up" {
		migration.RunMigrations(conn)
	} else {
		migration.RevertTheLastMigration(conn)
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
