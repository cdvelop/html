---
PLAN: "fix: el shell de Document() debe emitir el meta viewport"
TAG: v0.0.7
---

## Antes de escribir código: lee [AGENTS.md](../AGENTS.md)

**Es vinculante, no orientativo.**

| # | Principio | Cómo se aplica aquí |
|---|---|---|
| 3 | Illegal states unrepresentable | Un documento sin `<meta name="viewport">` es un estado ilegal en móvil. Deja de ser escribible: no hay opción para omitirlo |
| 8 | Closed by default | El valor correcto es el que sale **escribiendo nada**. No se añade un campo `Viewport` que haya que acordarse de rellenar |
| 9 | Lego pieces, never forks | El shell del documento lo emite esta librería. `assetmin` no debe volver a escribirlo a mano, y `devbrowser` no debe inyectarlo en caliente |
| 4 | One way to do each thing | Una sola etiqueta viewport en el ecosistema, con un solo valor |

---

## 1. El hueco

`Document()` construye el `<head>` en `document.go:24-40`. Emite `charset`,
`title`, favicon, hoja de estilo, script y el `Head` extra del llamador.

**No emite `<meta name="viewport">`.**

```go
// document.go:24-25
head := NewElement("head")
head.Child(NewElement("meta").NoCloseTag().Attr("charset", "utf-8"))
// ...y aquí nunca aparece el viewport
```

### 1.1 Qué provoca exactamente

Sin esa etiqueta, Safari iOS —y todo navegador móvil— asume un viewport virtual
de **980px** y escala la página entera para que quepa en la pantalla. El
resultado no es "un poco distinto": es la app completa reducida, con texto
ilegible, como una web de escritorio vista desde el teléfono. Es la diferencia
más severa y más fácil de confundir con un problema de CSS, porque el CSS está
bien: nunca se le pidió al navegador que usara el ancho del dispositivo.

Y es silenciosa: en el escritorio, donde se desarrolla, no se nota nada.

### 1.2 Por qué esta librería y no otra

El comentario de `document.go:48-49` ya lo dice:

```go
// DocumentString renders the full HTML document including <!DOCTYPE html>.
// This is what assetmin writes to disk as index.html.
```

Esta librería es la dueña declarada del shell. Que hoy `assetmin` además
mantenga su propia copia escrita a mano (`assetmin/html.go:44`) es una deuda de
principio 9 que se corrige en el plan de aquel repo; no cambia quién es el dueño.

Tampoco puede resolverlo `tinywasm/css`: una etiqueta `<meta>` no es CSS. Ni
`tinywasm/devbrowser`: inyectarla en el navegador de desarrollo haría que el
emulador mostrara un HTML que el servidor real no sirve — el peor tipo de
diferencia, la que sólo aparece en producción.

---

## 2. Cambio

Una línea en `Document()`, inmediatamente después del `charset` —el orden
importa: ambas metaetiquetas deben estar en los primeros bytes del `<head>`:

```go
head.Child(NewElement("meta").NoCloseTag().
    Attr("name", "viewport").
    Attr("content", "width=device-width, initial-scale=1, viewport-fit=cover"))
```

### 2.1 Los tres valores, uno por uno

| Valor | Qué hace | Qué pasa sin él |
|---|---|---|
| `width=device-width` | El viewport CSS mide lo que mide la pantalla | Viewport virtual de 980px: la app se ve reducida |
| `initial-scale=1` | Sin zoom inicial | En iOS, la orientación horizontal puede entrar con zoom aplicado |
| `viewport-fit=cover` | El documento ocupa **también** el área bajo el notch/Dynamic Island, y **activa** `env(safe-area-inset-*)` | Los `env(safe-area-inset-*)` valen **0px** en todos los dispositivos. Cualquier CSS de safe areas queda como código muerto |

`viewport-fit=cover` es el que menos se conoce y el que más consecuencias tiene
aguas abajo: es la **precondición** de los tokens de safe area que
`tinywasm/css` va a declarar. Sin esta línea, aquel trabajo no tiene efecto
observable.

### 2.2 Lo que no se añade

**No** se añade `user-scalable=no` ni `maximum-scale=1`. Bloquean el zoom del
usuario: es un fallo de accesibilidad (WCAG 1.4.4), y desde iOS 10 Safari lo
ignora de todos modos. Un valor que el navegador ignora es ruido que aparenta
control.

**No** se añade un campo `Viewport` a `DocumentOptions`. Por principio 8, el
valor correcto se obtiene escribiendo nada; un campo opcional convierte "la app
se ve bien en el móvil" en algo que hay que acordarse de activar, y el modo de
fallo es silencioso. Quien necesite un viewport exótico ya tiene `opts.Head`.

---

## 3. Verificación

En `tests/`, siguiendo el patrón `uc_*_test.go` y con aserciones de stdlib:

1. `DocumentString(DocumentOptions{})` —opciones vacías— contiene
   `name="viewport"` y contiene `viewport-fit=cover`. Es el test que fija el
   principio 8: sale bien sin configurar nada.
2. El `<meta charset>` aparece **antes** que el viewport, y ambos antes que
   cualquier `<link>` o `<script>`.
3. La salida contiene exactamente **una** etiqueta viewport (guarda contra que
   un consumidor la duplique vía `opts.Head` sin darse cuenta; si aparece dos
   veces, el test lo dice).
4. La salida **no** contiene `user-scalable` ni `maximum-scale` (guarda de
   accesibilidad: convierte §2.2 en algo que un cambio futuro no puede deshacer
   en silencio).

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
gotest
```

---

## 4. Alcance

### Dentro

- `document.go`: la etiqueta en `Document()`.
- `tests/`: los cuatro tests de §3.
- `README.md`: mencionar qué incluye el shell que produce `Document()`.

### Fuera

- El `index.html` que `assetmin` escribe a mano y su plantilla: son de aquel
  repo, tienen su propio plan.
- Los tokens de safe area (`env(safe-area-inset-*)`): son de `tinywasm/css`.
  Este plan sólo crea la condición para que funcionen.
- Reconciliar las dos implementaciones del shell (que `assetmin` consuma
  `DocumentString`): es la corrección de fondo del principio 9, más grande que
  este cambio y con otro dueño.

### Orden respecto a los planes hermanos

Este plan y el de `assetmin` son **prerrequisitos** del de `css`: los tokens de
safe area no producen ningún efecto observable hasta que el documento sale con
`viewport-fit=cover`. Se pueden ejecutar en paralelo, pero el de `css` no se
puede *verificar en un dispositivo* antes que estos dos.
