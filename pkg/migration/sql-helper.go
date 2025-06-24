package migration

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/mhvn092/movie-go/pkg/exception"
)

func parseSQLStatements(file *os.File) map[int]string {
	scanner := bufio.NewScanner(file)
	sqlStatements := make(map[int]string)
	var sqlBuilder strings.Builder
	statementID := 0
	hasContent := false
	currentDelimiter := ";"

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		hasContent = true

		if strings.HasPrefix(trimmedLine, `-- delimiter`) {
			parts := strings.Fields(trimmedLine)
			if len(parts) == 3 && parts[1] == "delimiter" {
				currentDelimiter = parts[2]
			}
			continue
		}

		sqlBuilder.WriteString(trimmedLine)
		sqlBuilder.WriteString(" ")

		if strings.HasSuffix(trimmedLine, currentDelimiter) {
			stmt := strings.TrimSpace(sqlBuilder.String())
			stmt = strings.TrimSuffix(stmt, currentDelimiter)
			sqlStatements[statementID] = strings.TrimSpace(stmt)
			sqlBuilder.Reset()
			statementID++
			currentDelimiter = ";"
		}
	}

	if err := scanner.Err(); err != nil {
		exception.ErrorExit(err, "Failed to scan migration file: %v\n")
	}

	// Add any remaining SQL statement that doesn't end with a semicolon
	if sqlBuilder.Len() > 0 {
		sqlStatements[statementID] = sqlBuilder.String()
	}

	if !hasContent {
		exception.ErrorExit(errors.New("migration file is empty"), "Empty file")
	}

	return sqlStatements
}

func getLastMigrationQuery() string {
	return `SELECT name FROM migrations order by applied_at DESC LIMIT 1;`
}

func getAllMigrationQuery() string {
	return `SELECT name FROM migrations ORDER BY applied_at ASC;`
}

func getUpdatingMigrationTableQuery(revert bool) string {
	migrationTableStatement := `INSERT INTO migrations (name) VALUES ($1)`
	if revert {
		migrationTableStatement = `DELETE FROM migrations where name = $1`
	}
	return migrationTableStatement
}

func checkExistenceOfMigrationTableQuery() string {
	return `CREATE TABLE If NOT EXISTS migrations (
			id serial PRIMARY KEY,
			name VARCHAR(100),
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
			);`
}
