# PLAN — Eliminar `regexp` de código agnóstico (`url_rewrite.go`)

> Este plan se despacha vía el workflow CodeJob. Ver skill: `agents-workflow`.
> **Estado:** PENDIENTE
> **Impacto estimado:** ~35 KB menos en cualquier binario wasm que importe `tinywasm/html`
> (`regexp` + `regexp/syntax` son ~35 KB en TinyGo).

## Prerequisito (PRIMERO — entorno del agente)

El agente NO puede probar en el navegador. Toda verificación es con `gotest` (backend stdlib
+ WASM con build tags), que no viene preinstalado en el entorno aislado del agente:

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

Usar `gotest` (sin argumentos para la suite completa); NO usar `go test` directo.

## Problema

`url_rewrite.go` usa `regexp` **sin build tag** (código agnóstico → compila también en
wasm). Importar `regexp` arrastra `regexp` + `regexp/syntax` (~35 KB en TinyGo), un coste
enorme para un worker de edge que no debería siquiera incluir `html`.

```go
// url_rewrite.go (agnóstico, SIN //go:build)
import "regexp"
var attrStartRegex = regexp.MustCompile(`(?i)\b(src|href)\s*=\s*`)
func RewriteAssetURLs(htmlStr, newRoot string) string { ... }
```

### Evidencia

En `goflare-demo`, `tinywasm/html` entra al edge vía `form/input` (ver
`tinywasm/form/docs/PLAN.md`) y aporta `regexp` (~35 KB) al `edge.wasm`.

## Principio rector

> `regexp` **solo** en código `!wasm` (backend puro). El código agnóstico no debe usar
> `regexp`: se reemplaza por lógica Go explícita. La validación de inputs ya sigue esta
> regla (es regexp-free); `url_rewrite.go` es la excepción a corregir.

## Causa raíz

`RewriteAssetURLs` reescribe URLs relativas de `href`/`src` en un fragmento HTML,
prependiendo un nuevo root. Es una operación de **SSR/backend** (se reescribe el HTML antes
de servirlo); no tiene sentido en el cliente wasm ni en un worker de edge.

## Solución propuesta

**Opción A (recomendada si `RewriteAssetURLs` es solo backend/SSR):** marcar el archivo
`!wasm`.
```go
//go:build !wasm

package html
```
Verificar antes que ningún consumidor wasm la llame (frontend/edge). Si alguno la usa,
es un error de diseño aparte (no debería reescribir assets en el cliente).

**Opción B (si debe ser agnóstico):** reescribir sin `regexp`. El patrón
`(?i)\b(src|href)\s*=\s*` es trivial de escanear a mano: recorrer la cadena buscando,
case-insensitive, los tokens `src`/`href` precedidos de límite de palabra, seguidos de
espacios opcionales, `=`, espacios opcionales. Mantiene el comportamiento sin `regexp`.

> Preferir **A** si la función es SSR-only: es la que aplica el principio más limpiamente
> (concurrencia/parseo pesado pertenece al backend, igual que `sync`/`reflect` ya están
> tras `!wasm` en `tinywasm/fmt`).

## Pasos

1. Confirmar que `RewriteAssetURLs` no tiene llamadores en código wasm:
   ```bash
   grep -rn "RewriteAssetURLs" --include=*.go ~/Dev/Project/tinywasm | grep -v _test
   ```
2. Si es backend-only → añadir `//go:build !wasm` a `url_rewrite.go`.
3. Si es agnóstico → reemplazar el regex por un escáner manual en Go (sin `regexp`).
4. Quitar el import `"regexp"` del archivo agnóstico.

## Verificación

```bash
# En un proyecto que importe html en wasm:
GOOS=js GOARCH=wasm go list -deps <pkg> | grep -E '^regexp$'   # → vacío
grep -rln '"regexp"' *.go | xargs -I{} sh -c 'head -1 {} | grep -q "go:build !wasm" || echo "AGNOSTIC: {}"'
# → no debe listar url_rewrite.go
```

## Nota

Una vez que `tinywasm/form` desacople metadata de rendering (ver su `docs/PLAN.md`),
`tinywasm/html` dejará de entrar al edge. Aun así, esta corrección es necesaria por
principio (evita la fuga de `regexp` en el frontend y en futuros consumidores).
