/*
	Trend Micro Deep Discovery Analyzer API SDK
	(c) 2021-2025 by Mikhail Kondrashin (mkondrashin@gmail.com)

	cache.go

	Cached Vision One Sandbox results
*/

package vone

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

//const layout = "2006-01-02T15:04:05Z"

// Cache - database to cache Analyzer check results
type Cache struct {
	dbPath string
	db     *sql.DB
}

// NewCache - open existing or create new cache
func NewCache(db *sql.DB, dbPath string) (*Cache, error) {
	stmt := `CREATE TABLE IF NOT EXISTS hashes (
		type TEXT NOT NULL,
		md5 TEXT NOT NULL UNIQUE,
		sha1 TEXT NOT NULL UNIQUE,
		sha256 TEXT NOT NULL UNIQUE,
		arguments TEXT,
		AnalysisCompletionDateTime TEXT NOT NULL,
		RiskLevel INTEGER,
		DetectionNames TEXT,
		ThreatTypes TEXT,
		TrueFileType TEXT,
		updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS hash_idx ON hashes (sha1);
		`
	if _, err := db.Exec(stmt); err != nil {
		return nil, fmt.Errorf("%s: %w", dbPath, err)
	}
	return &Cache{
		dbPath: dbPath,
		db:     db,
	}, nil
}

