package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
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
	settings, err := LoadAppSettings()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(settings.URL)
	target := getStrippingTargetPrefix(settings.URL)
	hub := *newSSEHub()
	go hub.run()

	if len(target) > 0 {
		http.Handle(fmt.Sprintf("/%s/css/", target), http.StripPrefix(fmt.Sprintf("/%s", target), http.FileServer(http.Dir("templates"))))
		http.Handle(fmt.Sprintf("/%s/js/", target), http.StripPrefix(fmt.Sprintf("/%s", target), http.FileServer(http.Dir("templates"))))

		http.HandleFunc(fmt.Sprintf("/%s/sse/message", target), func(w http.ResponseWriter, r *http.Request) {
			sendSSEMessage(w, r, &hub)
		})
		http.HandleFunc(fmt.Sprintf("/%s/sse/", target), func(w http.ResponseWriter, r *http.Request) {
			registerSSEClient(w, r, &hub)
		})
	} else {
		http.Handle("/css/", http.FileServer(http.Dir("templates")))
		http.Handle("/js/", http.FileServer(http.Dir("templates")))
		http.HandleFunc("/sse/message", func(w http.ResponseWriter, r *http.Request) {
			sendSSEMessage(w, r, &hub)
		})
		http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
			registerSSEClient(w, r, &hub)
		})
	}
	http.Handle("/", &templateHandler{filename: "index.html", serverUrl: settings.URL})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
func getStrippingTargetPrefix(url string) string {
	sURL := strings.Split(url, "/")
	if len(sURL) <= 3 {
		return ""
	}
	for i := len(sURL) - 1; i >= 3; i-- {
		if sURL[i] != "" {
			return sURL[i]
		}
	}
	return ""
}
