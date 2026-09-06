# webtyp/html
<img src="docs/img/badges.svg">

HTML element builder API for WebTyp components.

## Overview

`webtyp/html` provides declarative HTML element builders for use in WebTyp components.
It was separated from `webtyp/dom` so that `dom` can focus solely on DOM manipulation and syscall.

## Installation

    go get webtyp.com/html

## Usage (dot-import)

    import (
        . "webtyp.com/html"  // Div, Span, H1, Nav, Button...
        . "webtyp.com/dom"   // Event, Component, Render, Append
    )

    func (c *MyComponent) Render() *dom.Element {
        return Div().Class("container").Child(
            H1().Text("Welcome"),
            P().Text("Minimalist UI."),
            Button().Text("Click").On("click", func(e dom.Event) {}),
        )
    }

## Document shell

`Document` / `DocumentString` emit the full HTML shell that sitec writes as `index.html`:
`<!DOCTYPE html>`, `<html lang>`, `<head>` with `charset`, **viewport**
(`width=device-width, initial-scale=1, viewport-fit=cover`), optional title/favicon/CSS/JS,
and `<body>`. Empty `DocumentOptions{}` already includes the viewport — no field to set.

## SSR HTML Templates

Components that need a custom SSR template implement `RenderHTML() string` in their `html.go` file:

    //go:build !wasm
    package mycomponent

    import . "webtyp.com/html"

    func (c *MyComponent) RenderHTML() string {
        return Div().Set(clsRoot.AsAttr()).String()
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

- [webtyp/dom](https://github.com/webtyp/dom) — DOM manipulation, Element type, lifecycle interfaces
- [webtyp/svg](https://github.com/webtyp/svg) — SVG builders + icon sprite
- [webtyp/image](https://github.com/webtyp/image) — Image builders
