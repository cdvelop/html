package html

// HTML represents a server-side HTML fragment for SSR injection by assetmin.
// Follows the same pattern as *css.Stylesheet and []*js.Script.
// Build with New() or Raw().
type HTML struct {
	content string
}

// New builds an HTML fragment by serializing an Element tree.
// Use this in component html.go files: RenderHTML() *HTML { return html.New(Div(...)) }
func New(el interface{ String() string }) *HTML {
	return &HTML{content: el.String()}
}

// Raw wraps a pre-built HTML string.
// Use when you have a static template or embedded HTML.
func Raw(content string) *HTML {
	return &HTML{content: content}
}

// String returns the rendered HTML string.
// Called by assetmin extractor via .String() to get the content.
func (h *HTML) String() string {
	if h == nil {
		return ""
	}
	return h.content
}
