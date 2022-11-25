package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const maxRecords = 100000
const maxKeys = 100
const batchSize = 1000

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("could not read env %s", key)
	}
	return val
}

func newMySQL(idleConn, maxConn int) (*sql.DB, error) {
	dbHost := mustEnv("DATABASE_HOST")
	dbPort := mustEnv("DATABASE_PORT")
	dbName := mustEnv("DATABASE_NAME")
	user := mustEnv("DATABASE_USER")
	pwd := mustEnv("DATABASE_PASSWORD")
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

type genFunc func(tx *sql.Tx, size int) error

func genBatch(tx *sql.Tx, size int, genFunc genFunc) error {
	batchNum := int(math.Ceil(float64(size) / batchSize))
	for i := 1; i < batchNum; i++ {
		if err := genFunc(tx, batchSize); err != nil {
			return err
		}
	}

	remainingBatchNum := size - (batchNum-1)*batchSize
	if err := genFunc(tx, remainingBatchNum); err != nil {
		return err
	}
	return nil
}

func genRecords(tx *sql.Tx, size int) error {
	valueStrings := make([]string, 0, size)
	valueArgs := make([]interface{}, 0, size)
	for i := 1; i <= size; i++ {
		valueStrings = append(valueStrings, "(?)")
		valueArgs = append(valueArgs, randSeq(1024))
	}

	stmt := fmt.Sprintf("INSERT INTO records (value) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := tx.Exec(stmt, valueArgs...); err != nil {
		return err
	}
	return nil
}

func genKeys(tx *sql.Tx, size int) error {
	valueStrings := make([]string, 0, size)
	valueArgs := make([]interface{}, 0, size)
	for i := 1; i <= size; i++ {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, i, randSeq(255))
	}

	stmt := fmt.Sprintf("INSERT INTO secret_keys (identifier, value) VALUES %s", strings.Join(valueStrings, ","))
	if _, err := tx.Exec(stmt, valueArgs...); err != nil {
		return err
	}
	return nil
}

func deleteAll(tx *sql.Tx) error {
	if _, err := tx.Exec("DELETE FROM secret_keys"); err != nil {
		return err
	}
	if _, err := tx.Exec("DROP TABLE IF EXISTS records"); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM signatures"); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM progresses"); err != nil {
		return err
	}
	if _, err := tx.Exec(`CREATE TABLE records
	(
		id    INT AUTO_INCREMENT,
		value VARCHAR(1024),
		PRIMARY KEY (id)
	)`); err != nil {
		return err
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db, err := newMySQL(10, 10)
	if err != nil {
		log.Fatal("could not connect to db", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("could not start a transaction", err)
	}
	defer tx.Rollback()

	if err := deleteAll(tx); err != nil {
		log.Fatal("could not delete data", err)
	}

	if err := genBatch(tx, maxRecords, genRecords); err != nil {
		log.Fatal("could not create records", err)
	}

	if err := genBatch(tx, maxKeys, genKeys); err != nil {
		log.Fatal("could not create keys", err)
	}

	tx.Commit()

	log.Println("Done!")
}
