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
		"<table/>": {{}},
		"<table><td>Hello</td>": {{
			Headers: []string{"Col 1"},
			Rows:    [][]string{{"Hello"}},
		}},
		"<table><th>boo</th><td>Hello</td>": {{
			Headers: []string{"boo"},
			Rows:    [][]string{{"Hello"}},
		}},
		"<table><thead><th>boo</th></thead><tr/><td>Hello</td>": {{
			Headers: []string{"boo"},
			Rows:    [][]string{{""}, {"Hello"}},
		}},
		"<table><thead><td><a href=\"x\">Hello</a></td>": {{
			Headers: []string{"Col 1"},
			Rows:    [][]string{{"<a href=\"x\">Hello</a>"}},
		}},
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
