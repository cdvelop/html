//go:build wasm

package html_test

import (
	"testing"

	"github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

// mountCountComp counts how many times Init is called — used to detect double-mount.
type mountCountComp struct {
	dom.Element
	MountCount int
}

func (c *mountCountComp) Init(ctx dom.Ctx) {
	c.MountCount++
}

func (c *mountCountComp) Render() *dom.Element { return &c.Element }

func TestRender_OnMount_CalledOnce(t *testing.T) {
	SetupDOM(t)

	child := &mountCountComp{}
	child.SetID("mount-child")

	container := Div().Child(child).ID("mount-container")

	if err := dom.Render("root", container); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if child.MountCount != 1 {
		t.Errorf("Init called %d times, want 1 — double-mount bug", child.MountCount)
	}
}

func TestRender_OnMount_MultipleChildren(t *testing.T) {
	SetupDOM(t)

	a := &mountCountComp{}
	a.SetID("mc-a")
	b := &mountCountComp{}
	b.SetID("mc-b")
	c := &mountCountComp{}
	c.SetID("mc-c")

	container := Div().Child(a, b, c).ID("mc-container")

	if err := dom.Render("root", container); err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	for name, comp := range map[string]*mountCountComp{"a": a, "b": b, "c": c} {
		if comp.MountCount != 1 {
			t.Errorf("child %s: Init called %d times, want 1", name, comp.MountCount)
		}
	}
}

func TestAppend_OnMount_CalledOnce(t *testing.T) {
	SetupDOM(t)

	child := &mountCountComp{}
	child.SetID("append-child")
	container := Div().Child(child).ID("append-container")

	if err := dom.Append("root", container); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if child.MountCount != 1 {
		t.Errorf("Append Init called %d times, want 1", child.MountCount)
	}
}
