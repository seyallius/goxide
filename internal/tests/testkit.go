// Copyright (c) 2025 SeyedAli
// Licensed under the MIT License. See LICENSE file in the project root for details.

// Package tests. testkit.go - Bridges snapdb's generic lifecycle with
// goxide's specific schema and engine (database/sql).
package tests

import (
	"database/sql"
	"fmt"
	"sync"
	"testing"

	"github.com/seyallius/snapdb"
	"github.com/seyallius/snapdb/drivers/sqlite"
	_ "modernc.org/sqlite"
)

// ---------------------------------- Types, Variables & Constants ---------------------------------- //

var (
	engineInitOnce sync.Once
	testDB         *sql.DB
)

// -------------------------------------------- Public API ------------------------------------------ //

// RunGoxideTestMain sets up the SQLite test environment and runs the suite.
func RunGoxideTestMain(m *testing.M, opts ...snapdb.Option) {
	baseOpts := []snapdb.Option{
		snapdb.WithTestdataDir("testdata"),
		// snapdb.WithSQLitePath("testdata/snapdb.sqlite"),
	}
	baseOpts = append(baseOpts, opts...)

	snapdb.Run(m, sqlite.New(), GoxideSchemaInit, GoxideDataInit, GoxideEngineInit, baseOpts...)
}

// DB returns the live *sql.DB, or nil before setup completes.
func DB() *sql.DB { return testDB }

// -------------------------------------------- Lifecycle Callbacks --------------------------------- //

// GoxideSchemaInit maps to snapdb.SchemaInitializer. env.Engine() is nil here.
func GoxideSchemaInit(env *snapdb.Environment) error {
	db, err := sql.Open("sqlite", env.DSN())
	if err != nil {
		return fmt.Errorf("schema init: open db: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			email      TEXT NOT NULL UNIQUE,
			name       TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("schema init: create table: %w", err)
	}
	return nil
}

// GoxideDataInit maps to snapdb.DataInitializer.
func GoxideDataInit(_ *snapdb.Environment) error {
	return nil
}

// GoxideEngineInit maps to snapdb.EngineInitializer.
func GoxideEngineInit(env *snapdb.Environment) (snapdb.Engine, error) {
	var eng snapdb.Engine
	var err error

	engineInitOnce.Do(func() {
		var db *sql.DB
		db, err = sql.Open("sqlite", env.DSN())
		if err != nil {
			return
		}
		db.SetMaxOpenConns(1)
		testDB = db
		eng = &sqlEngine{db: db}
	})
	return eng, err
}

// ------------------------------------------- Engine Adapter -------------------------------------- //

// sqlEngine wraps *sql.DB to satisfy snapdb.Engine.
type sqlEngine struct{ db *sql.DB }

func (e *sqlEngine) Exec(q string, args ...any) (sql.Result, error) { return e.db.Exec(q, args...) }

func (e *sqlEngine) QueryString(q string, args ...any) ([]map[string]string, error) {
	rows, err := e.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var out []map[string]string
	for rows.Next() {
		raw := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range raw {
			ptrs[i] = &raw[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]string, len(cols))
		for i, c := range cols {
			if raw[i].Valid {
				row[c] = raw[i].String
			} else {
				row[c] = ""
			}
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (e *sqlEngine) Ping() error       { return e.db.Ping() }
func (e *sqlEngine) Close() error      { return e.db.Close() }
func (e *sqlEngine) ClearCache() error { return nil }
