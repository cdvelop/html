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
