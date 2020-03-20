package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type server struct {
	mux sync.Mutex
}

func main() {
	s := &server{}

	mux := http.NewServeMux()
	mux.Handle("/", s)
	mux.Handle("/inc", s)
	mux.Handle("/dec", s)
	mux.Handle("/get", s)

	log.Fatal(http.ListenAndServe(":3030", mux))
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		if !s.serveHome(w) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/inc":
		if !s.increment() {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/dec":
		if !s.decrement() {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/get":
		val := s.get()
		w.Write([]byte(fmt.Sprintf(`{"val": %v}`, val)))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(`{"error": "Not valid API endpoint"}`))
	}
}

func (s *server) serveHome(w http.ResponseWriter) bool {
	bytes, err := ioutil.ReadFile("index.html")
	if err != nil {
		return false
	}
	w.Write(bytes)
	return true
}

func (s *server) increment() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return false
}

func (s *server) decrement() bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return false
}

func (s *server) get() int {
	return 0
}
