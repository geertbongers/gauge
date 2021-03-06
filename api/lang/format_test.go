// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package lang

import (
	"testing"

	"encoding/json"
	"fmt"
	"reflect"

	"github.com/sourcegraph/go-langserver/pkg/lsp"
	"github.com/sourcegraph/jsonrpc2"
)

func TestFormat(t *testing.T) {
	specText := `# Specification Heading

## Scenario Heading

* Step text`

	openFilesCache = &files{cache: make(map[lsp.DocumentURI][]string)}
	openFilesCache.add("foo.spec", specText)

	want := []lsp.TextEdit{
		{
			Range: lsp.Range{
				Start: lsp.Position{Line: 0, Character: 0},
				End:   lsp.Position{Line: 5, Character: 57},
			},
			NewText: specText + "\n",
		},
	}

	b, _ := json.Marshal(lsp.DocumentFormattingParams{TextDocument: lsp.TextDocumentIdentifier{URI: "foo.spec"}, Options: lsp.FormattingOptions{}})
	p := json.RawMessage(b)

	got, err := format(&jsonrpc2.Request{Params: &p})
	if err != nil {
		t.Fatalf("Expected error == nil in format, got %s", err.Error())
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("format failed, want: `%v`, got: `%v`", want, got)
	}
}

func TestFormatParseError(t *testing.T) {
	specText := `Specification Heading
=====================

# Scenario Heading


* Step text`

	openFilesCache = &files{cache: make(map[lsp.DocumentURI][]string)}
	openFilesCache.add("foo.spec", specText)

	specFile := lsp.DocumentURI("foo.spec")

	b, _ := json.Marshal(lsp.DocumentFormattingParams{TextDocument: lsp.TextDocumentIdentifier{URI: specFile}, Options: lsp.FormattingOptions{}})
	p := json.RawMessage(b)

	expectedError := fmt.Errorf("failed to format document. Fix all the problems first")

	data, err := format(&jsonrpc2.Request{Params: &p})
	if data != nil {
		t.Fatalf("Expected data == nil in format, got %s", data)
	}
	if err.Error() != expectedError.Error() {
		t.Fatalf(" want : %s\ngot : %s", expectedError.Error(), err.Error())
	}

}
