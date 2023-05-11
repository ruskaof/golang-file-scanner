package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func SetupDbConnection(username, password, host, port, dbName string) (*sql.DB, error) {
	url := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		username,
		password,
		host,
		port,
		dbName,
	)

	db, err := sql.Open("postgres", url)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}
