package dbinit

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

const (
	maxDBConnections  = 10
	maxPingDBAttempts = 20
)

func GetPostgres() (*sql.DB, error) {
	pass := os.Getenv("pass")
	user := os.Getenv("user")
	dbName := os.Getenv("dbName")
	host := os.Getenv("hostPG")
	port := os.Getenv("portPG")
	sslMode := os.Getenv("sslMode")
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		user, dbName, pass, host, port, sslMode)
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxDBConnections)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	attemptsNumber := 0
	for range ticker.C {
		err = db.Ping()
		attemptsNumber++
		if err == nil {
			break
		}
		if attemptsNumber == maxPingDBAttempts {
			return nil, err
		}
	}
	return db, nil
}
