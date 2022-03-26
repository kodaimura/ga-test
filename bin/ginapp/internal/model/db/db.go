package db

import (
    "log"
    "database/sql"
    
    _ "github.com/mattn/go-sqlite3"

    "ginapp/internal/constants"
)


var db *sql.DB

func init() {
    var err error

    dbName := "./" + constants.AppName + ".db"
    db, err = sql.Open("sqlite3", dbName)

    if err != nil {
        log.Panic(err)
    }
}

func GetDB() *sql.DB {
    return db
}
