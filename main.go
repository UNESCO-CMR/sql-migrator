package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	// "os"
	"path/filepath"
	"strings"

	_ "github.com/denisenkom/go-mssqldb" // Import MSSQL driver
)

const (
	server   = "127.0.0.1"
	port     = 80
	database = "myschoolonline"
	user     = "root"
	password = ""
)

func main() {
	// Create the connection string
	connString := fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s", server, port, database, user, password)

	// Connect to the SQL server
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Connection error:", err.Error())
	}
	defer db.Close()

	// Directory containing SQL files
	sqlDirectory := "/Users/matt-klaus/Documents/UNESCO/Programming/sql-migrate-go/db_split/SQLDumpSplitterResult"

	// Get the list of SQL files in the directory
	sqlFiles, err := getSQLFiles(sqlDirectory)
	if err != nil {
		log.Fatal("Error reading SQL files:", err.Error())
	}

	// Execute the schema file first
	schemaFile := "dbstructure.sql"
	schemaPath := filepath.Join(sqlDirectory, schemaFile)
	schemaSQL, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatal("Error reading schema file:", err.Error())
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		log.Fatal("Error executing schema file:", err.Error())
	}
	fmt.Println("Successfully executed schema file:", schemaFile)

	// Execute the remaining SQL files
	for _, sqlFile := range sqlFiles {
		// if sqlFile != schemaFile {
		filePath := filepath.Join(sqlDirectory, sqlFile)
		sql, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println("Error reading SQL file:", sqlFile, "-", err.Error())
			continue
		}

		_, err = db.Exec(string(sql))
		if err != nil {
			log.Println("Error executing SQL file:", sqlFile, "-", err.Error())
			continue
		}

		fmt.Println("Successfully executed SQL file:", sqlFile)
		// }
	}

	fmt.Println("All SQL files executed successfully!")
}

// Helper function to get the list of SQL files in a directory
func getSQLFiles(directory string) ([]string, error) {
	var sqlFiles []string

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	return sqlFiles, nil
}
