package html

// HTMLProvider is an optional capability: components that expose
// a static SSR HTML template fragment for injection by assetmin.
//
// Implement in a component's html.go file (//go:build !wasm).
// Most components do NOT need this — only those with a static shell
// distinct from their dynamic Render() output.
//
// Example in mycomponent/html.go:
//
//	//go:build !wasm
//	package mycomponent
//	import . "github.com/tinywasm/html"
//
//	func (c *MyComponent) RenderHTML() string {
//	    return Div(clsRoot.AsAttr()).String()
//	}
//
// For raw static HTML:
//
//	func (c *MyComponent) RenderHTML() string {
//	    return `<div class="root"></div>`
//	}
type HTMLProvider interface {
	RenderHTML() string
}

// Page is one emitted document of a multi-page static site: where it is
// served from, the head metadata that makes it rank on its own, and its body
// markup.
//
// A site compiler turns each Page into one file. RenderHTML above is the
// single-page case (one document at "/"); Page is what a module declares when
// it owns several routes — a specialty per URL, an article per URL — each
// needing its own title, description and canonical, which is precisely what a
// single shared index.html cannot express.
type Page struct {
	// Path is the site-absolute route, e.g. "/" or "/especialidades/oftalmologia/".
	// A trailing slash means the compiler writes <path>/index.html so the URL
	// serves without an extension.
	Path string

	// Doc is the head metadata for THIS page. Title, Description and Canonical
	// are what differ per page; CSSURL/JSURL/FaviconURL are normally filled in
	// by the compiler, which is the only layer that knows the built asset names.
	Doc DocumentOptions

	// Body is the rendered markup placed inside <body>.
	Body string
}

// PagesProvider is an optional capability: a module that owns several routes
// declares them here instead of (or in addition to) RenderHTML.
//
// Implement in a //go:build !wasm file, like the other SSR producers:
//
//	func (m *Module) RenderPages() []Page {
//	    return []Page{{
//	        Path: "/especialidades/oftalmologia/",
//	        Doc:  DocumentOptions{Title: "Oftalmología en Chillán — …", Description: "…"},
//	        Body: body.String(),
//	    }}
//	}
type PagesProvider interface {
	RenderPages() []Page
}
