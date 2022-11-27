package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//var reateUserinfoTable string = ""

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/hello")
	checkError(err)

	//db.Exec()

	stmt, err := db.Prepare("INSERT INTO userinfo SET username=?,department=?,created=?")
	checkError(err)

	res, err := stmt.Exec("hxia", "cloudRAN", "2018-04-10")
	checkError(err)

	id, err := res.LastInsertId()
	checkError(err)

	fmt.Println(id)
}
