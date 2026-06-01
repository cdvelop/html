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

// Document builds the complete document shell as *Element (the <html> element).
// Does not include <!DOCTYPE html> — use DocumentString for the full output.
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

// DocumentString renders the full HTML document including <!DOCTYPE html>.
// This is what assetmin writes to disk as index.html.
func DocumentString(opts DocumentOptions, body ...any) string {
	return "<!DOCTYPE html>" + Document(opts, body...).String()
}
