package html_test

import (
	"testing"

	. "github.com/tinywasm/html"
)

func TestHTML_New(t *testing.T) {
	h := New(Div("content").Class("wrap"))
	if h.String() == "" {
		t.Fatal("expected non-empty")
	}
}

func TestHTML_Raw(t *testing.T) {
	raw := `<div class="raw">test</div>`
	h := Raw(raw)
	if h.String() != raw {
		t.Fatalf("got %q", h.String())
	}
}

func TestHTML_Nil_Safe(t *testing.T) {
	var h *HTML
	if h.String() != "" {
		t.Fatal("nil HTML.String() should return empty")
	}
}
