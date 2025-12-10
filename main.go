package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"webappsample/logs"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
	settings AppSettings
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// "sync.Once" executes only one time.
	t.once.Do(func() {
		// "Must()" wraps "ParseFiles()" results, so I can put it into "templateHandler.templ" directly
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, t.settings)
}

func main() {
	err := logs.ConfigLogger()
	if err != nil {
		log.Panicln(err.Error())
	}
	settings, err := LoadAppSettings()
	if err != nil {
		log.Println(err.Error())
	}
	urlPrefix := getStrippingTargetPrefix(settings.URL)
	hub := *newSSEHub()
	go hub.run()

	if len(urlPrefix) > 0 {
		http.Handle(fmt.Sprintf("%s/css/", urlPrefix), http.StripPrefix(fmt.Sprintf("%s", urlPrefix), http.FileServer(http.Dir("templates"))))
		http.Handle(fmt.Sprintf("%s/js/", urlPrefix), http.StripPrefix(fmt.Sprintf("%s", urlPrefix), http.FileServer(http.Dir("templates"))))
	} else {
		http.Handle("/css/", http.FileServer(http.Dir("templates")))
		http.Handle("/js/", http.FileServer(http.Dir("templates")))
	}
	http.HandleFunc(fmt.Sprintf("%s/sse/message", urlPrefix), func(w http.ResponseWriter, r *http.Request) {
		sendSSEMessage(w, r, &hub)
	})
	http.HandleFunc(fmt.Sprintf("%s/sse/", urlPrefix), func(w http.ResponseWriter, r *http.Request) {
		registerSSEClient(w, r, &hub)
	})
	http.Handle("/", &templateHandler{filename: "index.html", settings: settings})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
func getStrippingTargetPrefix(url string) string {
	sURL := strings.Split(url, "/")
	if len(sURL) <= 3 {
		return ""
	}
	for i := len(sURL) - 1; i >= 3; i-- {
		if sURL[i] != "" {
			return fmt.Sprintf("/%s", sURL[i])
		}
	}
	return ""
}
