# PLAN: tinywasm/html — HTML Element Builder API

## Repositorio
`github.com/tinywasm/html` — path local: `tinywasm/html/`  
Estado actual: repo inicializado con `html.go` vacío (solo `type Html struct{}`)

## Dependencias de ejecución
```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

## Prerequisito
Este plan se ejecuta **después** de que `tinywasm/dom` haya publicado su nueva versión v0.10.0 con:
- `RenderHTML()` → `String()` renombrado
- Builders eliminados de `dom/element.go`

---

## Objetivo

`tinywasm/html` provee todos los builders de elementos HTML que vivían en `dom/element.go`. Además introduce el tipo `*HTML` para que los componentes puedan proveer templates SSR tipados (patrón igual a `*css.Stylesheet`).

---

## Dependencia en go.mod

```
tinywasm/html → tinywasm/dom  (para *dom.Element)
tinywasm/dom  → NO importa tinywasm/html
```

Actualizar `tinywasm/html/go.mod`:
```
module github.com/tinywasm/html

go 1.25

require (
    github.com/tinywasm/dom v<nueva-version>
)
```

---

## Archivo: `html/builders.go`

Crear este archivo con TODOS los builders que se eliminaron de `dom/element.go`:

```go
package html

import . "github.com/tinywasm/dom"

// Block containers
func Div(children ...any) *Element        { return (&Element{}).Tag("div").Add(children...) }
func Span(children ...any) *Element       { return (&Element{}).Tag("span").Add(children...) }
func P(children ...any) *Element          { return (&Element{}).Tag("p").Add(children...) }
func Pre(children ...any) *Element        { return (&Element{}).Tag("pre").Add(children...) }
func Code(children ...any) *Element       { return (&Element{}).Tag("code").Add(children...) }
func Strong(children ...any) *Element     { return (&Element{}).Tag("strong").Add(children...) }
func Small(children ...any) *Element      { return (&Element{}).Tag("small").Add(children...) }
func Mark(children ...any) *Element       { return (&Element{}).Tag("mark").Add(children...) }

// Headings
func H1(children ...any) *Element { return (&Element{}).Tag("h1").Add(children...) }
func H2(children ...any) *Element { return (&Element{}).Tag("h2").Add(children...) }
func H3(children ...any) *Element { return (&Element{}).Tag("h3").Add(children...) }
func H4(children ...any) *Element { return (&Element{}).Tag("h4").Add(children...) }
func H5(children ...any) *Element { return (&Element{}).Tag("h5").Add(children...) }
func H6(children ...any) *Element { return (&Element{}).Tag("h6").Add(children...) }

// Lists
func Ul(children ...any) *Element { return (&Element{}).Tag("ul").Add(children...) }
func Ol(children ...any) *Element { return (&Element{}).Tag("ol").Add(children...) }
func Li(children ...any) *Element { return (&Element{}).Tag("li").Add(children...) }

// Semantic layout
func Nav(children ...any) *Element        { return (&Element{}).Tag("nav").Add(children...) }
func Section(children ...any) *Element    { return (&Element{}).Tag("section").Add(children...) }
func Main(children ...any) *Element       { return (&Element{}).Tag("main").Add(children...) }
func Article(children ...any) *Element    { return (&Element{}).Tag("article").Add(children...) }
func Header(children ...any) *Element     { return (&Element{}).Tag("header").Add(children...) }
func Footer(children ...any) *Element     { return (&Element{}).Tag("footer").Add(children...) }
func Aside(children ...any) *Element      { return (&Element{}).Tag("aside").Add(children...) }
func Details(children ...any) *Element    { return (&Element{}).Tag("details").Add(children...) }
func Summary(children ...any) *Element    { return (&Element{}).Tag("summary").Add(children...) }
func Dialog(children ...any) *Element     { return (&Element{}).Tag("dialog").Add(children...) }
func Figure(children ...any) *Element     { return (&Element{}).Tag("figure").Add(children...) }
func Figcaption(children ...any) *Element { return (&Element{}).Tag("figcaption").Add(children...) }

