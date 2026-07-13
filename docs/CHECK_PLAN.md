---
message: "Plan delete `Svg()` / `Use()` placeholders"
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Master plan: https://github.com/tinywasm/tinywasm/blob/main/docs/SVG_ICON_HARNESS_MASTER_PLAN.md
> Execution: tiny — LOCAL possible. No gate (independent of the svg publish).

## Context (zero-context summary)

The ecosystem's harness principle allows exactly ONE construction path per
concept. Sprite references are now built only by `svg.Icon.Render()` in
`tinywasm/svg`. `html/builders.go:78-80` still carries placeholders:

```go
// SVG placeholders to satisfy tests until tinywasm/svg is used
func Svg() *Element { return NewElement("svg") }
func Use() *Element { return NewElement("use").NoCloseTag() }
```

Their stated reason ("until tinywasm/svg is used") has expired. They let code
hand-build `<use href="#typo">`, bypassing the typed path. No non-test code in
the monorepo calls them (verified 2026-07-11).

## Stages

### Stage 1 — delete

Remove the comment and both functions from `builders.go` (lines 76-80 area).
Delete any test that ONLY exercises `Svg()`/`Use()`; if a test uses them
incidentally to build fixtures, rebuild the fixture with `NewElement("svg")` /
`NewElement("use").NoCloseTag()` inline.

### Stage 2 — verification

```bash
gotest                                            # never `go test`
grep -rn "func Svg()\|func Use()" --include='*.go' .   # empty
GOOS=js GOARCH=wasm go build ./...
```

## Anti-footguns

- Do NOT delete `NewElement` or `NoCloseTag` — they are the generic primitives.
- Do NOT touch anything else in `builders.go`.
- Never run `gopush` or `codejob`.

## Stages table

| # | Stage | Files | Done |
|---|---|---|---|
| 1 | Delete placeholders | `builders.go`, affected tests | ☐ |
| 2 | Verify | — | ☐ |
