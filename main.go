package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := ""
		page := 1
		limit := 20
		q, ok := r.URL.Query()["q"]
		if !ok || len(q[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		query = q[0]

		var err error
		p, ok := r.URL.Query()["page"]
		if ok && len(p[0]) > 0 {
			page, err = strconv.Atoi(p[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid page in URL params"))
				return
			}
		}
		l, ok := r.URL.Query()["limit"]
		if ok && len(p[0]) > 0 {
			limit, err = strconv.Atoi(l[0])
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid limit in URL params"))
				return
			}
		}

		results := searcher.Search(query, limit, page)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err = enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}
