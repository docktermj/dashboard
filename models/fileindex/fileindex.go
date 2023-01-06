package fileindex

// See https://www.alexedwards.net/blog/organising-database-access "Using an interface"

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const SelectAll = "SELECT id, volume, path, name, modified, size, sha256 FROM fileindex"
const SelectCount = "SELECT COUNT(*) as count FROM fileindex"

// A list of methods for "fileindex"
type Interface interface {
	AllDuplicatesSha256Count() (int, error)
	AllEverythingCount() (int, error)
	AllUniqueSha256Count() (int, error)
	ByID(string, *DatabaseQueryMetadata) ([]*FileIndex, error)
	ByIDCount(string) (int, error)
	ByModified(string, *DatabaseQueryMetadata) ([]*FileIndex, error)
	ByModifiedCount(string) (int, error)
	ByName(string, *DatabaseQueryMetadata) ([]*FileIndex, error)
	ByNameCount(string) (int, error)
	ByPath(string, *DatabaseQueryMetadata) ([]*FileIndex, error)
	ByPathCount(string) (int, error)
	BySha256(string, *DatabaseQueryMetadata) ([]*FileIndex, error)
	BySha256Count(string) (int, error)
	BySize(string, *DatabaseQueryMetadata) ([]*FileIndex, error)
	BySizeCount(string) (int, error)
	ByVolume(string, *DatabaseQueryMetadata) ([]*FileIndex, error)
	ByVolumeCount(string) (int, error)
	DuplicatesSha256(*DatabaseQueryMetadata) ([]*DuplicatesSha256, error)
	DuplicatesSha256Count(*DatabaseQueryMetadata) (int, error)
	Everything(*DatabaseQueryMetadata) ([]*FileIndex, error)
	EverythingCount(*DatabaseQueryMetadata) (int, error)
	UniqueSha256(*DatabaseQueryMetadata) ([]*FileIndex, error)
	UniqueSha256Count(*DatabaseQueryMetadata) (int, error)
}

type DB struct {
	*sql.DB
}

type DatabaseResult struct {
	Total    int64     `json:"total"`
	Filtered int64     `json:"filtered"`
	Data     FileIndex `json:"data"`
}

