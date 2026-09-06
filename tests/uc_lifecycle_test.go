//go:build wasm

package html_test

import (
	"testing"

	"webtyp.com/dom"
	. "webtyp.com/html"
)

func TestLifecycle(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("Get Existing Element", func(t *testing.T) {
		el, ok := GetRef("root")
		if !ok {
			t.Fatal("Expected to find root element")
		}
		if el.GetAttr("id") != "root" {
			t.Errorf("Expected id 'root', got %s", el.GetAttr("id"))
		}
	})

	t.Run("Get Non-Existing Element", func(t *testing.T) {
		_, ok := GetRef("non-existent")
		if ok {
			t.Error("Expected not to find non-existent element")
		}
	})

	t.Run("Get Body and Head", func(t *testing.T) {
		body, ok := GetRef("body")
		if !ok {
			t.Error("Expected to find body")
		}
		if body == nil {
			t.Error("Body elem should not be nil")
		}

		head, ok := GetRef("head")
		if !ok {
			t.Error("Expected to find head")
		}
		if head == nil {
			t.Error("Head elem should not be nil")
		}
	})

	t.Run("Mount Component", func(t *testing.T) {
		comp := &MockComponent{Element: *Div()}
		comp.SetID("comp1")
		err := dom.Render("root", comp)
		if err != nil {
			t.Fatalf("Mount failed: %v", err)
		}

		if !comp.Mounted {
			t.Error("Init was not called")
		}

		el, ok := GetRef("comp1")
		if !ok {
			t.Fatal("Component element not found in DOM")
		}
		if val := el.Value(); val != "" && val != "<undefined>" && val != "undefined" {
			t.Errorf("Expected empty value or <undefined> for div, got: %s", val)
		}
	})

	t.Run("Unmount Component", func(t *testing.T) {
		comp := &MockComponent{Element: *Div()}
		comp.SetID("comp2")
		dom.Render("root", comp)

		// Unmount via replacement
		dom.Render("root", Div().ID("root-placeholder"))

		if comp.Mounted {
			t.Error("Component should be unmounted after replacement")
		}

		_, ok := GetRef("comp2")
		if ok {
			t.Error("Element should be removed from DOM")
		}
	})

	t.Run("Mount Invalid Parent", func(t *testing.T) {
		comp := &MockComponent{Element: *Div()}
		comp.SetID("comp-invalid")
		err := dom.Render("invalid-parent-id", comp)
		if err == nil {
			t.Error("Expected error when mounting to invalid parent")
		}
	})

	t.Run("Get Cache Hit", func(t *testing.T) {
		_, ok := GetRef("root")
		if !ok {
			t.Fatal("Root not found")
		}

		el, ok := GetRef("root")
		if !ok {
			t.Fatal("Root not found in cache")
		}
		if el.GetAttr("id") != "root" {
			t.Error("Wrong element from cache")
		}
	})
}
