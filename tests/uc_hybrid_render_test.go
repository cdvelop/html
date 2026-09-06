//go:build wasm

package html_test

import (
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

// DynamicComp uses ViewRenderer (DSL)
type DynamicComp struct {
	dom.Element
}

func (c *DynamicComp) Render() *dom.Element {
	return Div().ID("dynamic").Text("Dynamic Content")
}

// StaticComp uses String()
type StaticComp struct {
	dom.Element
}

func (c *StaticComp) String() string {
	return `<div id="static">Static Content</div>`
}

func TestHybridRendering(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Render Dynamic Component (DSL)", func(t *testing.T) {
		c := &DynamicComp{}
		dom.Render("root", c)

		_, ok := GetRef("dynamic")
		if !ok {
			t.Fatal("Dynamic component not rendered")
		}
	})

	t.Run("Render Static Component (String)", func(t *testing.T) {
		c := &StaticComp{}
		dom.Render("root", c)

		_, ok := GetRef("static")
		if !ok {
			t.Fatal("Static component not rendered")
		}
	})
}
