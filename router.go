package main

import (
	"net/http"
	"time"
)

func homePage(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	p := &page{Tmpl: "layout", Post: &post{}, W: w}
	err := p.renderTemplate()
	return err
}

func newPost(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	p := &page{Tmpl: "new", Post: &post{}, W: w}
	err := p.renderTemplate()
	return err
}

func savePost(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	rtitle := r.FormValue("title")
	rbody := r.FormValue("body")
	now := time.Now()
	p := &post{
		Title:   rtitle,
		Created: now,
		Body:    []byte(rbody),
	}
	err := p.insert()
	http.Redirect(w, r, "/blog/"+rtitle, http.StatusFound)
	return err
}

func updatePost(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	title := r.URL.Path[len("/blog/update/"):]

	btx.title = title
	p := &post{Title: title}
	p, err := p.query()

	pa := &page{Tmpl: "edit", Post: p, W: w}
	err = pa.renderTemplate()
	return err
}

func viewPost(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	title := r.URL.Path[len("/blog/"):]
	p := &post{Title: title}
	p, err := p.query()

	pa := &page{Tmpl: "view", Post: p, W: w}
	err = pa.renderTemplate()
	return err
}

func saveUpdate(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	rtitle := r.FormValue("title")
	rbody := r.FormValue("body")
	now := time.Now()
	p := &post{
		Title:   rtitle,
		Created: now,
		Body:    []byte(rbody),
	}
	err := p.update(btx.title)
	http.Redirect(w, r, "/blog/"+rtitle, http.StatusFound)
	return err
}

func listPosts(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	var p blogs
	p.Posts = getAllPosts()

	pa := &page{Tmpl: "lists", Post: p, W: w}
	err := pa.renderTemplate()
	return err
}

func managePosts(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	var p blogs
	p.Posts = getAllPosts()

	pa := &page{Tmpl: "exists", Post: p, W: w}
	err := pa.renderTemplate()
	return err
}

func deletePost(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	title := r.URL.Path[len("/blog/delete/"):]
	var p post
	p.Title = title
	err := p.delete()
	http.Redirect(w, r, "/blogs/manage/", http.StatusFound)
	return err
}

// cleanup by delete database file and initialize it again
// all exist data will be lost
func cleanUp(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	err := cleanup()
	http.Redirect(w, r, "/blogs/manage/", http.StatusFound)
	return err
}

func gallerys(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	p := &page{Tmpl: "gallerys", Post: &post{}, W: w}
	err = p.renderTemplate()
	return err
}

func adminPage(btx *postContext, w http.ResponseWriter, r *http.Request) error {
	p := &page{Tmpl: "admin", Post: &post{}, W: w}
	err = p.renderTemplate()
	return err
}