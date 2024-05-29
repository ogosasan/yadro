package repository

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
	"strings"
	"yadro/internal/core/comics"
)

func Head(db_path string, comicsMap map[int]comics.Write, indexMap map[string][]string) {
	file, err := os.Create(db_path)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/adapter/repository/migrations",
		"sqlite3", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	Add(db, comicsMap, indexMap)
}

func Add(db *sql.DB, comicsMap map[int]comics.Write, indexMap map[string][]string) {
	for key, value := range comicsMap {
		records := `INSERT INTO database(id, Keywords, Url) VALUES (?, ?, ?)`
		_, err := db.Exec(records, key, strings.Join(value.Tscript, ","), value.Img)
		if err != nil {
			log.Fatal(err)
		}
	}
	for key, value := range indexMap {
		records := `INSERT INTO index_table(Keywords, Numbers) VALUES (?, ?)`
		_, err := db.Exec(records, key, strings.Join(value, ","))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func FetchRecords(db *sql.DB) (map[int]comics.Write, map[string][]int) {
	rows, err := db.Query("SELECT * FROM database")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	comicsMap := map[int]comics.Write{}
	for rows.Next() {
		var id int
		var keywords string
		var url string
		if err := rows.Scan(&id, &keywords, &url); err != nil {
			log.Fatal(err)
		}
		comicsMap[id] = comics.Write{Tscript: strings.Split(keywords, ","), Img: url}
	}
	rowss, err := db.Query("SELECT * FROM index_table")
	if err != nil {
		log.Fatal(err)
	}
	defer rowss.Close()
	indexMap := map[string][]int{}
	for rowss.Next() {
		var keyword string
		var numbers string
		values := make([]int, 0, len(numbers))
		if err := rowss.Scan(&keyword, &numbers); err != nil {
			log.Fatal(err)
		}
		tmp := strings.Split(numbers, ",")
		for _, raw := range tmp {
			v, _ := strconv.Atoi(raw)
			values = append(values, v)
		}
		indexMap[keyword] = values
	}
	return comicsMap, indexMap
}
