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
	Attributes map[string]string
	Headers    []string
	Rows       [][]string
}

// Parse parses a html fragment or whole document looking for HTML
// tables. It converts all cells into text, stripping away any HTML content.
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

func innerText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	result := ""
	for x := n.FirstChild; x != nil; x = x.NextSibling {
		result += innerText(x)
	}
	return result
}

func parse(n *html.Node, tables *[]*Table) {
	strip := strings.TrimSpace
	switch n.DataAtom {
	case atom.Table:
		t := &Table{}
		for _, at := range n.Attr {
			if t.Attributes == nil {
				t.Attributes = map[string]string{}
			}
			t.Attributes[at.Key] = at.Val
		}
		*tables = append(*tables, t)
	case atom.Th:
		t := (*tables)[len(*tables)-1]
		t.Headers = append(t.Headers, strip(innerText(n)))
	case atom.Tr:
		t := (*tables)[len(*tables)-1]
		t.Rows = append(t.Rows, []string{})
	case atom.Td:
		t := (*tables)[len(*tables)-1]
		l := len(t.Rows) - 1
		t.Rows[l] = append(t.Rows[l], strip(innerText(n)))
		return
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		parse(child, tables)
	}
}

func addMissingColumns(t *Table) *Table {
	cols := len(t.Headers)
	rows := make([][]string, 0, len(t.Rows))
	for _, row := range t.Rows {
		if len(row) > 0 {
			rows = append(rows, row)
		}
		if len(row) > cols {
			cols = len(row)
		}
	}
	for len(t.Headers) < cols {
		name := "Col " + strconv.Itoa(len(t.Headers)+1)
		t.Headers = append(t.Headers, name)
	}
	for kk := range rows {
		for len(rows[kk]) < cols {
			rows[kk] = append(rows[kk], "")
		}
	}
	t.Rows = rows
	return t
}
