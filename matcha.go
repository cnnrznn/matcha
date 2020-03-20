package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
)

type server struct {
	mux sync.Mutex
	pool *sql.DB
}

func newServer() *server {
	s := &server{}
	pool, err := sql.Open("postgres",
		"user=cnnrznn password=matchatime dbname=matcha")
	if err != nil {
		log.Fatal("Could not open database connection")
	}
	s.pool = pool

	return s
}

func main() {
	s := newServer()

	mux := http.NewServeMux()
	mux.Handle("/", s)
	mux.Handle("/inc", s)
	mux.Handle("/get", s)

	log.Fatal(http.ListenAndServe(":3030", mux))
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	switch req.URL.Path {
	case "/":
		log.Println("Hit: /")
		if !s.serveHome(w) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/inc":
		log.Println("Hit: /inc")
		if !s.increment(req) {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "/get":
		log.Println("Hit: /get")
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

func (s *server) increment(req *http.Request) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	src := req.URL.Query()["src"]
	query := fmt.Sprintf("INSERT INTO log VALUES (now(), '%v');", src)
	_, err := s.pool.Query(query)

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (s *server) get() int {
	result := 0
	rows, _ := s.pool.Query("SELECT COUNT(*) FROM log")
	defer rows.Close()

	rows.Next()
	rows.Scan(&result)

	return result
}
