package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var dbPath string = "test.db"
var DsnSuffix string = "?_fk=true&_busy_timeout=30000"

var CreateUserSQL string = "create table if not exists User(UID integer not null primary key, Name text not null, Password text, UserGroup text);"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	var dbconnection *sql.DB
	dbconnection, err = sql.Open("sqlite3", fmt.Sprintf("%s%s", dbPath, DsnSuffix))
	if err != nil {
		panic(err)
	}

	stmt, err := dbconnection.Prepare("INSERT INTO User(UID, Name, Password, UserGroup) values(?,?,?,?)")
	checkError(err)

	res, err := stmt.Exec(1003, "hxia", "nokia123", "0")
	checkError(err)

	id, err := res.LastInsertId()
	checkError(err)

	fmt.Println(id)
}
