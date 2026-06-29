//go:build wasm

package html_test

import (
	"testing"

	"github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

// RenderableComp implements ViewRenderer for testing the full rendering pipeline.
type RenderableComp struct {
	dom.Element
	label string
}

func (c *RenderableComp) Render() *dom.Element {
	return Div().
		Class("test-comp").
		Child(
			Span().
				ID(c.GetID() + "-label").
				Text(c.label),
		)
}

func (c *RenderableComp) String() string { return c.Element.String() }

func TestBodyHeadResolution(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Append ViewRenderer to body", func(t *testing.T) {
		comp := &RenderableComp{label: "body-append"}
		comp.SetID("body-append-vr")
		err := dom.Append("body", comp)
		if err != nil {
			t.Fatalf("Append to body failed: %v", err)
		}

		el, ok := GetRef("body-append-vr-label")
		if !ok {
			t.Fatal("Appended component not found in body")
		}
		if val := el.GetAttr("id"); val != "body-append-vr-label" {
			t.Errorf("Expected id 'body-append-vr-label', got '%s'", val)
		}
	})

	t.Run("Append to head", func(t *testing.T) {
		comp := &MockComponent{}
		comp.SetID("head-item")
		err := dom.Append("head", comp)
		if err != nil {
			t.Fatalf("Append to head failed: %v", err)
		}

		_, ok := GetRef("head-item")
		if !ok {
			t.Fatal("Component not found in head")
		}
	})
}
