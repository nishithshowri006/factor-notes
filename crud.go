package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type EntryDB struct {
	Entry
	day string

	id      int
	weekNum int
}

func createFile(db *sql.DB) {
	os.Remove("example.db")
	var err error
	db, err = sql.Open("sqlite", "example.db")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func createTable(db *sql.DB) {
	crtSmt := `CREATE TABLE IF NOT EXISTS journal(
	id INTEGER PRIMARY KEY,
	date DATETIME ,
	weekNum INTEGER NOT NULL,
	content TEXT,
	day VARCHAR(8) NOT NULL
	)`
	_, err := db.Exec(crtSmt)
	if err != nil {
		log.Fatal(err)
	}
}

func prefillYearData(db *sql.DB) {
	// var entries []EntryDB
	vals := make([]interface{}, 366*4)
	today := time.Now()
	midnight := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.Local)
	sqlSmt := `Insert INTO 
	journal(id, date, weekNum, day)
	VALUES 
	`
	for i, j := 0, 0; i < 365*4; i += 4 {
		dt := midnight.AddDate(0, 0, j)
		_, week := dt.ISOWeek()
		// entry := EntryDB{id: dt.YearDay(), weekNum: week, day: dt.Weekday().String(), Entry: Entry{date: dt}}
		vals[i], vals[i+1], vals[i+2], vals[i+3] = dt.YearDay(), dt, week, dt.Weekday().String()
		sqlSmt += `(?,?,?,?),`
		j++
	}
	sqlSmt = sqlSmt[:len(sqlSmt)-1]
	_, err := db.Exec(sqlSmt, vals...)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateValue(db *sql.DB, e Entry) {
	sqlSmt := `Update journal
	SET content = ?
	where id = ?
	`
	_, err := db.Exec(sqlSmt, e.content, e.date.YearDay())
	if err != nil {
		log.Fatal(err)
	}
}

func getWeekData(db *sql.DB, date time.Time, counter *int) ([]Entry, error) {
	_, week := date.ISOWeek()
	sqlSmt := "select date,IFNULL(content,'') from journal where weekNum = ?"
	var entries []Entry
	rows, err := db.Query(sqlSmt, week+*counter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.date, &entry.content); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
