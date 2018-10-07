// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package htmltables

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"strconv"
	"strings"
)

// Table holds a simple table of headers and rows.
type Table struct {
	Headers []string
	Rows    [][]string
}

// Parse parses a html fragment or whole document looking for HTML
// tables.
func Parse(s string) ([]*Table, error) {
	node, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	tables := []*Table{}
	parse(node, &tables)
	for kk, t := range tables {
		tables[kk] = addMissingColumns(t)
	}
	return tables, nil
}

func innerHTML(n *html.Node) string {
	clone := html.Node{
		FirstChild: n.FirstChild,
		LastChild:  n.LastChild,
		Type:       html.ElementNode,
		DataAtom:   atom.Body,
		Data:       "body",
	}
	writer := &strings.Builder{}
	if err := html.Render(writer, &clone); err != nil {
		return err.Error()
	}
	s := writer.String()
	s = s[strings.Index(s, "<body>")+6:]
	s = s[:strings.Index(s, "</body>")]
	return s
}

func parse(n *html.Node, tables *[]*Table) {
	switch n.DataAtom {
	case atom.Table:
		if len(*tables) == 0 {
			*tables = append(*tables, &Table{})
		}
	case atom.Th:
		t := (*tables)[len(*tables)-1]
		t.Headers = append(t.Headers, innerHTML(n))
	case atom.Tr:
		t := (*tables)[len(*tables)-1]
		t.Rows = append(t.Rows, []string{})
	case atom.Td:
		t := (*tables)[len(*tables)-1]
		l := len(t.Rows) - 1
		t.Rows[l] = append(t.Rows[l], innerHTML(n))
		return
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		parse(child, tables)
	}
}

func addMissingColumns(t *Table) *Table {
	cols := len(t.Headers)
	for _, row := range t.Rows {
		if len(row) > cols {
			cols = len(row)
		}
	}
	for len(t.Headers) < cols {
		name := "Col " + strconv.Itoa(len(t.Headers)+1)
		t.Headers = append(t.Headers, name)
	}
	for kk, row := range t.Rows {
		for len(row) < cols {
			row = append(row, "")
			t.Rows[kk] = row
		}
	}
	return t
}
