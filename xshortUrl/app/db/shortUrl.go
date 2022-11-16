package db

import (
	"database/sql"
)

type ShortUrl struct {
	Uid       int64
	ShortCode string
	UrlStr    string
	Time     	int64 
}

func (su *ShortUrl) Insert(db *sql.DB) error {
	sql := "insert into shorturl(shortcode, urlstr, time) values(?,?,?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// time.Now().UnixMilli()
	_, err = stmt.Exec(su.ShortCode, su.UrlStr, su.Time)
	if err != nil {
		return err
	}
	return nil
}

func (su *ShortUrl) Select(db *sql.DB) error {
	sql := "select urlstr from shorturl where time = ?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRow(su.Time)
	err = row.Scan(&su.UrlStr)
	if err != nil {
		return err
	}
	return nil
}