// Tables
func Table(children ...any) *Element { return (&Element{}).Tag("table").Add(children...) }
func Thead(children ...any) *Element { return (&Element{}).Tag("thead").Add(children...) }
func Tbody(children ...any) *Element { return (&Element{}).Tag("tbody").Add(children...) }
func Tfoot(children ...any) *Element { return (&Element{}).Tag("tfoot").Add(children...) }
func Tr(children ...any) *Element    { return (&Element{}).Tag("tr").Add(children...) }
func Th(children ...any) *Element    { return (&Element{}).Tag("th").Add(children...) }
func Td(children ...any) *Element    { return (&Element{}).Tag("td").Add(children...) }

// Form-adjacent (non-form elements)
func Fieldset(children ...any) *Element { return (&Element{}).Tag("fieldset").Add(children...) }
func Legend(children ...any) *Element   { return (&Element{}).Tag("legend").Add(children...) }
func Label(children ...any) *Element    { return (&Element{}).Tag("label").Add(children...) }
func Button(children ...any) *Element   { return (&Element{}).Tag("button").Add(children...) }
func Canvas(children ...any) *Element   { return (&Element{}).Tag("canvas").Add(children...) }
func Style(children ...any) *Element    { return (&Element{}).Tag("style").Add(children...) }
func Script(children ...any) *Element   { return (&Element{}).Tag("script").Add(children...) }

// Special
func A(href string, children ...any) *Element {
    return (&Element{}).Tag("a").Attr("href", href).Add(children...)
}
func Input(typ string) *Element {
    return (&Element{}).Tag("input").Attr("type", typ).SetVoid(true)
}
func Option(value, text string) *Element {
    return (&Element{}).Tag("option").Attr("value", value).Add(text)
}
func SelectedOption(value, text string) *Element {
    return (&Element{}).Tag("option").Attr("value", value).Attr("selected", "").Add(text)
}
func Br() *Element { return (&Element{}).Tag("br").SetVoid(true) }
func Hr() *Element { return (&Element{}).Tag("hr").SetVoid(true) }
```

> **NOTA:** Los métodos `Tag()` y `SetVoid()` deben existir en `dom.Element` o los builders deben replicar la construcción interna. Verificar la API de `dom.Element` al momento de implementar. Si `Element` no expone `Tag()`, usar el constructor interno de dom que ya está disponible: ver `dom/element.go` para la forma exacta de instanciar un `*Element` con un tag. La implementación puede ser: crear función helper privada `el(tag string) *Element` que instancie igual que lo hace dom internamente.

---

## Archivo: `html/html.go`

Reemplazar el contenido actual (que solo tiene `type Html struct{}`) con:

```go
package html

// HTML represents a server-side HTML fragment for SSR injection by assetmin.
// Follows the same pattern as *css.Stylesheet and []*js.Script.
// Build with New() or Raw().
type HTML struct {
    content string
}

// New builds an HTML fragment by serializing an Element tree.
// Use this in component html.go files: RenderHTML() *HTML { return html.New(Div(...)) }
func New(el interface{ String() string }) *HTML {
    return &HTML{content: el.String()}
}

// Raw wraps a pre-built HTML string.
// Use when you have a static template or embedded HTML.
func Raw(content string) *HTML {
    return &HTML{content: content}
}

// String returns the rendered HTML string.
// Called by assetmin extractor via .String() to get the content.
func (h *HTML) String() string {
    if h == nil {
        return ""
    }
    return h.content
}
```

---

## Archivo: `html/providers.go`

```go
package html

// HTMLProvider is an optional capability: components that provide
// a custom SSR HTML template fragment for injection by assetmin.
//
// Implement this in a component's html.go file (//go:build !wasm).
// Only needed when the component has a custom SSR HTML template distinct
// from its Render() output. Most components do NOT need this.
//
// Example in mycomponent/html.go:
//
//  //go:build !wasm
//  package mycomponent
//  import "github.com/tinywasm/html"
//
//  func (c *MyComponent) RenderHTML() *html.HTML {
//      return html.New(Div(clsRoot.AsAttr(), /* ... */))
//  }
type HTMLProvider interface {
    RenderHTML() *html.HTML
}
```

> IMPORTANTE: Este `RenderHTML() *html.HTML` es DISTINTO de `dom.Component.String() string`. El primero es un template SSR custom para assetmin; el segundo es la serialización interna del árbol DOM.

---

## Tests: `html/builders_test.go`

```go
package html_test

