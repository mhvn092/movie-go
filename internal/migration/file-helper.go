package migration

import (
	"github.com/mhvn092/movie-go/internal/util"
	"os"
	"path/filepath"
	"sort"
)

func ReadMigrationsFromDir() []os.DirEntry {
	pwd, _ := os.Getwd()
	files, err := os.ReadDir(pwd + "/migrations/up")
	if err != nil {
		util.ErrorExit(err, "could not read the migrations directory")
	}

	// Sort the migration files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	return files
}

func ReadMigrationFile(filename string, revert bool) *os.File {
	pwd, _ := os.Getwd()
	folder := "up"
	if revert {
		folder = "down"
	}
	migrationFilePath := filepath.Join(pwd+"migrations/"+folder, filename)
	file, err := os.Open(migrationFilePath)
	if err != nil {
		util.ErrorExit(err, "could not open the migration file")
	}
	return file
}
