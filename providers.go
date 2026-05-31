package html

// HTMLProvider is an optional capability: components that provide
// a custom SSR HTML template fragment for injection by assetmin.
//
// Implement this in a component's html.go file (//go:build !wasm).
// Only needed when the component has a custom SSR HTML template distinct
// from its Render() output. Most components do NOT need this.
//
// Example in mycomponent/html.go:
//
//	//go:build !wasm
//	package mycomponent
//	import "github.com/tinywasm/html"
//
//	func (c *MyComponent) RenderHTML() *html.HTML {
//	    return html.New(Div(clsRoot.AsAttr(), /* ... */))
//	}
type HTMLProvider interface {
	RenderHTML() *HTML
}
