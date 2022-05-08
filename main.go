package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once      sync.Once
	filename  string
	templ     *template.Template
	serverUrl string
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// "sync.Once" executes only one time.
	t.once.Do(func() {
		// "Must()" wraps "ParseFiles()" results, so I can put it into "templateHandler.templ" directly
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, t.serverUrl)
}

func main() {
	http.Handle("/css/", http.FileServer(http.Dir("templates")))
	http.Handle("/js/", http.FileServer(http.Dir("templates")))
	http.HandleFunc("/websocket", websocketHandler)
	http.HandleFunc("/sse", sseHandler)
	http.HandleFunc("/webreq", webRequestHandler)

	http.Handle("/pages/sse", &templateHandler{filename: "sse.html", serverUrl: "http://localhost:8080/sse"})
	http.Handle("/pages/webreq", &templateHandler{filename: "webReq.html", serverUrl: "http://localhost:8080/webreq"})
	http.Handle("/", &templateHandler{filename: "index.html", serverUrl: "ws://localhost:8080/websocket"})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
