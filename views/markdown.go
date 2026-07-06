package views

import (
	"bytes"
	"html/template"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"
)

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

func MarkdownToHTML(content string) string {
	var out bytes.Buffer
	if err := markdown.Convert([]byte(content), &out); err != nil {
		return template.HTMLEscapeString(content)
	}

	return out.String()
}
