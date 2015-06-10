package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
)

type postContext struct {
	title string
}

type blogHandler struct {
	*postContext
	h func(*postContext, http.ResponseWriter, *http.Request) error
}

func (bh blogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := bh.h(bh.postContext, w, r)
	if err != nil {
		log.Fatalln(err)
	}
}

var db *sql.DB

var err = errors.New("Open database fail")

func init() {
	db, err = sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		log.Fatalln(err)
	}
	exist := `select * from blog`
	_, err = db.Exec(exist)
	if err != nil {
		sqlStmt := `CREATE TABLE blog (id INTEGER NOT NULL PRIMARY KEY, title TEXT NOT NULL, created TIMESTAMP, body BLOB);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func main() {
	// Intialize database file or table.
	// db.InitDB()

	//authenticator := httpauth.NewBasicAuthenticator("localhost", core.Secret)

	bctx := &postContext{title: ""}
	fs := http.FileServer(http.Dir("static"))

	r := mux.NewRouter()
	r.Handle("/", blogHandler{bctx, homePage})
	r.Handle("/admin", blogHandler{bctx, adminPage})
	r.Handle("/cleanup", blogHandler{bctx, cleanUp})
	r.Handle("/blogs", blogHandler{bctx, listPosts})
	r.Handle("/gallerys", blogHandler{bctx, gallerys})
	r.Handle("/blog/{title}", blogHandler{bctx, viewPost})
	r.Handle("/blog/new/", blogHandler{bctx, newPost})
	r.Handle("/blog/save/", blogHandler{bctx, savePost})
	r.Handle("/blog/update/{title}", blogHandler{bctx, updatePost})
	r.Handle("/blog/saveupdate/", blogHandler{bctx, saveUpdate})
	r.Handle("/blogs/manage/", blogHandler{bctx, managePosts})
	r.Handle("/blog/delete/{title}", blogHandler{bctx, deletePost})
	r.Handle("/static/", http.StripPrefix("/static", fs))
	/*
		n := negroni.New()
		n.Use(auth.Basic("admin", "hello"))
		n.UseHandler(r)
		n.Run(":8008")
		r.Handle("/", n)
	*/
	authHandler := httpauth.SimpleBasicAuth("admin", "hello")
	http.Handle("/", authHandler(r))
	http.ListenAndServe(":8080", nil)
}
