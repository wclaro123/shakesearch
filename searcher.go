package main

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"sort"
	"strings"
)

type Searcher struct {
	CompleteWorks      string
	CompleteWorksLower string
	Lines              Lines
	SuffixArray        *suffixarray.Index
	SuffixArrayLower   *suffixarray.Index
}

type SearchResult struct {
	LinesText []string `json:"lines_text"`
	FromLine  int      `json:"from_line"`
	ToLine    int      `json:"to_line"`
}

type Response struct {
	TotalQuantity int            `json:"total_quantity"`
	PageQuantity  int            `json:"page_quantity"`
	SearchResults []SearchResult `json:"search_results"`
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}
	datLower := bytes.ToLower(dat)
	s.CompleteWorks = string(dat)
	s.CompleteWorksLower = string(datLower)
	s.SuffixArray = suffixarray.New(dat)
	s.SuffixArrayLower = suffixarray.New(datLower)
	lines := strings.Split(string(dat), "\r\n")
	var from, to int
	for i, line := range lines {
		from = to
		l := len(line)
		to += l
		s.Lines = append(s.Lines, Line{
			Text:  line,
			Index: i,
			From:  from,
			To:    to,
		})
		to += 2
	}

	return nil
}

func (s *Searcher) Search(query string, limit, page int) (result Response) {
	result = s.FullMatch(query, limit, page)
	if result.TotalQuantity == 0 {
		return s.MultiWords(query, limit, page)
	}
	return
}

func (s Searcher) FullMatch(query string, limit, page int) (result Response) {
	idxs := s.Lookup(query)

	if len(idxs) == 0 {
		return
	}

	initialPos := (page - 1) * limit
	if len(idxs) < initialPos {
		return
	}
	finalPos := page * limit
	if finalPos > len(idxs) {
		finalPos = len(idxs)
	}

	result.TotalQuantity = len(idxs)
	result.PageQuantity = finalPos - initialPos
	sort.Ints(idxs)
	for i := initialPos; i < finalPos; i++ {
		line := s.Lines.FindLine(idxs[i])
		q := Query{Query: query, Index: idxs[i]}
		line.Query = q
		line.Queries = Queries{q}
		s.Lines[line.Index] = line
		r := s.Lines.BuildSearchResult(line.Index)
		result.SearchResults = append(result.SearchResults, r)
	}
	return
}

func (s Searcher) MultiWords(query string, limit, page int) (response Response) {
	words := strings.Fields(query)
	if len(words) < 2 {
		return
	}
	var (
		arrLines []Lines
		results  []SearchResult
	)

	for _, word := range words {
		idxs := s.Lookup(word)
		sort.Ints(idxs)
		arrLines = append(arrLines, s.GetLines(idxs, word))
	}
	if arrLines == nil {
		return
	}
	lines := s.IntersectLines(arrLines)

	initialPos := (page - 1) * limit
	if len(lines) < initialPos {
		return
	}

	finalPos := page * limit
	if finalPos > len(lines) {
		finalPos = len(lines)
	}

	response.TotalQuantity = len(lines)
	response.PageQuantity = finalPos - initialPos

	lines = lines[initialPos:finalPos]

	for _, line := range lines {
		r := s.Lines.BuildSearchResult(line.Index)
		results = append(results, r)
	}
	response.SearchResults = results

	return
}

func (s Searcher) Lookup(query string) []int {
	return s.SuffixArrayLower.Lookup([]byte(strings.ToLower(query)), -1)
}

func (s Searcher) GetLines(idxs []int, query string) (lines Lines) {
	for _, idx := range idxs {
		line := s.Lines.FindLine(idx)
		q := Query{Query: query, Index: idx}
		line.Query = q
		line.Queries = Queries{q}
		s.Lines[line.Index] = line
		lines = append(lines, line)
	}
	return
}

func (s Searcher) IntersectLines(arrLines []Lines) (ln Lines) {
	lineMap := make(map[int]Line)
	for _, line := range arrLines[0] {
		lineMap[line.Index] = line
	}

	for i, lines := range arrLines[1:] {
		tempMap := make(map[int]Line)
		for _, line := range lines {

			if lnMap, found := lineMap[line.Index]; found {
				lnMap.Queries = append(lnMap.Queries, line.Query)
				tempMap[line.Index] = lnMap
				s.Lines[lnMap.Index] = lnMap
				if i == len(arrLines)-2 {
					ln = append(ln, line)
				}
			}
		}
		lineMap = tempMap
	}

	return
}
