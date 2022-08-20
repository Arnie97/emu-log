// Package models defines the database access objects.
package models

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/arnie97/emu-log/common"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "emu-log.db"
)

type Table interface {
	Schema() string
}

var (
	tables []Table
	dbConn *sql.DB
	dbOnce sync.Once
)

func Register(t Table) {
	tables = append(tables, t)
}

// Migrate creates schemas for all tables.
func Migrate(db *sql.DB) (err error) {
	for _, t := range tables {
		if _, err = db.Exec(t.Schema()); err != nil {
			return err
		}
	}
	return nil
}

// DB returns an initialized database singleton.
func DB() *sql.DB {
	dbOnce.Do(func() {
		var err error
		dbConn, err = sql.Open(
			"sqlite3",
			filepath.Join(common.AppPath(), dbFile),
		)
		common.Must(err)
		// TODO: defer dbConn.Close()

		common.Must(Migrate(dbConn))
	})

	return dbConn
}

// CountRecords takes a tableName and returns the number of rows in the table.
func CountRecords(tableName string, fields ...string) (count int) {
	field := "*"
	if len(fields) != 0 {
		field = fields[0]
	}
	row := DB().QueryRow(fmt.Sprintf(
		`SELECT COUNT(%s) FROM %s`, field, tableName,
	))
	common.Must(row.Scan(&count))
	return
}

// In fills the SQL query statement with IN(args...) if args is not empty,
// or with emptyDefault otherwise.
func In(query string, args []string, emptyDefault string) (
	rewritten string, argv []interface{}) {

	if len(args) == 0 {
		rewritten = fmt.Sprintf(query, emptyDefault)
		return
	}
	placeholders := strings.Repeat(", ?", len(args))[2:]
	rewritten = fmt.Sprintf(query, fmt.Sprintf("IN (%s)", placeholders))

	argv = make([]interface{}, len(args))
	for i := range args {
		argv[i] = args[i]
	}
	return
}
