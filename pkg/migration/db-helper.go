package migration

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhvn092/movie-go/pkg/exception"
	"log"
	"strings"
)

func ensureMigrationTable(conn *pgxpool.Pool) {
	row, err := conn.Query(context.Background(), checkExistenceOfMigrationTableQuery())
	defer row.Close()
	if err != nil {
		exception.ErrorExit(err, "could not query the migrations table")
	}
	fmt.Println("ensured migrations table exist")
}

func readAllMigrationsFromDb(conn *pgxpool.Pool) map[string]bool {
	var appliedMigrations = make(map[string]bool)
	rows, err := conn.Query(context.Background(), getAllMigrationQuery())
	if err != nil {
		exception.ErrorExit(err, "could not query the migrations table")
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			exception.ErrorExit(err, "could not read the row from the migrations table")
		}
		appliedMigrations[name] = true
	}
	if err := rows.Err(); err != nil {
		exception.ErrorExit(err, "could not read the rows from the migrations table")
	}
	return appliedMigrations
}

func readTheLastMigrationsFromDb(conn *pgxpool.Pool) string {
	rows, err := conn.Query(context.Background(), getLastMigrationQuery())
	if err != nil {
		exception.ErrorExit(err, "could not query the migrations table")
	}
	defer rows.Close()

	var name string

	if rows.Next() {
		if err := rows.Scan(&name); err != nil {
			exception.ErrorExit(err, "could not read the row from the migrations table")
		}
	} else {
		if err := rows.Err(); err != nil {
			exception.ErrorExit(err, "error occurred during row iteration")
		}
		exception.ErrorExit(errors.New("no rows found"), "could not find any migrations")
	}
	return name
}

func revertTheLastCommitedMigration(conn *pgxpool.Pool) {
	lastMigrationName := readTheLastMigrationsFromDb(conn)
	applyMigration(conn, lastMigrationName+".sql", true)
}

func readMigrationsFromDirAndApply(conn *pgxpool.Pool, appliedMigrations map[string]bool) {
	files := readMigrationsFromDirSorted()
	if len(files) == 0 {
		fmt.Println("no migration to run")
		return
	}
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
	file := readMigrationFile(filename, revert)
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
				exception.ErrorExit(err, "could not rollback transaction")
			}
			exception.ErrorExit(err, fmt.Sprintf("Transaction rolled back for migration %s", filename))
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				exception.ErrorExit(err, fmt.Sprintf("Failed to commit transaction for migration %s", filename))
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
	migrationTableStatement := getUpdatingMigrationTableQuery(revert)

	_, err = tx.Exec(ctx, migrationTableStatement, nameWithoutExtension)
	if err != nil {
		return // This will cause the deferred function to roll back the transaction
	}

	fmt.Printf("Applied migration: %s\n", filename)
}

func RunMigrations(conn *pgxpool.Pool) {
	ensureMigrationTable(conn)
	appliedMigrations := readAllMigrationsFromDb(conn)
	readMigrationsFromDirAndApply(conn, appliedMigrations)
}

func RevertTheLastMigration(conn *pgxpool.Pool) {
	ensureMigrationTable(conn)
	revertTheLastCommitedMigration(conn)
}
