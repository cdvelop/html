package html

import (
	. "github.com/tinywasm/dom"
)

// Block containers
func Div(children ...any) *Element        { return NewElement("div").Add(children...) }
func Span(children ...any) *Element       { return NewElement("span").Add(children...) }
func P(children ...any) *Element          { return NewElement("p").Add(children...) }
func Pre(children ...any) *Element        { return NewElement("pre").Add(children...) }
func Code(children ...any) *Element       { return NewElement("code").Add(children...) }
func Strong(children ...any) *Element     { return NewElement("strong").Add(children...) }
func Small(children ...any) *Element      { return NewElement("small").Add(children...) }
func Mark(children ...any) *Element       { return NewElement("mark").Add(children...) }

// Headings
func H1(children ...any) *Element { return NewElement("h1").Add(children...) }
func H2(children ...any) *Element { return NewElement("h2").Add(children...) }
func H3(children ...any) *Element { return NewElement("h3").Add(children...) }
func H4(children ...any) *Element { return NewElement("h4").Add(children...) }
func H5(children ...any) *Element { return NewElement("h5").Add(children...) }
func H6(children ...any) *Element { return NewElement("h6").Add(children...) }

// Lists
func Ul(children ...any) *Element { return NewElement("ul").Add(children...) }
func Ol(children ...any) *Element { return NewElement("ol").Add(children...) }
func Li(children ...any) *Element { return NewElement("li").Add(children...) }

// Semantic layout
func Nav(children ...any) *Element        { return NewElement("nav").Add(children...) }
func Section(children ...any) *Element    { return NewElement("section").Add(children...) }
func Main(children ...any) *Element       { return NewElement("main").Add(children...) }
func Article(children ...any) *Element    { return NewElement("article").Add(children...) }
func Header(children ...any) *Element     { return NewElement("header").Add(children...) }
func Footer(children ...any) *Element     { return NewElement("footer").Add(children...) }
func Aside(children ...any) *Element      { return NewElement("aside").Add(children...) }
func Details(children ...any) *Element    { return NewElement("details").Add(children...) }
func Summary(children ...any) *Element    { return NewElement("summary").Add(children...) }
func Dialog(children ...any) *Element     { return NewElement("dialog").Add(children...) }
func Figure(children ...any) *Element     { return NewElement("figure").Add(children...) }
func Figcaption(children ...any) *Element { return NewElement("figcaption").Add(children...) }

// Tables
func Table(children ...any) *Element { return NewElement("table").Add(children...) }
func Thead(children ...any) *Element { return NewElement("thead").Add(children...) }
func Tbody(children ...any) *Element { return NewElement("tbody").Add(children...) }
func Tfoot(children ...any) *Element { return NewElement("tfoot").Add(children...) }
func Tr(children ...any) *Element    { return NewElement("tr").Add(children...) }
func Th(children ...any) *Element    { return NewElement("th").Add(children...) }
func Td(children ...any) *Element    { return NewElement("td").Add(children...) }

// Form-adjacent (non-form elements)
func Fieldset(children ...any) *Element { return NewElement("fieldset").Add(children...) }
func Legend(children ...any) *Element   { return NewElement("legend").Add(children...) }
func Label(children ...any) *Element    { return NewElement("label").Add(children...) }
func Button(children ...any) *Element   { return NewElement("button").Add(children...) }
func Canvas(children ...any) *Element   { return NewElement("canvas").Add(children...) }
func Style(children ...any) *Element    { return NewElement("style").Add(children...) }
func Script(children ...any) *Element   { return NewElement("script").Add(children...) }

// Special
func A(href string, children ...any) *Element {
	return NewElement("a").Attr("href", href).Add(children...)
}
func Input(typ string) *Element {
	return NewElement("input").NoCloseTag().Attr("type", typ)
}
func Option(value, text string) *Element {
	return NewElement("option").Attr("value", value).Add(text)
}
func SelectedOption(value, text string) *Element {
	return NewElement("option").Attr("value", value).Attr("selected", "").Add(text)
}
func Br() *Element { return NewElement("br").NoCloseTag() }
func Hr() *Element { return NewElement("hr").NoCloseTag() }

// SVG placeholders to satisfy tests until tinywasm/svg is used
func Svg(children ...any) *Element { return NewElement("svg").Add(children...) }
func Use(children ...any) *Element { return NewElement("use").NoCloseTag().Add(children...) }
