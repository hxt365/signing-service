package db

import (
	"database/sql"
	"fmt"

	"coordinator/utils"
	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(idleConn, maxConn int) (*sql.DB, error) {
	dbHost := utils.MustEnv("DATABASE_HOST")
	dbPort := utils.MustEnv("DATABASE_PORT")
	dbName := utils.MustEnv("DATABASE_NAME")
	user := utils.MustEnv("DATABASE_USER")
	pwd := utils.MustEnv("DATABASE_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pwd, dbHost, dbPort, dbName)

	return newDB("mysql", dsn, idleConn, maxConn)
}

func newDB(dialect, dsn string, idleConn, maxConn int) (*sql.DB, error) {
	db, err := sql.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(idleConn)
	db.SetMaxOpenConns(maxConn)

	return db, nil
}