// Add - add Analyzer check result to cache database
func (c *Cache) Add(ctx context.Context, data *SandboxAnalysisResultsResponse) error {
	stmt := `INSERT OR REPLACE INTO hashes (
		type,
		md5,
		sha1,
		sha256,
		arguments,
		AnalysisCompletionDateTime,
		RiskLevel,
		DetectionNames,
		ThreatTypes,
		TrueFileType 
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := c.db.ExecContext(ctx, stmt,
		data.Type,
		strings.ToUpper(data.Digest.MD5),
		strings.ToUpper(data.Digest.SHA1),
		strings.ToUpper(data.Digest.SHA256),
		data.Arguments,
		data.AnalysisCompletionDateTime.String(),
		data.RiskLevel,
		strings.Join(data.DetectionNames, ","),
		strings.Join(data.ThreatTypes, ","),
		data.TrueFileType)

	if err != nil && strings.Contains(err.Error(), "pq: syntax error") {
		// its postresql
		stmt := `INSERT INTO hashes (type,
		md5,
		sha1,
		sha256,
		arguments,
		AnalysisCompletionDateTime,
		RiskLevel,
		DetectionNames,
		ThreatTypes,
		TrueFileType
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) ON CONFLICT (sha1) DO UPDATE SET
		arguments=$5,
		AnalysisCompletionDateTime=$6,
		RiskLevel=$7,
		DetectionNames=$8,
		ThreatTypes=$9,
		TrueFileType=$10`
		_, err = c.db.Exec(stmt,
			data.Type,
			strings.ToUpper(data.Digest.MD5),
			strings.ToUpper(data.Digest.SHA1),
			strings.ToUpper(data.Digest.SHA256),
			data.Arguments,
			data.AnalysisCompletionDateTime.String(),
			data.RiskLevel,
			strings.Join(data.DetectionNames, ","),
			strings.Join(data.ThreatTypes, ","),
			data.TrueFileType,
		)
	}
	return c.error("Add Exec", err)
}

// Delete - delete entity from hashes table
func (c *Cache) Delete(sha1 string) error {
	stmt := "DELETE FROM hashes where sha1=$1"
	_, err := c.db.Exec(stmt, strings.ToUpper(sha1))
	return c.error("Delete Exec", err)
}

// Cleanup - remove data from cache that was put there before
// time provided
func (c *Cache) Cleanup(ctx context.Context, date time.Time) error {
	//	stmt := "DELETE FROM hashes where strftime('%s', updated) < $1"
	//	_, err := c.DB.Exec(stmt, fmt.Sprint(date.Unix()))
	stmt := "DELETE FROM hashes where updated < $1"
	if _, err := c.db.ExecContext(ctx, stmt, date); err != nil {
		return c.error("Cleanup hashes", err)
	}
	return nil
}

var ErrNotFound = errors.New("not found")

// Query - get cached Analyzer check result for SHA1 of file
func (c *Cache) Query(ctx context.Context, sha1 string) (*SandboxAnalysisResultsResponse, time.Time, error) {
	stmt := `SELECT type,
		md5,
		sha1,
		sha256,
		arguments,
		AnalysisCompletionDateTime,
		RiskLevel,
		DetectionNames,
		ThreatTypes,
		TrueFileType,
		updated FROM hashes WHERE sha1=$1`
	rows, err := c.db.Query(stmt, strings.ToUpper(sha1))
	if err != nil {
		return nil, time.Time{}, c.error("Query", err)
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, time.Time{}, c.error("Query", rows.Err())
	}
	if !rows.Next() {
		return nil, time.Time{}, nil // fmt.Errorf("Query %s: %w", sha1, ErrNotFound)
	}
	return c.ScanSandboxAnalysisResultsResponse(rows)
}

func (c *Cache) ScanSandboxAnalysisResultsResponse(rows *sql.Rows) (*SandboxAnalysisResultsResponse, time.Time, error) {
	var data SandboxAnalysisResultsResponse
	var detectionNames, threatTypes string
	var analysisCompletionDateTime string
	var updated string
	err := rows.Scan(&data.Type, &data.Digest.MD5, &data.Digest.SHA1, &data.Digest.SHA256,
		&data.Arguments, &analysisCompletionDateTime, &data.RiskLevel, &detectionNames,
		&threatTypes, &data.TrueFileType, &updated)
	if err != nil {
		return nil, time.Time{}, c.error("Query row.Scan", err)
	}
	t, err := time.Parse(timeFormat, analysisCompletionDateTime)
	if err != nil {
		t, err = time.Parse(timeFormatZ, analysisCompletionDateTime)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("ScanSandboxAnalysisResults Parse \"%s\": %w",
				analysisCompletionDateTime, err)
		}
	}
	data.AnalysisCompletionDateTime = VisionOneTime(t)

	data.DetectionNames = strings.Split(detectionNames, ",")
	data.ThreatTypes = strings.Split(threatTypes, ",")
	date, err := time.Parse(timeFormatZ, updated)
	if err != nil {
		return nil, time.Time{}, c.error("IterateCache time.Parse \""+updated+"\"", err)
	}
	return &data, date, nil
}

// Count - return number of entities in cache database
func (c *Cache) Count(ctx context.Context) (int, error) {
	return c.countEntities(ctx, "hashes")
}

// сountEntities - count entities in given table
func (c *Cache) countEntities(ctx context.Context, table string) (int, error) {
	Select := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
	rows, err := c.db.Query(Select)
	if err != nil {
		return -1, c.error("countEntities", err)
	}
	if rows.Err() != nil {
		return 0, c.error("countEntities", rows.Err())
	}
	defer rows.Close()
	if !rows.Next() {
		return 0, c.error("countEntities", fmt.Errorf("rows.Next() returned false"))
	}
	var value int
	err = rows.Scan(&value)
	if err != nil {
		return -1, err
	}
	return value, nil
}

// Close database - should be called when database is not in use anymore
func (c *Cache) Close() error {
	return c.db.Close()
}

// IterateCache - perform provided function fo each database entity
func (c *Cache) IterateCache(ctx context.Context, f func(data *SandboxAnalysisResultsResponse, updated time.Time) error) error {
	Select := `SELECT
	    type,
		md5,
		sha1,
		sha256,
		arguments,
		AnalysisCompletionDateTime,
		RiskLevel,
		DetectionNames,
		ThreatTypes,
		TrueFileType,
		updated
		FROM hashes`
	rows, err := c.db.QueryContext(ctx, Select)
	if err != nil {
		return c.error("IterateCache Query", err)
	}
	defer rows.Close()
	if rows.Err() != nil {
		return c.error("IterateCache Query", rows.Err())
	}
	for rows.Next() {
		data, updated, err := c.ScanSandboxAnalysisResultsResponse(rows)
		if err != nil {
			return c.error("IterateCache", err)
		}
		err = f(data, updated)
		if err != nil {
			return c.error("IterateCache callback", err)
		}
	}
	return nil
}

func (c *Cache) error(message string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %s: %w", c.dbPath, message, err)
}
