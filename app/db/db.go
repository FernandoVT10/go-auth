package db

import (
	"database/sql"

	"github.com/FernandoVT10/go-auth/app/utils"
	_ "github.com/mattn/go-sqlite3"
)

const SQLITE_DB_PATH = "./sqlite.db"

var db *sql.DB

func Connect() error {
    _db, err := sql.Open("sqlite3", SQLITE_DB_PATH)
    if err != nil {
        return err
    }
    db = _db;
    return nil
}

func Initialize() error {
    stmt := `
    CREATE TABLE IF NOT EXISTS Users (
        Id INTEGER NOT NULL PRIMARY KEY,
        Username VARCHAR(20),
        Password VARCHAR(100)
    )
    `
    _, err := db.Exec(stmt)
    return err
}

func Close() {
    if db != nil {
        db.Close()
    }
}

// These are only wrappers to the actual functions
func Exec(query string, args...any) (sql.Result, error) {
    if db == nil {
        utils.LogFatal("db is nil")
    }

    return db.Exec(query, args...)
}

func QueryRow(query string, args...any) *sql.Row {
    if db == nil {
        utils.LogFatal("db is nil")
    }

    return db.QueryRow(query, args...)
}
