// Copyright (C) 2018 Ramesh Vyaghrapuri. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package htmltables_test

import (
	"github.com/tvastar/htmltables"
	"reflect"
	"testing"
)

type Table = htmltables.Table

func TestParse(t *testing.T) {
	cases := map[string][]*Table{
		"":         {},
		"goop":     {},
		"<table/>": {{Headers: []string(nil), Rows: [][]string{}}},
		"<table><td>Hello</td>": {{
			Headers: []string{"Col 1"},
			Rows:    [][]string{{"Hello"}},
		}},
		"<table><th>boo</th><td>Hello</td>": {{
			Headers: []string{"boo"},
			Rows:    [][]string{{"Hello"}},
		}},
		"<table><tr><td>1</td><td>2</td></tr><td>3": {{
			Headers: []string{"Col 1", "Col 2"},
			Rows:    [][]string{{"1", "2"}, {"3", ""}},
		}},
		"<table><thead><th>boo</th></thead><tr/><td>Hello</td>": {{
			Headers: []string{"boo"},
			Rows:    [][]string{{"Hello"}},
		}},
		"<table><thead><td><a href=\"x\">Hello</a></td>": {{
			Headers: []string{"Col 1"},
			Rows:    [][]string{{"Hello"}},
		}},
		"<table><thead/><td>Hello</td></th></table><table><thead/><td>Hello2</td></th>": {
			{
				Headers: []string{"Col 1"},
				Rows:    [][]string{{"Hello"}},
			},
			{
				Headers: []string{"Col 1"},
				Rows:    [][]string{{"Hello2"}},
			},
		},
	}

	for caseName, expected := range cases {
		t.Run(caseName, func(t *testing.T) {
			tx, err := htmltables.Parse(caseName)
			if err != nil {
				t.Fatal("Unexpected parse failure", err)
			}
			if !reflect.DeepEqual(expected, tx) {
				t.Errorf("Unexpected result %#v", tx[0])
			}
		})
	}
}
