//go:build wasm

package html_test

import (
	"syscall/js"
	"testing"

	"github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

// FocusUpdater is a component with a text input that uses Bind to handle reactivity.
type FocusUpdater struct {
	dom.Element
	filterTerm *dom.SignalString
}

func (c *FocusUpdater) Init(ctx dom.Ctx) {
	c.filterTerm = dom.NewString("")
}

func (c *FocusUpdater) Render() *dom.Element {
	return Div().Child(
		Input("text").ID(c.GetID()+"-input").Bind(c.filterTerm),
	)
}

// installFocusSpy replaces HTMLElement.prototype.focus with a counter spy.
func installFocusSpy(t *testing.T) (count *int, cleanup func()) {
	t.Helper()
	n := 0
	proto := js.Global().Get("HTMLElement").Get("prototype")
	original := proto.Get("focus")
	spy := js.FuncOf(func(this js.Value, args []js.Value) any {
		n++
		original.Call("call", this) // invoke original with correct this binding
		return nil
	})
	proto.Set("focus", spy)
	return &n, func() {
		proto.Set("focus", original)
		spy.Release()
	}
}

func TestUpdate_FocusCalledOnceWhenActive(t *testing.T) {
	SetupDOM(t)

	c := &FocusUpdater{}
	if err := dom.Render("root", c); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	id := c.GetID()
	doc := js.Global().Get("document")

	inputEl := doc.Call("getElementById", id+"-input")
	inputEl.Call("focus")
	inputEl.Set("value", "a")

	_, cleanup := installFocusSpy(t)
	defer cleanup()

	// Dispatch an input event — this triggers Bind's handler → c.filterTerm.Set()
	event := js.Global().Get("InputEvent").New("input", map[string]any{"bubbles": true})
	inputEl.Call("dispatchEvent", event)
}

func TestUpdate_PreservesActiveElementFocus(t *testing.T) {
	SetupDOM(t)

	c := &FocusUpdater{}
	if err := dom.Render("root", c); err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	id := c.GetID()
	doc := js.Global().Get("document")

	inputEl := doc.Call("getElementById", id+"-input")
	if inputEl.IsNull() {
		t.Fatal("input element not found")
	}
	inputEl.Call("focus")
	inputEl.Set("value", "hi")
	inputEl.Call("setSelectionRange", 2, 2)

	event := js.Global().Get("InputEvent").New("input", map[string]any{
		"bubbles": true,
	})
	inputEl.Call("dispatchEvent", event)

	activeID := doc.Get("activeElement").Get("id").String()
	if activeID != id+"-input" {
		t.Errorf("focus lost after update: activeElement.id=%q, want %q", activeID, id+"-input")
		return
	}

	selStart := doc.Get("activeElement").Get("selectionStart").Int()
	if selStart != 2 {
		t.Errorf("cursor reset after update: selectionStart=%d, want 2", selStart)
	}
}
