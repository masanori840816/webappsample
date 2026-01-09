package main

import (
	"encoding/json"
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
	groups := *NewChatGroups()
	defer func() {
		groups.close <- 0
	}()
	go groups.Run()

	http.Handle(fmt.Sprintf("%s/css/", urlPrefix), http.StripPrefix(fmt.Sprintf("%s", urlPrefix), http.FileServer(http.Dir("templates"))))
	http.Handle(fmt.Sprintf("%s/js/", urlPrefix), http.StripPrefix(fmt.Sprintf("%s", urlPrefix), http.FileServer(http.Dir("templates"))))
	http.Handle(fmt.Sprintf("%s/videos/", urlPrefix), http.StripPrefix(fmt.Sprintf("%s", urlPrefix), http.FileServer(http.Dir("templates"))))
	http.HandleFunc(fmt.Sprintf("%s/sse/message", urlPrefix), func(w http.ResponseWriter, r *http.Request) {
		message, err := GetClientMessage(w, r)
		if err != nil {
			j, _ := json.Marshal(GetFailed("Failed receiving the messages"))
			w.Write(j)
			return
		}
		hub := groups.GetOrCreateRoom(message.GroupName)
		SendSSEMessage(w, hub, message)
	})
	http.HandleFunc(fmt.Sprintf("%s/sse/", urlPrefix), func(w http.ResponseWriter, r *http.Request) {
		groupName, err := GetParam(r, "group")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Missing or invalid groupname parameter", http.StatusBadRequest)
			return
		}
		hub := groups.GetOrCreateRoom(groupName)
		registerSSEClient(w, r, hub)
	})
	http.Handle("/", &templateHandler{filename: "index.html", settings: settings})
	http.Handle("/video", &templateHandler{filename: "video.html", settings: settings})
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
