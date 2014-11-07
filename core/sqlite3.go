package core

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ouyanggh/goblog/models"
)

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func InitSqlite3DB() {
	os.Remove("./sqlite3.db")

	db, err := sql.Open("sqlite3", "./sqlite3.db")
	LogFatal(err)
	defer db.Close()

	sqlStmt := `CREATE TABLE blog (id INTEGER NOT NULL PRIMARY KEY, title TEXT NOT NULL, body BLOB);`
	_, err = db.Exec(sqlStmt)
	LogFatal(err)
}

func SqliteInsert(p *models.Post) {
	now := time.Now().Unix()
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	LogFatal(err)

	tx, err := db.Begin()
	LogFatal(err)

	stmtq, err := db.Prepare("SELECT id FROM blog WHERE id = ?")
	LogFatal(err)
	defer stmtq.Close()

	var nid int64
	err = stmtq.QueryRow(now).Scan(&nid)
	if err == nil {
		now = now + 1
	}

	stmt, err := tx.Prepare("INSERT INTO blog VALUES(?, ?, ?)")
	LogFatal(err)
	defer stmt.Close()

	_, err = stmt.Exec(now, p.Title, p.Body)
	LogFatal(err)
	tx.Commit()
}

func SqliteDelete(title string) {
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	LogFatal(err)

	tx, err := db.Begin()
	LogFatal(err)

	stmt, err := tx.Prepare("DELETE FROM blog WHERE title = ?")
	LogFatal(err)
	defer stmt.Close()

	_, err = stmt.Exec(title)
	LogFatal(err)
	tx.Commit()
}

func SqliteUpdate(np *models.Post, title string) {
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	LogFatal(err)

	tx, err := db.Begin()
	LogFatal(err)

	stmt, err := tx.Prepare("UPDATE blog SET title = ?, body = ? WHERE title = ?")
	LogFatal(err)
	defer stmt.Close()

	_, err = stmt.Exec(np.Title, np.Body, title)
	LogFatal(err)
	tx.Commit()
}

func SqliteQuery(title string) string {
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	LogFatal(err)

	stmt, err := db.Prepare("SELECT body FROM blog WHERE title = ?")
	LogFatal(err)
	defer stmt.Close()

	var body string
	err = stmt.QueryRow(title).Scan(&body)
	LogFatal(err)

	return body
}

func SqliteQueryAll() (titles map[string][]byte) {
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	LogFatal(err)

	rows, err := db.Query("SELECT title, body FROM blog")
	LogFatal(err)
	defer rows.Close()
	titles = make(map[string][]byte)
	for rows.Next() {
		var title string
		var body []byte
		rows.Scan(&title, &body)
		titles[title] = body
	}
	return titles
}
