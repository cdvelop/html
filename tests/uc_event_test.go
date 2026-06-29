//go:build wasm

package html_test

import (
	"syscall/js"
	"testing"

	"github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

// EventComponent registers listeners
type EventComponent struct {
	dom.Element
	clickCount  int
	customCount int
}

func (c *EventComponent) Render() *dom.Element {
	return Div().
		ID(c.GetID()).
		On("click", func(e dom.Event) {
			c.clickCount++
			e.PreventDefault()
			e.StopPropagation()
		}).
		On("custom-test", func(e dom.Event) {
			c.customCount++
		})
}

func TestEvents(t *testing.T) {
	doc := SetupDOM(t)

	t.Run("Basic Event Handling", func(t *testing.T) {
		comp := &MockComponent{}
		comp.SetID("comp-basic-event")
		dom.Render("root", comp)
		el, _ := GetRef("comp-basic-event")

		clicked := false
		el.On("click", func(e dom.Event) {
			clicked = true
		})

		rawEl := doc.Call("getElementById", "comp-basic-event")
		clickEvent := js.Global().Get("MouseEvent").New("click")
		rawEl.Call("dispatchEvent", clickEvent)

		if !clicked {
			t.Error("Click handler not called")
		}
	})

	t.Run("Complex Event Handling and Cleanup", func(t *testing.T) {
		comp := &EventComponent{}
		comp.SetID("comp-events")

		dom.Render("root", comp)

		// Trigger events
		rawEl := doc.Call("getElementById", "comp-events")
		clickEvent := js.Global().Get("MouseEvent").New("click")
		rawEl.Call("dispatchEvent", clickEvent)

		customEvent := js.Global().Get("CustomEvent").New("custom-test")
		rawEl.Call("dispatchEvent", customEvent)

		if comp.clickCount != 1 || comp.customCount != 1 {
			t.Errorf("Events not triggered correctly: %d, %d", comp.clickCount, comp.customCount)
		}

		// Unmount via replacement
		dom.Render("root", Div().ID("cleanup-placeholder"))

		_, found := GetRef("comp-events")
		if found {
			t.Error("Component element still in DOM after Render replacement")
		}
	})

	t.Run("Event Target Value", func(t *testing.T) {
		js.Global().Get("document").Call("getElementById", "root").Set("innerHTML", `<input id="test-input" value="initial">`)
		el, _ := GetRef("test-input")

		var targetVal string
		el.On("input", func(e dom.Event) {
			targetVal = e.TargetValue()
		})

		rawEl := doc.Call("getElementById", "test-input")
		inputEvent := js.Global().Get("Event").New("input", map[string]interface{}{
			"bubbles": true,
		})
		rawEl.Call("dispatchEvent", inputEvent)

		if targetVal != "initial" {
			t.Errorf("Expected target value 'initial', got '%s'", targetVal)
		}
	})
}
