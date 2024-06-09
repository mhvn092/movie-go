package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhvn092/movie-go/internal/util"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var appliedMigrations = make(map[string]bool)

func main() {
	conn := util.InitDb()
	defer conn.Close()
	ensureMigrationTable(conn)
	readAllMigrationsFromDb(conn)
}

func ensureMigrationTable(conn *pgxpool.Pool) {
	_, err := conn.Query(context.Background(), `CREATE TABLE If NOT EXISTS migrations (
			id serial PRIMARY KEY,
			name VARCHAR(100),
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			);`)
	if err != nil {
		util.ErrorExit(err, "could not query the migrations table")
	}
}

func readAllMigrationsFromDb(conn *pgxpool.Pool) {
	rows, err := conn.Query(context.Background(), `SELECT name FROM migrations ORDER BY applied_at ASC;`)
	if err != nil {
		util.ErrorExit(err, "could not query the migrations table")
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			util.ErrorExit(err, "could not read the row from the migrations table")
		}
		appliedMigrations[name] = true
	}
	if err := rows.Err(); err != nil {
		util.ErrorExit(err, "could not read the rows from the migrations table")
	}
}

func readMigrationsFromDirAndApply(conn *pgxpool.Pool) {
	pwd, _ := os.Getwd()
	files, err := os.ReadDir(pwd + "/migrations/up")
	if err != nil {
		util.ErrorExit(err, "could not read the migrations directory")
	}

	// Sort the migration files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			var nameWithoutExtension = strings.TrimSuffix(file.Name(), ".sql")
			if _, applied := appliedMigrations[nameWithoutExtension]; !applied {
				applyMigration(conn, file.Name(), false)
			}
		}
	}
}

func applyMigration(conn *pgxpool.Pool, filename string, revert bool) {
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
	defer file.Close()
	// Split the SQL statements by `;` and execute them within a transaction
	sqlStatements := parseSQLStatements(file)

	runMigrationStatementsInTransaction(sqlStatements, conn, filename, revert)
}

func runMigrationStatementsInTransaction(sqlStatements map[int]string, conn *pgxpool.Pool, filename string, revert bool) {
	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to begin transaction for migration %s: %v\n", filename, err)
	}

	defer func() {
		if err != nil {
			e := tx.Rollback(ctx)
			if e != nil {
				util.ErrorExit(err, "could not rollback transaction")
			}
			util.ErrorExit(err, fmt.Sprintf("Transaction rolled back for migration %s", filename))
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				util.ErrorExit(err, fmt.Sprintf("Failed to commit transaction for migration %s", filename))
			}
		}
	}()

	for _, statement := range sqlStatements {
		_, err = tx.Exec(ctx, statement)
		if err != nil {
			return // This will cause the deferred function to roll back the transaction
		}
	}

	var nameWithoutExtension = strings.TrimSuffix(filename, ".sql")
	// Log the applied migration
	migrationTableStatement := `INSERT INTO migrations (name) VALUES ($1)`
	if revert {
		migrationTableStatement = `DELETE FROM migration where name = $1`
	}
	_, err = tx.Exec(ctx, migrationTableStatement, nameWithoutExtension)
	if err != nil {
		return // This will cause the deferred function to roll back the transaction
	}

	fmt.Printf("Applied migration: %s\n", filename)
}

func readTheLastMigrationsFromDb(conn *pgxpool.Pool) string {
	rows, err := conn.Query(context.Background(), `SELECT name FROM migrations order by applied_at DESC LIMIT 1;`)
	if err != nil {
		util.ErrorExit(err, "could not query the migrations table")
	}
	defer rows.Close()
	var name string
	if err := rows.Scan(&name); err != nil {
		util.ErrorExit(err, "could not read the row from the migrations table")
	}
	return name
}

func revertTheLastCommitedMigration(conn *pgxpool.Pool) {
	lastMigrationName := readTheLastMigrationsFromDb(conn)
	applyMigration(conn, lastMigrationName+".sql", true)
}

func parseSQLStatements(file *os.File) map[int]string {
	scanner := bufio.NewScanner(file)
	sqlStatements := make(map[int]string)
	var sqlBuilder strings.Builder
	statementID := 0
	hasContent := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			hasContent = true
			sqlBuilder.WriteString(trimmedLine)
			sqlBuilder.WriteString(" ")
			// If the line ends with a semicolon, consider it the end of a statement
			if strings.HasSuffix(trimmedLine, ";") {
				sqlStatements[statementID] = sqlBuilder.String()
				sqlBuilder.Reset()
				statementID++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		util.ErrorExit(err, "Failed to scan migration file: %v\n")
	}

	// Add any remaining SQL statement that doesn't end with a semicolon
	if sqlBuilder.Len() > 0 {
		sqlStatements[statementID] = sqlBuilder.String()
	}

	if !hasContent {
		util.ErrorExit(errors.New("migration file is empty"), "Empty file")
	}

	return sqlStatements
}
