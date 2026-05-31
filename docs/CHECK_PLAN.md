# PLAN: tinywasm/html — Simplificación v2 (eliminar *HTML)

## Estado previo (CHECK_PLAN ejecutado)

El CHECK_PLAN ya fue ejecutado. El paquete `tinywasm/html` tiene:
- `builders.go` ✅ — todos los builders usando `NewElement()` y `.NoCloseTag()`
- `html.go` ⚠️ — tiene `*HTML` (newtype de string) que se va a eliminar
- `providers.go` ⚠️ — `HTMLProvider` con `RenderHTML() *HTML` que se va a simplificar
- `builders_test.go` ✅ — tests de builders
- `html_test.go` ⚠️ — tests de `*HTML` que se van a eliminar

## Problema

`*HTML` es un newtype de `string` que no agrega valor:

```go
type HTML struct { content string }
func New(el interface{ String() string }) *HTML { return &HTML{el.String()} }
func (h *HTML) String() string { return h.content }
```

`*dom.Element` ya implementa `String() string`. Envolver el resultado en `*HTML` es una indirección sin beneficio. Además `RenderHTML() *HTML` tiene colisión semántica con `String() string` — ambos "renderizan a HTML".

## Objetivo

Eliminar `*HTML`. Simplificar `HTMLProvider` para que devuelva `string` directamente.

---

## Cambio 1: Eliminar `html/html.go`

Eliminar el archivo completamente. Contiene solo `*HTML`, `New()`, `Raw()`.

```bash
rm tinywasm/html/html.go
```

---

## Cambio 2: Eliminar `html/html_test.go`

Los tests prueban `*HTML` que ya no existirá.

```bash
rm tinywasm/html/html_test.go
```

---

## Cambio 3: Simplificar `html/providers.go`

**Reemplazar todo el contenido con:**

```go
package html

// HTMLProvider is an optional capability: components that expose
// a static SSR HTML template fragment for injection by assetmin.
//
// Implement in a component's html.go file (//go:build !wasm).
// Most components do NOT need this — only those with a static shell
// distinct from their dynamic Render() output.
//
// Example in mycomponent/html.go:
//
//	//go:build !wasm
//	package mycomponent
//	import . "github.com/tinywasm/html"
//
//	func (c *MyComponent) RenderHTML() string {
//	    return Div(clsRoot.AsAttr()).String()
//	}
//
// For raw static HTML:
//
//	func (c *MyComponent) RenderHTML() string {
//	    return `<div class="root"></div>`
//	}
type HTMLProvider interface {
    RenderHTML() string
}
```

---

## Cambio 4: Actualizar `assetmin` — detección de HTMLProvider

Ver `tinywasm/assetmin/docs/PLAN.md` Cambio 4 — actualizar para detectar `HTMLProvider` con `RenderHTML()` en lugar de `*html.HTML`.

---

## Cambio 5: Actualizar componentes que implementen RenderHTML()

Buscar en todo el ecosistema:

```bash
grep -r "RenderHTML() \*html.HTML\|RenderHTML() \*HTML\|html\.New(" --include="*.go" .
```

Para cada ocurrencia, migrar:

```go
// ANTES:
func (c *MyComponent) RenderHTML() *html.HTML {
    return html.New(Div(clsRoot.AsAttr()))
}

// DESPUÉS:
func (c *MyComponent) RenderHTML() string {
    return Div(clsRoot.AsAttr()).String()
}
```

---

## Verificación

```bash
cd tinywasm/html
go build ./...
gotest
```

No debe quedar ninguna referencia a `*HTML`, `New()`, `Raw()`, ni `RenderHTML() *HTML`.

---

## Impacto en otros planes

| Paquete | Cambio requerido |
|---------|-----------------|
| `tinywasm/assetmin` | Cambio 4: detectar `RenderHTML()` en vez de `*html.HTML` |
| `tinywasm/components` | Si algún componente implementa `RenderHTML() *HTML` → migrar a `RenderHTML() string` |
| `tinywasm/layout` | Igual que components |
| `tinywasm/form` | Sin impacto (no usa HTMLProvider) |

Ver `tinywasm/docs/MASTER_PLAN.md` para orden global de ejecución.
