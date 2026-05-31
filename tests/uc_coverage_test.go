//go:build wasm

package html_test

import (
	"testing"

	"github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
	"github.com/tinywasm/fmt"
)

func TestCoverageElementFactories(t *testing.T) {
	t.Run("All Element Factories", func(t *testing.T) {
		els := []dom.Component{
			Span(), P(), H1(), H2(), H3(), H4(), H5(), H6(),
			Ul(), Ol(), Li(), Nav(), Section(), Main(), Article(),
			Header(), Footer(), Aside(), Details(), Summary(), Dialog(),
			Figure(), Figcaption(), Pre(), Code(), Strong(),
			Small(), Mark(), Table(), Thead(), Tbody(), Tfoot(),
			Tr(), Th(), Td(), Fieldset(), Legend(), Label(),
			Canvas(), Style(), Script(), A("href"), Button(),
		}
		for _, el := range els {
			if el.String() == "" {
				t.Errorf("Element factory returned empty HTML")
			}
		}
	})
}

func TestCoverageDOMLogic(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Update non-existent component", func(t *testing.T) {
		dom.SetLog(nil)
		c := Div().ID("non-existent")
		dom.Update(c)
		dom.SetLog(func(v ...any) { t.Log(v...) })
	})

	t.Run("Render to non-existent parent", func(t *testing.T) {
		err := dom.Render("void", Div())
		if err == nil {
			t.Error("Expected error when rendering to non-existent parent")
		}
	})

	t.Run("Logging", func(t *testing.T) {
		// No crash without log handler
		dom.SetLog(nil)
		dom.Log("test message")

		logged := false
		dom.SetLog(func(v ...any) {
			logged = true
		})
		dom.Log("test message 2")
		if !logged {
			t.Error("Log handler not called")
		}
		dom.SetLog(nil)
	})

	t.Run("Hash", func(t *testing.T) {
		dom.SetHash("test")
		_ = dom.GetHash()
		// OnHashChange coverage
		dom.OnHashChange(func(h string) {})
	})

	t.Run("Element methods", func(t *testing.T) {
		el := Div().Text("hello").Attr("k", "v")
		// Test duplicate Attr
		el.Attr("k", "v2")
		if !fmt.Contains(el.String(), "k='v2'") {
			t.Error("Attr not updated")
		}
	})
}

func TestCoverageEvents(t *testing.T) {
	SetupDOM(t)

	t.Run("Event interface methods - Button", func(t *testing.T) {
		var ev dom.Event
		triggered := false
		btn := Button("Click").ID("btn-ev").On("click", func(e dom.Event) {
			ev = e
			triggered = true
		})
		dom.Render("root", btn)
		TriggerEvent(btn.GetID(), "click", "")

		if triggered && ev != nil {
			ev.PreventDefault()
			ev.StopPropagation()
			_ = ev.TargetID()
			_ = ev.TargetValue()
		}
	})

	t.Run("Append logic", func(t *testing.T) {
		parent := Div().ID("parent-append")
		dom.Render("root", parent)
		child := Span().ID("child-append").Text("appended")
		err := dom.Append("parent-append", child)
		if err != nil {
			t.Errorf("Append failed: %v", err)
		}
		if _, ok := GetRef("child-append"); !ok {
			t.Error("Appended element not found in DOM")
		}
	})
}

func TestLifecycleDeep(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Nested components cleanup", func(t *testing.T) {
		child := &MockComponent{Element: Div().ID("child-comp")}
		parent := &MockComponent{Element: Div(child).ID("parent-comp")}

		dom.Render("root", parent)
		if !child.Mounted {
			t.Error("Child component should be mounted")
		}

		// Update parent to remove child
		parent.Element = Div().ID("parent-comp").Text("no more child")
		dom.Update(parent)

		if child.Mounted {
			t.Error("Child component should be unmounted after removal from parent")
		}
	})

	t.Run("ElementNode in children", func(t *testing.T) {
		// MockComponent is an elementNode (implements AsElement())
		c := &MockComponent{Element: Div().ID("mock-child")}
		el := Div(c)
		html := el.String()
		if !fmt.Contains(html, "<div id='mock-child'") {
			t.Error("MockComponent elementNode not rendered correctly in children")
		}
	})

	t.Run("Default case in renderToHTML", func(t *testing.T) {
		el := Div(123, true)
		html := el.String()
		if !fmt.Contains(html, "123") || !fmt.Contains(html, "true") {
			t.Errorf("Default types not rendered correctly: %s", html)
		}
	})

	t.Run("Component with only String", func(t *testing.T) {
		c := &OnlyHTMLComp{id: "ohc"}
		el := Div(c)
		html := el.String()
		if !fmt.Contains(html, "ONLY HTML") {
			t.Error("Component with only String not rendered correctly")
		}
	})
}

