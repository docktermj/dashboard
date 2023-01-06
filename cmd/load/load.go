// The version package simply prints the version of the go-fileindex binary file.
package load

import (
	"bufio"
	"database/sql"
	"encoding/json"
	log "github.com/docktermj/go-logger/logger"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"os"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type FileMetadata struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Volume   string `json:"volume"`
	Size     int64  `json:"size"`
	SHA256   string `json:"sha256"`
	Modified int64  `json:"modified"`
}

// ----------------------------------------------------------------------------
// Helper functions
// ----------------------------------------------------------------------------

func viperAsJson() string {
	viperConfig := viper.AllSettings()
	viperByteArray, err := json.Marshal(viperConfig)
	if err != nil {
		log.Fatalf("Unable to marshal viper config to JSON: %v", err)
	}
	return string(viperByteArray)
}

// ----------------------------------------------------------------------------
// Load
// ----------------------------------------------------------------------------

// Load a SQLite database from a file of JSONlines.
func Load(jsonFileName string, sqliteFileName string) {

	// Prepare SQLite database.

	database, err := sql.Open("sqlite3", sqliteFileName)
	if err != nil {
		log.Errorf("%s cannot be opened. Error:  %v", sqliteFileName, err)
		return
	}

	// Create database table, if it doesn't already exist.

	createSqlStatement := `CREATE TABLE IF NOT EXISTS fileindex (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      path TEXT,
      name TEXT,
      volume TEXT,
      sha256 TEXT,
      size INTEGER,
      modified INTEGER
      )`

	createSql, err := database.Prepare(createSqlStatement)
	if err != nil {
		log.Errorf("Cannot prepare SQL: %s  Error:  %v", createSqlStatement, err)
		return
	}
	_, err = createSql.Exec()
	if err != nil {
		log.Errorf("Cannot exec SQL: %s  Error:  %v", createSqlStatement, err)
		return
	}

	// Prepare SQL for insert.

	insertSqlStatement := "INSERT INTO fileindex (path, name, volume, size, sha256, modified) VALUES (?, ?, ?, ?, ?, ?)"
	insertSql, err := database.Prepare(insertSqlStatement)
	if err != nil {
		log.Errorf("Cannot prepare SQL: %s  Error:  %v", insertSqlStatement, err)
		return
	}

	// Open input file.

	jsonFile, err := os.Open(jsonFileName)
	defer jsonFile.Close()
	if err != nil {
		log.Errorf("%s cannot be opened. Error:  %v", jsonFile, err)
		return
	}

	// Read JSONlines.

	scanner := bufio.NewScanner(jsonFile)
	for scanner.Scan() {

		// Read JSONLine into FileMetadata structure.

		var fileMetadata FileMetadata
		if err := json.Unmarshal(scanner.Bytes(), &fileMetadata); err != nil {
			log.Errorf("Unmarshal error: %v", err)
		}

		// Insert into database.

		insertSql.Exec(
			fileMetadata.Path,
			fileMetadata.Name,
			fileMetadata.Volume,
			fileMetadata.Size,
			fileMetadata.SHA256,
			fileMetadata.Modified,
		)
	}
	if scanner.Err() != nil {
		log.Errorf("Scanner error: %v", scanner.Err())
	}

}

// ----------------------------------------------------------------------------
// Command pattern "Execute" function.
// ----------------------------------------------------------------------------

// The Command sofware design pattern's Execute() method.
func Execute() {

	// Print context parameters.
	// TODO: Formalize entry parameters

	log.Info(viperAsJson())

	// Get parameters from viper.

	var jsonFileName = viper.GetString("json_file_name")
	var sqliteFileName = viper.GetString("sqlite_file_name")

	// Check parameters.

	errors := 0

	if jsonFileName == "" {
		errors += 1
		log.Warn("--json-file-name not set.")
	}

	if sqliteFileName == "" {
		errors += 1
		log.Warn("--sqlite-file-name not set.")
	}

	// If any errors, exit.

	if errors > 0 {
		return
	}

	// Perform command.

	Load(jsonFileName, sqliteFileName)

}
