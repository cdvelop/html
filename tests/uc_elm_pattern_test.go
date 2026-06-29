//go:build wasm

package html_test

import (
	"testing"

	"github.com/tinywasm/dom"
	. "github.com/tinywasm/html"
)

type CounterElm struct {
	dom.Element
	count *dom.SignalString
}

func (c *CounterElm) Init(ctx dom.Ctx) {
	c.count = dom.NewString("0")
}

func (c *CounterElm) Render() *dom.Element {
	return Div().Child(
		Span().ID("count-val").BindText(c.count),
	)
}

func TestElmPattern(t *testing.T) {
	_ = SetupDOM(t)

	t.Run("State Update and Re-render", func(t *testing.T) {
		c := &CounterElm{}
		c.SetID("counter-elm")
		dom.Render("root", c)

		// Check initial render
		_, ok := GetRef("count-val")
		if !ok {
			t.Fatal("Counter value not found")
		}

		// Perform update via signal
		c.count.Set("1")

		// Verify re-render occurred (no error)
		_, ok = GetRef("count-val")
		if !ok {
			t.Fatal("Counter value lost after update")
		}
	})
}
