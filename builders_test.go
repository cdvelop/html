package html_test

import (
	"strings"
	"testing"

	. "github.com/tinywasm/html"
)

func TestDiv_Basic(t *testing.T) {
	el := Div(H2("Hello"), P("World")).Class("root")
	got := el.String()
	if !strings.Contains(got, `class='root'`) {
		t.Errorf("expected class root, got %q", got)
	}
	if !strings.Contains(got, "<h2>Hello</h2>") {
		t.Errorf("expected h2, got %q", got)
	}
	if !strings.Contains(got, "<p>World</p>") {
		t.Errorf("expected p, got %q", got)
	}
}

func TestA_HasHref(t *testing.T) {
	el := A("/home", "Home")
	got := el.String()
	if !strings.Contains(got, `href='/home'`) {
		t.Errorf("expected href, got %q", got)
	}
	if !strings.Contains(got, "Home") {
		t.Errorf("expected text, got %q", got)
	}
}

func TestInput_IsVoid(t *testing.T) {
	el := Input("text")
	got := el.String()
	// void element: no closing tag
	if strings.Contains(got, "</input>") {
		t.Error("input should be void")
	}
}

func TestBr_IsVoid(t *testing.T) {
	got := Br().String()
	if strings.Contains(got, "</br>") {
		t.Error("br should be void")
	}
}
