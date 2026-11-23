package handlers

import (
	"bytes"
	"html/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"
)

var markdownHTML = bluemonday.UGCPolicy().
	AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6").
	AllowStyles("color", "background-color", "font-weight", "font-style", "text-decoration", "display", "white-space").
	OnElements("span", "pre", "code", "div", "table", "td")

var markdown = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("gruvbox"),
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
				chromahtml.TabWidth(4),
			),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
	),
	goldmark.WithRendererOptions(
		goldhtml.WithHardWraps(),
		goldhtml.WithXHTML(),
		goldhtml.WithUnsafe(),
	),
)

func renderPublicMarkdown(content string) string {
	var out bytes.Buffer
	if err := markdown.Convert([]byte(content), &out); err != nil {
		return template.HTMLEscapeString(content)
	}

	return markdownHTML.Sanitize(out.String())
}
