package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
	"math"
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

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		page, ok := r.URL.Query()["p"]

		if !ok {
			page = []string{"1"}
		}
		results := searcher.Search(query[0], page[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)

	// Convert data to lower case.
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(string(dat))))

	return nil
}

func (s *Searcher) Search(query string, page string) []string {

	pageNo, err := strconv.Atoi(page)

	if err != nil {
		return []string{}
	}
	// Search with query in lowercase.
	idxs := s.SuffixArray.Lookup([]byte(strings.ToLower(query)), 20 * pageNo)
	results := []string{}
	for _, idx := range idxs[20 * (pageNo - 1):int(math.Min(float64(20 * pageNo), float64(len(idxs))))] {
		// Show 2 sentences with that query, not random words.
		var offsetFirst = strings.LastIndex(s.CompleteWorks[:idx], ".")
		var offsetMid = strings.Index(s.CompleteWorks[idx:], ".")
		var offsetLast = strings.Index(s.CompleteWorks[idx + offsetMid + 1:], ".")

		results = append(results, s.CompleteWorks[offsetFirst + 2:idx + offsetMid + offsetLast + 2])
	}
	return results
}
