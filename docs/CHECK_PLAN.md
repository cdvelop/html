# tinywasm/html — Plan: Typed Builder (remove `...any` constructors)

> **Master:** tinywasm/docs/PLAN.md · **Engine:** tinywasm/dom/docs/PLAN.md
> **Module:** `github.com/tinywasm/html`
> **Type:** Breaking-aligned migration. Closes the `any` hole in the tag builder.
>
> **Read first:** tinywasm/docs/ARNES_DE_CONSTRUCCION.md — the construction-harness principle this
> serves (typed over `any`, like `tinywasm/json`).

---

## Prerequisites

```bash
# Canonical test runner (WASM tests run against a real DOM). External agents have no global gotest.
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

## Development Rules

- **Documentation First:** update `docs/` and re-index `README.md` before code.
- **No Go stdlib:** use `github.com/tinywasm/fmt`. DOM types come from `github.com/tinywasm/dom`
  (dot-imported, as today). Never `syscall/js`.
- **No new types — reuse `fmt`.** Attributes are `fmt.KeyValue` (already declared in `fmt/parse.go`);
  children are `*dom.Element`. Do **not** declare a `Node`/`Text`/attr wrapper type.
- **Tests:** `gotest` (never `go test`); stdlib assertions only; dual WASM/stdlib. Publish with `gopush 'msg'`.

---

## Context

`builders.go` defines **46 tag constructors** as a generic `any` slot delegating to `Element.Add`:

```go
func Div(children ...any) *Element  { return NewElement("div").Add(children...) }
func Span(children ...any) *Element { return NewElement("span").Add(children...) }
// …44 more (P, Pre, Code, Strong, H1-6, Ul, Ol, Li, Nav, Button, …)
```

The engine plan removes `Element.Add(...any)` (it is the only `any` left in the builder, and it is the
generic slot the harness forbids). So these constructors must change. This is the **house pattern**:
`tinywasm/json`'s writer is already typed-per-primitive with `any` only at the I/O boundary — `html`
aligns to the same shape.

---

## Change 1 — Constructors become no-arg

The 46 `Tag(children ...any) *Element` become **no-arg**, composing via the typed `*Element` methods
from `dom` (`Text`, `Child`, `Attr`, `Class`, `Set(...fmt.KeyValue)`, and the `Bind*` family):

```go
func Div() *Element    { return NewElement("div") }
func Span() *Element   { return NewElement("span") }
func Button() *Element { return NewElement("button") }
// …all 46 the same shape
```

Constructors that take **semantic** arguments stay as-is (the arg is meaningful, not a generic slot):
`Input(typ string)`, `Option(value, text string)`, `SelectedOption(value, text string)`,
`Br()`, `Hr()`.

## Change 2 — Authoring style: nested args → typed methods

Composition moves from constructor arguments to chained typed methods:

```go
// Before (any):
Ul(Li("a"), Li("b"))
Span(clsErr, "error")
Button(clsTsBtn.AsAttr(), "Save")

// After (typed, no any):
Ul().Child(Li().Text("a"), Li().Text("b"))
Span().Class(clsErr).Text("error")
Button().Set(clsTsBtn.AsAttr()).Text("Save")   // Set takes ...fmt.KeyValue (reused type)
```

- **text** → `.Text(string)`
- **child element(s)** → `.Child(...*Element)`
- **attribute / class** → `.Attr(k,v)` / `.Class(...string)` / `.Set(...fmt.KeyValue)`
- **anything reactive** → `.BindText`/`.Bind`/`.BindClass*`/`.BindChildren` (requires a signal)

`document.go`, `providers.go`, `url_rewrite.go`: update any internal use of the old variadic
constructors / `Add` to the typed methods.

## Change 3 — Migrate consumers

Every call site that passed children to a constructor or called `.Add(...)` must move to the typed
methods. Known consumers: `components/*`, `form/*`, `user/*`, `layout/*` (their own plans reference
this). Each consumer plan owns its edits; this plan owns `html` itself.

---

## Documentation (do FIRST)

- **`docs/` + `README.md`:** document the typed builder (no-arg constructors + `.Text/.Child/.Attr/
  .Class/.Set`), state "no `any`, reuse `fmt.KeyValue`", and re-index every `docs/` file.
- **`AGENTS.md` (NEW):** add the construction-harness block (typed over `any`, reuse `fmt` types, no
  new types) — same as the other ecosystem libs.

## Tests — frequent use cases (`gotest`)

- **stdlib:** a small tree built with the typed methods renders the expected HTML string
  (`Div().Class("x").Child(Span().Text("hi"))` → `<div class="x"><span>hi</span></div>`).
- **stdlib:** `Set(...fmt.KeyValue)` applies class/id/attr correctly (same mapping the old `Add` did).
- **build:** the package compiles with **no `...any` constructor** remaining (grep guard in a test or
  doc check).

## Done When

- No `Tag(children ...any)` constructors remain; the 46 are no-arg and the semantic-arg ones unchanged.
- Composition is via typed `*Element` methods; **no `Add(...any)`**; attributes reuse `fmt.KeyValue`
  (no new types declared).
- **Docs:** typed builder documented, `README.md` re-indexed, `AGENTS.md` created.
  **Tests:** the use-case tests pass under `gotest`.