type FileIndex struct {
	Id       string `json:"id"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Volume   string `json:"volume"`
	Size     int64  `json:"size"`
	SHA256   string `json:"sha256"`
	Modified int64  `json:"modified"`
}

type DuplicatesSha256 struct {
	Count  int64  `json:"count"`
	SHA256 string `json:"sha256"`
}

type DatabaseQueryMetadata struct {
	Limit          int
	Start          int
	OrderColumn    int
	OrderDirection string // asc, desc
	Search         string
}

// ----------------------------------------------------------------------------
// Database requests for count
// ----------------------------------------------------------------------------

func (database *DB) count(sqlStatement string) (int, error) {
	var count int
	row := database.QueryRow(sqlStatement)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (database *DB) byGenericCount(sqlStatement string, sqlParameters ...interface{}) (int, error) {
	var count int
	row := database.QueryRow(sqlStatement, sqlParameters...)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (database *DB) AllEverythingCount() (int, error) {
	return database.count("SELECT COUNT(*) as count FROM fileindex")
}

func (database *DB) AllDuplicatesSha256Count() (int, error) {
	return database.count("SELECT COUNT(*) FROM (SELECT COUNT(*) count, sha256 FROM fileindex GROUP BY sha256 HAVING COUNT(*) > 1)")
}

func (database *DB) AllUniqueSha256Count() (int, error) {
	return database.count("SELECT COUNT(*) FROM (SELECT COUNT(*) count, sha256 FROM fileindex GROUP BY sha256 HAVING COUNT(*) = 1)")
}

// ----------------------------------------------------------------------------
// DuplicatesSha256 and Count
// ----------------------------------------------------------------------------

func (database *DB) DuplicatesSha256(metadata *DatabaseQueryMetadata) ([]*DuplicatesSha256, error) {

	// Determine SQL format string.

	formatString := "SELECT COUNT(*) count, sha256 FROM fileindex %s GROUP BY sha256 HAVING COUNT(*) > 1 ORDER BY %d %s LIMIT ?, ?"

	if metadata.Search != "" {
		formatString = "SELECT COUNT(*) as count, sha256 FROM fileindex WHERE sha256 LIKE '%%%s%%' GROUP BY sha256 HAVING COUNT(*) > 1 ORDER BY %d %s LIMIT ?, ?"
	}

	sqlStatement := fmt.Sprintf(
		formatString,
		metadata.Search,
		metadata.OrderColumn,
		metadata.OrderDirection)

	rows, err := database.Query(sqlStatement, metadata.Start, metadata.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*DuplicatesSha256, 0)
	for rows.Next() {
		result := new(DuplicatesSha256)
		err := rows.Scan(&result.Count, &result.SHA256)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (database *DB) DuplicatesSha256Count(metadata *DatabaseQueryMetadata) (int, error) {

	// Determine SQL format string.

	formatString := "SELECT COUNT(*) as count from (SELECT sha256 FROM fileindex %s GROUP BY sha256 HAVING COUNT(*) > 1)"
	if metadata.Search != "" {
		formatString = "SELECT COUNT(*) as count from (SELECT sha256 FROM fileindex WHERE sha256 LIKE '%%%s%%' GROUP BY sha256 HAVING COUNT(*) > 1)"
	}

	// Run SQL statement.

	sqlStatement := fmt.Sprintf(formatString,
		metadata.Search)

	return database.count(sqlStatement)
}

// ----------------------------------------------------------------------------
// UniqueSha256 and Count
// ----------------------------------------------------------------------------

func (database *DB) UniqueSha256(metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {

	// Determine SQL format string.

	formatString := "%s %s GROUP BY sha256 HAVING COUNT(*) = 1 ORDER BY %d %s LIMIT ?, ?"

	if metadata.Search != "" {
		formatString = "%s WHERE name LIKE '%%%s%%' GROUP BY sha256 HAVING COUNT(*) = 1 ORDER BY %d %s LIMIT ?, ?"
	}

	// Construct SQL statement.

	sqlStatement := fmt.Sprintf(formatString,
		SelectAll,
		metadata.Search,
		metadata.OrderColumn,
		metadata.OrderDirection)

	return database.byGeneric(sqlStatement, metadata.Start, metadata.Limit)
}

func (database *DB) UniqueSha256Count(metadata *DatabaseQueryMetadata) (int, error) {

	// Determine SQL format string.

	formatString := "SELECT COUNT(*) as count FROM (SELECT COUNT(*) as count from (%s %s GROUP BY sha256 HAVING COUNT(*) = 1))"

	if metadata.Search != "" {
		formatString = "SELECT COUNT(*) as count FROM (SELECT COUNT(*) as count from (%s WHERE name LIKE '%%%s%%' GROUP BY sha256 HAVING COUNT(*) = 1))"
	}

	// Run SQL statement.

	sqlStatement := fmt.Sprintf(formatString,
		SelectCount,
		metadata.Search)

	return database.count(sqlStatement)
}

// ----------------------------------------------------------------------------
// Everything and Count
// ----------------------------------------------------------------------------

func (database *DB) Everything(metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {

	// Determine SQL format string.

	formatString := "%s %s ORDER BY %d %s LIMIT ?, ?" // Tricky code: 2nd %s may be replaced with empty string.

	if metadata.Search != "" {
		formatString = "%s WHERE name LIKE '%%%s%%' ORDER BY %d %s LIMIT ?, ?"
	}

	// Run SQL statement.

	sqlStatement := fmt.Sprintf(formatString,
		SelectAll,
		metadata.Search,
		metadata.OrderColumn,
		metadata.OrderDirection)
	rows, err := database.Query(sqlStatement, metadata.Start, metadata.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Transform rows returned to array of FileIndex.

	data := make([]*FileIndex, 0)
	for rows.Next() {
		result := new(FileIndex)
		err := rows.Scan(
			&result.Id,
			&result.Volume,
			&result.Path,
			&result.Name,
			&result.Modified,
			&result.Size,
			&result.SHA256)
		if err != nil {
			return nil, err
		}
		data = append(data, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (database *DB) EverythingCount(metadata *DatabaseQueryMetadata) (int, error) {

	// Determine SQL format string.

	formatString := "%s %s" // Tricky code: 2nd %s may be replaced with empty string.

	if metadata.Search != "" {
		formatString = "%s WHERE name LIKE '%%%s%%'"
	}

	// Run SQL statement.

	sqlStatement := fmt.Sprintf(formatString,
		SelectCount,
		metadata.Search)

	return database.count(sqlStatement)
}

// ----------------------------------------------------------------------------
// ByXXX() methods
//   - All use FileIndex structure.
// ----------------------------------------------------------------------------

func (database *DB) createSqlStatementForByMethods(metadata *DatabaseQueryMetadata, selectClause string) string {

	// Determine SQL format string.

	formatString := "%s %s %s ORDER BY %d %s LIMIT %d, %d" // Tricky code: 3rd %s may be replaced with empty string.

	if metadata.Search != "" {
		formatString = "%s %s AND name LIKE '%%%s%%' ORDER BY %d %s LIMIT %d, %d"
	}

	// Verify metadata.OrderDirection in ["asc", "desc"]

	orderDirection := "asc"
	if metadata.OrderDirection == "desc" {
		orderDirection = metadata.OrderDirection
	}

	// Construct SQL statement.

	sqlStatement := fmt.Sprintf(formatString,
		SelectAll,
		selectClause,
		metadata.Search,
		metadata.OrderColumn,
		orderDirection,
		metadata.Start,
		metadata.Limit,
	)
	return sqlStatement
}

func (database *DB) createSqlStatementForByMethodsCount(selectClause string) string {
	return fmt.Sprintf("SELECT COUNT(*) as count FROM (%s %s)", SelectAll, selectClause)
}

func (database *DB) byGeneric(sqlStatement string, sqlParameters ...interface{}) ([]*FileIndex, error) {

	// Execute database query.

	rows, err := database.Query(sqlStatement, sqlParameters...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Transform the SQL Query results into a FileIndex structure.

	results := make([]*FileIndex, 0)
	for rows.Next() {
		result := new(FileIndex)
		err := rows.Scan(
			&result.Id,
			&result.Volume,
			&result.Path,
			&result.Name,
			&result.Modified,
			&result.Size,
			&result.SHA256,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil

}

func (database *DB) ByID(id string, metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {
	sqlStatement := database.createSqlStatementForByMethods(metadata, "WHERE id=?")
	return database.byGeneric(sqlStatement, id)
}

func (database *DB) ByIDCount(id string) (int, error) {
	sqlStatement := database.createSqlStatementForByMethodsCount("WHERE id=?")
	return database.byGenericCount(sqlStatement, id)
}

func (database *DB) ByModified(modified string, metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {
	sqlStatement := database.createSqlStatementForByMethods(metadata, "WHERE modified=?")
	return database.byGeneric(sqlStatement, modified)
}

func (database *DB) ByModifiedCount(modified string) (int, error) {
	sqlStatement := database.createSqlStatementForByMethodsCount("WHERE modified=?")
	return database.byGenericCount(sqlStatement, modified)
}

func (database *DB) ByName(name string, metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {
	sqlStatement := database.createSqlStatementForByMethods(metadata, "WHERE name=?")
	return database.byGeneric(sqlStatement, name)
}

func (database *DB) ByNameCount(name string) (int, error) {
	sqlStatement := database.createSqlStatementForByMethodsCount("WHERE name=?")
	return database.byGenericCount(sqlStatement, name)
}

func (database *DB) ByPath(path string, metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {
	sqlStatement := database.createSqlStatementForByMethods(metadata, "WHERE path=?")
	return database.byGeneric(sqlStatement, path)
}

func (database *DB) ByPathCount(path string) (int, error) {
	sqlStatement := database.createSqlStatementForByMethodsCount("WHERE path=?")
	return database.byGenericCount(sqlStatement, path)
}

func (database *DB) BySha256(sha256 string, metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {
	sqlStatement := database.createSqlStatementForByMethods(metadata, "WHERE sha256=?")
	return database.byGeneric(sqlStatement, sha256)
}

func (database *DB) BySha256Count(sha256 string) (int, error) {
	sqlStatement := database.createSqlStatementForByMethodsCount("WHERE sha256=?")
	return database.byGenericCount(sqlStatement, sha256)
}

func (database *DB) BySize(size string, metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {
	sqlStatement := database.createSqlStatementForByMethods(metadata, "WHERE size=?")
	return database.byGeneric(sqlStatement, size)
}

func (database *DB) BySizeCount(size string) (int, error) {
	sqlStatement := database.createSqlStatementForByMethodsCount("WHERE size=?")
	return database.byGenericCount(sqlStatement, size)
}

func (database *DB) ByVolume(volume string, metadata *DatabaseQueryMetadata) ([]*FileIndex, error) {
	sqlStatement := database.createSqlStatementForByMethods(metadata, "WHERE volume=?")
	return database.byGeneric(sqlStatement, volume)
}

func (database *DB) ByVolumeCount(volume string) (int, error) {
	sqlStatement := database.createSqlStatementForByMethodsCount("WHERE volume=?")
	return database.byGenericCount(sqlStatement, volume)
}