import (
    "strings"
    "testing"

    . "github.com/tinywasm/html"
)

func TestDiv_Basic(t *testing.T) {
    el := Div(H2("Hello"), P("World")).Class("root")
    got := el.String()
    if !strings.Contains(got, `class="root"`) { t.Error("expected class root") }
    if !strings.Contains(got, "<h2>Hello</h2>") { t.Error("expected h2") }
    if !strings.Contains(got, "<p>World</p>") { t.Error("expected p") }
}

func TestA_HasHref(t *testing.T) {
    el := A("/home", "Home")
    got := el.String()
    if !strings.Contains(got, `href="/home"`) { t.Error("expected href") }
    if !strings.Contains(got, "Home") { t.Error("expected text") }
}

func TestInput_IsVoid(t *testing.T) {
    el := Input("text")
    got := el.String()
    // void element: no closing tag
    if strings.Contains(got, "</input>") { t.Error("input should be void") }
}

func TestBr_IsVoid(t *testing.T) {
    got := Br().String()
    if strings.Contains(got, "</br>") { t.Error("br should be void") }
}
```

## Tests: `html/html_test.go`

```go
package html_test

import "testing"
import . "github.com/tinywasm/html"

func TestHTML_New(t *testing.T) {
    h := New(Div("content").Class("wrap"))
    if h.String() == "" { t.Fatal("expected non-empty") }
}

func TestHTML_Raw(t *testing.T) {
    raw := `<div class="raw">test</div>`
    h := Raw(raw)
    if h.String() != raw { t.Fatalf("got %q", h.String()) }
}

func TestHTML_Nil_Safe(t *testing.T) {
    var h *HTML
    if h.String() != "" { t.Fatal("nil HTML.String() should return empty") }
}
```

---

## Tests recibidos de tinywasm/dom

Los archivos en `html/tests/` fueron migrados de `dom/tests/`. Están escritos correctamente pero NO compilan hasta que los builders existan.

### Qué necesitan para compilar:

1. El módulo `html/go.mod` debe ser `module github.com/tinywasm/html` y debe tener:
   ```
   require github.com/tinywasm/dom v<version>
   ```

2. Todos los builders referenciados en los tests deben existir en el paquete html:
   - Div, Span, P, H1-H6, Ul, Ol, Li, Nav, Section, Main, Article, Header, Footer, Aside
   - Button, Input, A, Label, Option, SelectedOption, Canvas, Style, Script
   - Details, Summary, Dialog, Figure, Figcaption, Pre, Code, Strong, Small, Mark
   - Table, Thead, Tbody, Tfoot, Tr, Th, Td, Fieldset, Legend, Br, Hr

3. Verificar que los tests pasan:
   ```bash
   cd tinywasm/html
   gotest
   ```

### Tests que cubren DOM lifecycle (usan dom.Render, dom.Update, dom.Append):

- uc_builder_test.go → TestBuilderAndUpdate (Render, Update, Append)
- uc_child_listeners_test.go → TestChildListenersAfterParentUpdate
- uc_coverage_test.go → TestCoverageDOMLogic, TestCoverageEvents, TestLifecycleDeep, TestCoverageCleanup
- uc_declarative_wiring_test.go → declarative component wiring
- uc_double_mount_test.go → TestDoubleMountPrevention
- uc_element_test.go → TestElementMethods
- uc_elm_pattern_test.go → TestElmPattern
- uc_event_test.go → event handling
- uc_fluent_test.go → TestFluentBuilder (pure builder test, no DOM lifecycle)
- uc_focus_preserve_test.go → TestFocusPreservation
- uc_hybrid_render_test.go → TestHybridRendering (String vs Render)
- uc_lifecycle_test.go → TestOnMount, TestOnUnmount
- uc_reference_mutation_test.go → TestReferenceMutation
- uc_render_body_test.go → TestRenderToBody
- uc_selectsearch_test.go → TestSelectSearch simulation
- uc_self_update_test.go → TestSelfUpdate
- uc_update_lifecycle_test.go → TestUpdateLifecycle

### uc_common_test.go

Este archivo define helpers compartidos por todos los demás tests:
- `MockComponent` (implementa dom.Component con String())
- `SetupDOM(t)` — inicializa el DOM WASM para tests
- `TriggerEvent(id, type, value)` — dispara eventos JS
- `GetRef(id)` — obtiene referencia a elemento por ID
- `TestReference` / `MockEvent` — implementaciones de dom.Reference y dom.Event

No hay tests de funcionalidad en este archivo — solo infraestructura.

---

## Ejemplo de referencia: `html/web/client.go`

El archivo `tinywasm/html/web/client.go` fue movido desde `tinywasm/dom/web/client.go`. Es el ejemplo de referencia que muestra el patrón Elm y componentes estáticos usando `tinywasm/html` + `tinywasm/dom`.

Actualizar ese archivo para que:
1. El import cambie de `. "github.com/tinywasm/dom"` a:
   ```go
   import (
       . "github.com/tinywasm/html"  // builders: Div, Span, H1, Button, etc.
       . "github.com/tinywasm/dom"   // lifecycle: Event, Component, Render, Append
   )
   ```
2. Todos los builders (`dom.Div`, `dom.Span`, etc.) pasen a usar el dot import de html: simplemente `Div`, `Span`, etc.
3. Cualquier `.RenderHTML()` → `.String()`

El `dom/README.md` ya fue actualizado para apuntar a `tinywasm/html` — verificar que la URL del README apunte correctamente después de publicar.

---

## Verificación

```bash
cd tinywasm/html
go build ./...
gotest
```

---

## Uso en componentes (patrón post-migración)

```go
// mycomponent/mycomponent.go
import (
    . "github.com/tinywasm/html"  // Div, Span, H1, Button, etc.
    . "github.com/tinywasm/dom"   // Event, Component, Render, Append, etc.
    . "github.com/tinywasm/css"   // Class, Rule, etc.
)

