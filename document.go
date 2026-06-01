package html

import (
	. "github.com/tinywasm/dom"
)

// DocumentOptions configures the document shell.
type DocumentOptions struct {
	Lang       string // default "en"
	Title      string
	CSSURL     string // <link rel="stylesheet">
	JSURL      string // <script src defer>
	FaviconURL string // <link rel="icon">
	Head       string // extra markup for <head> (optional)
}

// Document builds the complete document shell as *Element.
// It returns the <html> element. Note that since dom.Element
// represents a single tag, the <!DOCTYPE html> declaration
// is not included in the returned *Element itself.
//
// To render the full document including doctype, use:
//
//	"<!DOCTYPE html>" + doc.String()
func Document(opts DocumentOptions, body ...any) *Element {
	if opts.Lang == "" {
		opts.Lang = "en"
	}

	head := NewElement("head")
	head.Add(NewElement("meta").NoCloseTag().Attr("charset", "utf-8"))

	if opts.Title != "" {
		head.Add(NewElement("title").Add(opts.Title))
	}
	if opts.FaviconURL != "" {
		head.Add(NewElement("link").NoCloseTag().Attr("rel", "icon").Attr("href", opts.FaviconURL))
	}
	if opts.CSSURL != "" {
		head.Add(NewElement("link").NoCloseTag().Attr("rel", "stylesheet").Attr("href", opts.CSSURL))
	}
	if opts.JSURL != "" {
		head.Add(NewElement("script").Attr("src", opts.JSURL).Attr("defer", ""))
	}
	if opts.Head != "" {
		head.Add(opts.Head)
	}

	bodyEl := NewElement("body").Add(body...)

	return NewElement("html").Attr("lang", opts.Lang).Add(head, bodyEl)
}
