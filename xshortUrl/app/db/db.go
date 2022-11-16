package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var MyDB *sql.DB

type Db struct {
	username string
	password string
	address string
	dbname string
}

func NewDB(username, pwd, addr, dbname string) *Db {
	return &Db{
		username: username,
		password: pwd,
		address: addr,
		dbname: dbname,
	}
}

func InitConn(d *Db) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", d.username, d.password, d.address, d.dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("connect mysql error: %w", err))
	}

	// 判断数据库是否连接成功
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("ping mysql server error: %w", err))
	}

	MyDB = db
}