// mycomponent/html.go  (//go:build !wasm — solo si necesita template SSR custom)
import "github.com/tinywasm/html"

func (c *MyComponent) RenderHTML() *html.HTML {
    return html.New(Div(clsRoot.AsAttr()))
}
```

## Documentación a Actualizar

### `html/README.md`

El README actual dice solo "HTML managment api for TinyWasm App". Reemplazar con contenido completo:

```markdown
# tinywasm/html

HTML element builder API for TinyWasm components.

## Overview

`tinywasm/html` provides declarative HTML element builders for use in TinyWasm components.
It was separated from `tinywasm/dom` so that `dom` can focus solely on DOM manipulation and syscall.

## Installation

    go get github.com/tinywasm/html

## Usage (dot-import)

    import (
        . "github.com/tinywasm/html"  // Div, Span, H1, Nav, Button...
        . "github.com/tinywasm/dom"   // Event, Component, Render, Append
    )

    func (c *MyComponent) Render() *dom.Element {
        return Div(
            H1("Welcome"),
            P("Minimalist UI."),
            Button("Click").On("click", func(e dom.Event) {}),
        ).Class("container")
    }

## SSR HTML Templates

Components that need a custom SSR template implement `RenderHTML() *HTML` in their `html.go` file:

    //go:build !wasm
    package mycomponent

    import "github.com/tinywasm/html"

    func (c *MyComponent) RenderHTML() *html.HTML {
        return html.New(Div(clsRoot.AsAttr()))
    }

## Available Builders

Block: Div, Span, P, Pre, Code, Strong, Small, Mark
Headings: H1, H2, H3, H4, H5, H6
Lists: Ul, Ol, Li
Semantic: Nav, Section, Main, Article, Header, Footer, Aside, Details, Summary, Dialog, Figure, Figcaption
Tables: Table, Thead, Tbody, Tfoot, Tr, Th, Td
Form-adjacent: Fieldset, Legend, Label, Button, Canvas, Style, Script
Special: A(href), Input(type), Option, SelectedOption, Br, Hr

## Related Packages

- [tinywasm/dom](https://github.com/tinywasm/dom) — DOM manipulation, Element type, lifecycle interfaces
- [tinywasm/svg](https://github.com/tinywasm/svg) — SVG builders + icon sprite
- [tinywasm/image](https://github.com/tinywasm/image) — Image builders
```

Ver `tinywasm/docs/MASTER_PLAN.md` para el orden global de ejecución.
