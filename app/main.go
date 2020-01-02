package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/buoyantio/emojivoto/emojivoto-emoji-svc/emoji"
)

var Version string

type Emoji struct {
	Name      string `json:"name"`
	Character string `json:"character"`
}

func logger(server string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", server, r.Method, r.URL.Path)
	})
}

func mainWebHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	emoji := Emoji{
		Name:      ":doughnut:",
		Character: emoji.NewAllEmoji().WithShortcode(":doughnut:").Unicode,
	}

	s, _ := json.Marshal(emoji)
	w.Write([]byte(s))
}

func mainApiHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	s := `{"message": "Not Found"}`
	w.Write([]byte(s))
}

func main() {
	webserver := http.NewServeMux()
	api := http.NewServeMux()

	webserver.HandleFunc("/", mainWebHandler)
	api.HandleFunc("/", mainApiHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8081", logger("api", api)))
	}()

	fmt.Printf("%s", Version)
	log.Fatal(http.ListenAndServe(":8080", logger("web", webserver)))
}
