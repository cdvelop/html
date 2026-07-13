package html

import (
	. "github.com/tinywasm/dom"
)

// Block containers
func Div() *Element    { return NewElement("div") }
func Span() *Element   { return NewElement("span") }
func P() *Element      { return NewElement("p") }
func Pre() *Element    { return NewElement("pre") }
func Code() *Element   { return NewElement("code") }
func Strong() *Element { return NewElement("strong") }
func Small() *Element  { return NewElement("small") }
func Mark() *Element   { return NewElement("mark") }

// Headings
func H1() *Element { return NewElement("h1") }
func H2() *Element { return NewElement("h2") }
func H3() *Element { return NewElement("h3") }
func H4() *Element { return NewElement("h4") }
func H5() *Element { return NewElement("h5") }
func H6() *Element { return NewElement("h6") }

// Lists
func Ul() *Element { return NewElement("ul") }
func Ol() *Element { return NewElement("ol") }
func Li() *Element { return NewElement("li") }

// Semantic layout
func Nav() *Element        { return NewElement("nav") }
func Section() *Element    { return NewElement("section") }
func Main() *Element       { return NewElement("main") }
func Article() *Element    { return NewElement("article") }
func Header() *Element     { return NewElement("header") }
func Footer() *Element     { return NewElement("footer") }
func Aside() *Element      { return NewElement("aside") }
func Details() *Element    { return NewElement("details") }
func Summary() *Element    { return NewElement("summary") }
func Dialog() *Element     { return NewElement("dialog") }
func Figure() *Element     { return NewElement("figure") }
func Figcaption() *Element { return NewElement("figcaption") }

// Tables
func Table() *Element { return NewElement("table") }
func Thead() *Element { return NewElement("thead") }
func Tbody() *Element { return NewElement("tbody") }
func Tfoot() *Element { return NewElement("tfoot") }
func Tr() *Element    { return NewElement("tr") }
func Th() *Element    { return NewElement("th") }
func Td() *Element    { return NewElement("td") }

// Form-adjacent (non-form elements)
func Fieldset() *Element { return NewElement("fieldset") }
func Legend() *Element   { return NewElement("legend") }
func Label() *Element    { return NewElement("label") }
func Button() *Element   { return NewElement("button") }
func Canvas() *Element   { return NewElement("canvas") }
func Style() *Element    { return NewElement("style") }
func Script() *Element   { return NewElement("script") }

// Special
func A(href string) *Element {
	return NewElement("a").Attr("href", href)
}
func Input(typ string) *Element {
	return NewElement("input").NoCloseTag().Attr("type", typ)
}
func Option(value, text string) *Element {
	return NewElement("option").Attr("value", value).Text(text)
}
func SelectedOption(value, text string) *Element {
	return NewElement("option").Attr("value", value).Attr("selected", "").Text(text)
}
func Br() *Element { return NewElement("br").NoCloseTag() }
func Hr() *Element { return NewElement("hr").NoCloseTag() }

