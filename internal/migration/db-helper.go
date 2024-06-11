package migration

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mhvn092/movie-go/internal/util"
	"log"
	"strings"
)

func ReadTheLastMigrationsFromDb(conn *pgxpool.Pool) string {
	rows, err := conn.Query(context.Background(), getLastMigrationQuery())
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

func RevertTheLastCommitedMigration(conn *pgxpool.Pool) {
	lastMigrationName := ReadTheLastMigrationsFromDb(conn)
	ApplyMigration(conn, lastMigrationName+".sql", true)
}

func ApplyMigration(conn *pgxpool.Pool, filename string, revert bool) {
	file := ReadMigrationFile(filename, revert)
	defer file.Close()
	// Split the SQL statements by `;` and execute them within a transaction
	sqlStatements := ParseSQLStatements(file)

	RunMigrationStatementsInTransaction(sqlStatements, conn, filename, revert)
}

func RunMigrationStatementsInTransaction(sqlStatements map[int]string, conn *pgxpool.Pool, filename string, revert bool) {
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
	migrationTableStatement := getUpdatingMigrationTableQuery(revert)

	_, err = tx.Exec(ctx, migrationTableStatement, nameWithoutExtension)
	if err != nil {
		return // This will cause the deferred function to roll back the transaction
	}

	fmt.Printf("Applied migration: %s\n", filename)
}

func EnsureMigrationTable(conn *pgxpool.Pool) {
	_, err := conn.Query(context.Background(), checkExistenceOfMigrationTableQuery())
	if err != nil {
		util.ErrorExit(err, "could not query the migrations table")
	}
}

func ReadAllMigrationsFromDb(conn *pgxpool.Pool) map[string]bool {
	var appliedMigrations = make(map[string]bool)
	rows, err := conn.Query(context.Background(), getAllMigrationQuery())
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
	return appliedMigrations
}

func ReadMigrationsFromDirAndApply(conn *pgxpool.Pool, appliedMigrations map[string]bool) {
	files := ReadMigrationsFromDirSorted()

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			var nameWithoutExtension = strings.TrimSuffix(file.Name(), ".sql")
			if _, applied := appliedMigrations[nameWithoutExtension]; !applied {
				ApplyMigration(conn, file.Name(), false)
			}
		}
	}
}
