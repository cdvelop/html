# PLAN: tinywasm/html — Document builder + RewriteAssetURLs (Tier 2 desde assetmin)

## Repositorio
`github.com/tinywasm/html` — path local: `tinywasm/html/`

## Contexto
Refactor de responsabilidades: assetmin delega lógica de tipo-específico a su paquete. Dos
piezas que hoy viven en assetmin pertenecen a `tinywasm/html` (ver assetmin PLAN Cambio 8).

---

## Cambio 1: `Document()` — builder del shell `index.html`

Hoy `assetmin/html.go` (`NewHtmlHandler`) construye el documento HTML a mano (doctype, `<head>`
con links css/js/favicon, `<body>`). Eso es templating de página → pertenece a `tinywasm/html`.

Crear `html/document.go`:
```go
package html

// DocumentOptions configura el shell del documento.
type DocumentOptions struct {
    Lang       string // default "en"
    Title      string
    CSSURL     string // <link rel="stylesheet">
    JSURL      string // <script src defer>
    FaviconURL string // <link rel="icon">
    Head       string // markup extra para <head> (opcional)
}

// Document construye el shell completo del documento como *Element.
// El <body> queda con un punto de inyección (assetmin inserta sprite + HTML de componentes).
func Document(opts DocumentOptions, body ...any) *Element {
    // <!DOCTYPE html><html lang><head>...links...</head><body>...body...</body></html>
}
```

assetmin importará `html` y usará `html.Document(...)` en vez del template a mano. Las URLs
(css/js/favicon) las sigue aportando assetmin; los puntos de inyección dinámica (sprite,
RenderHTML de componentes) los rellena assetmin sobre el resultado.

> Detalle: el doctype no es un `Element` con tag normal. Resolver con un `Element` raíz que
> serialice `<!DOCTYPE html>` como prefijo, o un helper `RawPrefix`. Definir al implementar.

---

## Cambio 2: `RewriteAssetURLs()` — mover desde assetmin

`assetmin/urlRewrite.go` (`rewriteAssetUrls`) reescribe los atributos `href`/`src` de HTML para
reapuntar a un nuevo root. Es manipulación de HTML → pertenece a `tinywasm/html`.

Mover a `html/url_rewrite.go`:
```go
package html

// RewriteAssetURLs reescribe href/src de un fragmento HTML, conservando solo el nombre de
// archivo y anteponiendo newRoot.
func RewriteAssetURLs(htmlStr, newRoot string) string { /* lógica actual de rewriteAssetUrls */ }
```

assetmin llamará `html.RewriteAssetURLs(...)` en `inspect.go` (donde hoy usa `rewriteAssetUrls`).
Eliminar `urlRewrite.go` de assetmin.

---

## Tests
- `html/document_test.go` — verificar doctype, links css/js/favicon, lang, body.
- `html/url_rewrite_test.go` — mover los casos de `assetmin` que cubrían `rewriteAssetUrls`.

## Verificación
```bash
cd tinywasm/html
go build ./...
gotest
gopush
```

Ver `tinywasm/docs/MASTER_PLAN.md` para el orden global.
