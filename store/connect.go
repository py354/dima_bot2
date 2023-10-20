package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type DataBase struct {
	*sql.DB
}

var DB DataBase

func Connect(user, password, dbname, addr string) {
	connectURL := "postgres://" + user + ":" + password + "@" + addr + "/" + dbname
	log.Println(connectURL)
	db_, err := sql.Open("postgres", connectURL)

	if err != nil {
		panic(err)
	}

	err = db_.Ping()
	if err != nil {
		panic(err)
	}

	DB = DataBase{db_}
	DB.SetMaxOpenConns(5)
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
