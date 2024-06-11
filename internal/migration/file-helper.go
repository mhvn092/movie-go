package migration

import (
	"github.com/mhvn092/movie-go/internal/util"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func readMigrationsFromDirSorted() []os.DirEntry {
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

func getMigrationFilePath(filename string, revert bool) string {
	pwd, _ := os.Getwd()
	folder := "up"
	if revert {
		folder = "down"
	}
	return filepath.Join(pwd+"migrations/"+folder, filename)
}

func readMigrationFile(filename string, revert bool) *os.File {
	migrationFilePath := getMigrationFilePath(filename, revert)

	file, err := os.Open(migrationFilePath)
	if err != nil {
		util.ErrorExit(err, "could not open the migration file")
	}
	return file
}

func getTheMigrationNameIndex() int {
	files := readMigrationsFromDirSorted()
	return len(files)
}

func CreateMigrationFile(name string) {
	index := getTheMigrationNameIndex()
	trimmedName := strings.TrimSpace(name)
	spacedArray := strings.Split(trimmedName, " ")
	if len(spacedArray) > 1 {
		trimmedName = strings.Join(spacedArray, "_")
	}
	finalName := strconv.Itoa(index) + "_" + trimmedName + ".sql"
	upPath := getMigrationFilePath(finalName, false)
	downPath := getMigrationFilePath(finalName, true)

	_, err := os.Create(upPath)
	if err != nil {
		util.ErrorExit(err, "could not create up migration file")
	}

	_, err = os.Create(downPath)
	if err != nil {
		util.ErrorExit(err, "could not create down migration file")
	}
}
