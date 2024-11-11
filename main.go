package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var db *sql.DB

type Entry struct {
	date    time.Time
	content string
}

//implementing list.item interface

func (e Entry) FilterValue() string {
	return e.content
}

func (e Entry) Title() string {
	return e.date.Format("Mon Jan _2")
}

func (e Entry) Description() string {
	return e.content
}

var journal *Model
var counter int

func main() {
	var err error
	// createFile(db)
	db, err = sql.Open("sqlite", "example.db")
	createTable(db)
	defer db.Close()
	// prefillYearData(db)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	journal = New()
	p := tea.NewProgram(journal)
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
