package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/buoyantio/emojivoto/emojivoto-emoji-svc/emoji"
	_ "github.com/buoyantio/emojivoto/emojivoto-emoji-svc/emoji"
)

var (
	Version string
	emojis  = Emojis{
		Data: make(map[string]string),
	}
)

type Emojis struct {
	sync.Mutex
	Data map[string]string
}

func logger(server string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s", server, r.Method, r.URL.Path)
	})
}

func methodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Set("Content-Type", "application/json")
	s := `{"message": "Method Not Allowed"}`
	w.Write([]byte(s))
}

func mainWebHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	emojis.Lock()
	s, _ := json.Marshal(emojis.Data)
	emojis.Unlock()
	w.Write([]byte(s))
}

func mainApiHandler(w http.ResponseWriter, r *http.Request) {
	methodNotAllowed(w)
}

func emojisApiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		emojis.Lock()
		s, _ := json.Marshal(emojis.Data)
		emojis.Unlock()
		w.Write([]byte(s))
		return
	}
	methodNotAllowed(w)
}

func emojiApiHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/emojis/")
	if r.Method == "GET" {
		emojis.Lock()
		if emojis.Data[id] != "" {
			w.WriteHeader(http.StatusOK)
			s, _ := json.Marshal(emojis.Data[id])
			w.Write([]byte(s))
			return
		}
		emojis.Unlock()
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		return
	}
	if r.Method == "PUT" {
		e := emoji.NewAllEmoji().WithShortcode(fmt.Sprintf(":%s:", id))
		if e != nil {
			w.WriteHeader(http.StatusOK)
			emojis.Lock()
			emojis.Data[id] = e.Unicode
			s, _ := json.Marshal(emojis.Data[id])
			emojis.Unlock()
			w.Write([]byte(s))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
		return
	}
	if r.Method == "DELETE" {
		emojis.Lock()
		if emojis.Data[id] != "" {
			w.WriteHeader(http.StatusNoContent)
			delete(emojis.Data, id)
			w.Write([]byte(`{}`))
			return
		}
		emojis.Unlock()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
		return
	}
	methodNotAllowed(w)
}

func main() {
	webserver := http.NewServeMux()
	api := http.NewServeMux()

	webserver.HandleFunc("/", mainWebHandler)
	api.HandleFunc("/", mainApiHandler)
	api.HandleFunc("/emojis", emojisApiHandler)
	api.HandleFunc("/emojis/", emojiApiHandler)

	go func() {
		log.Fatal(http.ListenAndServe(":8081", logger("api", api)))
	}()

	fmt.Printf("Starting webserver version %s\n", Version)
	log.Fatal(http.ListenAndServe(":8080", logger("web", webserver)))
}
