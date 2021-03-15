package main

import (
	"fmt"
	"sort"
	"strings"
)

type Line struct {
	Text    string
	Index   int
	From    int
	To      int
	Query   Query
	Queries Queries
}

type Lines []Line

type Query struct {
	Query string
	Index int
}

type Queries []Query

func (l Lines) FindLine(index int) (ln Line) {
	for _, line := range l {
		if index <= line.To && index >= line.From {
			return line
		}
	}
	return ln
}

func (l Lines) BuildSearchResult(lineIndex int) (result SearchResult) {
	const numLines = 2
	var from, to int
	var lines []string

	if lineIndex > numLines {
		from = lineIndex - numLines
	}
	to = lineIndex + numLines
	if to >= len(l) {
		to = len(l) - 1
	}
	for i := from; i <= to; i++ {
		fullText := l[i].Text
		if i == lineIndex {
			fullText = highlightLine(l[i])
		}
		fullText = strings.ReplaceAll(fullText, " ", "&nbsp;")
		lines = append(lines, fullText)
	}

	return SearchResult{
		LinesText: lines,
		FromLine:  from + 1,
		ToLine:    to + 1,
	}
}

func (q Queries) Len() int           { return len(q) }
func (q Queries) Less(i, j int) bool { return q[i].Index > q[j].Index }
func (q Queries) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func highlightLine(line Line) string {
	queries := line.Queries
	sort.Sort(queries)
	for _, query := range queries {
		queryLineIndex := query.Index - line.From
		textLeft := line.Text[:queryLineIndex]
		textQuery := fmt.Sprintf("%s%s%s", "<b>", line.Text[queryLineIndex:queryLineIndex+len(query.Query)], "</b>")
		textRight := line.Text[queryLineIndex+len(query.Query):]
		line.Text = fmt.Sprintf("%s%s%s", textLeft, textQuery, textRight)
	}
	line.Text = strings.ReplaceAll(line.Text, " ", "&nbsp;")
	return line.Text
}
