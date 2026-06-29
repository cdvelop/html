package html_test

import (
	"strings"
	"testing"

	. "github.com/tinywasm/html"
)

func TestDocument(t *testing.T) {
	opts := DocumentOptions{
		Lang:       "es",
		Title:      "Mi App",
		CSSURL:     "/styles.css",
		JSURL:      "/bundle.js",
		FaviconURL: "/favicon.ico",
		Head:       "<meta name=\"theme-color\" content=\"#000\">",
	}

	doc := Document(opts, Div().Text("Hola Mundo"))

	got := doc.String()

	checks := []string{
		"lang='es'",
		"<title>Mi App</title>",
		"rel='stylesheet'",
		"href='/styles.css'",
		"src='/bundle.js'",
		"defer",
		"rel='icon'",
		"href='/favicon.ico'",
		"<meta charset='utf-8'",
		"theme-color",
		"<div>Hola Mundo</div>",
	}

	for _, check := range checks {
		if !strings.Contains(got, check) {
			t.Errorf("expected to find %q in output, but didn't.\nGot: %s", check, got)
		}
	}
}

func TestDocumentString(t *testing.T) {
	got := DocumentString(DocumentOptions{Title: "App"})
	if !strings.HasPrefix(got, "<!DOCTYPE html>") {
		t.Errorf("expected <!DOCTYPE html> prefix, got: %s", got[:30])
	}
	if !strings.Contains(got, "<html") {
		t.Error("expected <html element")
	}
}