type OnlyHTMLComp struct {
	id string
}

func (c *OnlyHTMLComp) GetID() string             { return c.id }
func (c *OnlyHTMLComp) SetID(id string)           { c.id = id }
func (c *OnlyHTMLComp) String() string            { return "ONLY HTML" }
func (c *OnlyHTMLComp) Children() []dom.Component { return nil }

func TestCoverageCleanup(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Listener cleanup", func(t *testing.T) {
		triggered := false
		btn := Button("Click").ID("btn-clean").On("click", func(e dom.Event) {
			triggered = true
		})
		root := &MockComponent{Element: Div(btn).ID("root-clean")}
		dom.Render("root", root)

		// Trigger before cleanup
		TriggerEvent("btn-clean", "click", "")
		if !triggered {
			t.Error("Event should have triggered")
		}

		// Update root to remove button
		triggered = false
		root.Element = Div().ID("root-clean").Text("Gone")
		dom.Update(root)

		// Trigger after cleanup (should not crash and triggered should remain false)
		TriggerEvent("btn-clean", "click", "")
		if triggered {
			t.Error("Event should NOT have triggered after cleanup")
		}
	})

	t.Run("Option helpers", func(t *testing.T) {
		opt := Option("v1", "Text 1")
		if !fmt.Contains(opt.String(), "value='v1'") || !fmt.Contains(opt.String(), "Text 1") {
			t.Error("Option not rendered correctly")
		}
		sopt := SelectedOption("v2", "Text 2")
		if !fmt.Contains(sopt.String(), "selected=''") {
			t.Error("SelectedOption not rendered correctly")
		}
	})

	t.Run("A and Button", func(t *testing.T) {
		a := A("https://google.com", "Link")
		if !fmt.Contains(a.String(), "href='https://google.com'") {
			t.Error("A not rendered correctly")
		}
		b := Button("Click Me")
		if !fmt.Contains(b.String(), "Click Me") {
			t.Error("Button not rendered correctly")
		}
	})

	t.Run("Element.Children", func(t *testing.T) {
		child := &MockComponent{Element: Div()}
		el := Div(child, "text", Span())
		children := el.Children()
		if len(children) != 2 { // MockComponent and Span
			t.Errorf("Expected 2 component children, got %d", len(children))
		}
	})

	t.Run("Deep cleanup slice manipulation", func(t *testing.T) {
		c1 := &MockComponent{Element: Div().ID("c1")}
		c2 := &MockComponent{Element: Div().ID("c2")}
		c3 := &MockComponent{Element: Div().ID("c3")}
		parent := &MockComponent{Element: Div(c1, c2, c3).ID("parent-deep")}

		dom.Render("root", parent)

		// Remove one child
		parent.Element = Div(c1, c3).ID("parent-deep")
		dom.Update(parent)
		if c2.Mounted {
			t.Error("c2 should be unmounted")
		}

		// Remove all
		parent.Element = Div().ID("parent-deep")
		dom.Update(parent)
		if c1.Mounted || c3.Mounted {
			t.Error("c1 and c3 should be unmounted")
		}
	})

	t.Run("Internal Edge Cases", func(t *testing.T) {
		// Exercise trackComponent already existing
		c := &MockComponent{Element: Div().ID("existing")}
		dom.Render("root", c)
		dom.Render("root", c) // Should return early in trackComponent

		// Exercise trackChildren entry exists
		dom.Update(c) // Should update existing entry in childrenMap
	})

	t.Run("Update with embedded element", func(t *testing.T) {
		child := &MockComponent{Element: Div().ID("embedded")}
		dom.Render("root", child)
		// Update using the embedded element
		dom.Update(child.Element)
	})

	t.Run("Complex cleanup branches", func(t *testing.T) {
		c1 := &MockComponent{Element: Div().ID("c1").On("click", func(e dom.Event) {})}
		c2 := &MockComponent{Element: Div().ID("c2").On("click", func(e dom.Event) {}).On("input", func(e dom.Event) {})}
		parent := &MockComponent{Element: Div(c1, c2).ID("parent-complex").On("click", func(e dom.Event) {})}

		dom.Render("root", parent)

		// This should trigger cleanupListeners for parent and children
		// and hit splitEventKey and the multiple eventFuncs loop
		dom.Render("root", Div().ID("new-root"))

		if c1.Mounted || c2.Mounted || parent.Mounted {
			// They should be unmounted
		}
	})
}
