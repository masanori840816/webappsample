package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// "sync.Once" executes only one time.
	t.once.Do(func() {
		// "Must()" wraps "ParseFiles()" results, so I can put it into "templateHandler.templ" directly
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, "ws://localhost:8080/websocket")
}

func main() {
	http.Handle("/css/", http.FileServer(http.Dir("templates")))
	http.Handle("/js/", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/websocket", websocketHandler)
	// In this sample, "ServeHTTP()" is called twice.
	// The second time is for loading "favicon.ico"
	http.Handle("/", &templateHandler{filename: "index.html"})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
