package handlers

import (
	"bytes"
	"html/template"
	"io"
	"net/url"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldhtml "github.com/yuin/goldmark/renderer/html"
	"golang.org/x/net/html"
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

func renderPublicMarkdown(content, siteURL string) string {
	var out bytes.Buffer
	if err := markdown.Convert([]byte(content), &out); err != nil {
		return template.HTMLEscapeString(content)
	}

	return rewritePublicLinks(markdownHTML.Sanitize(out.String()), siteURL)
}

func rewritePublicLinks(markup, siteURL string) string {
	site, err := url.Parse(siteURL)
	if err != nil || site.Host == "" {
		return markup
	}
	tokens := html.NewTokenizer(strings.NewReader(markup))
	var out strings.Builder
	out.Grow(len(markup))
	for {
		switch tokens.Next() {
		case html.ErrorToken:
			if tokens.Err() != io.EOF {
				return markup
			}
			return out.String()
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokens.Token()
			if token.Data == "a" {
				for i := range token.Attr {
					if token.Attr[i].Key != "href" {
						continue
					}
					href := token.Attr[i].Val
					if sameSiteLink(site, href) {
						// UGCPolicy adds nofollow to all links. Keep that policy
						// for other sites, but make our own editorial links crawlable.
						token.Attr = withoutNofollow(token.Attr)
					} else {
						token.Attr[i].Val = canonicalOutboundURL(href)
					}
					break
				}
			}
			out.WriteString(token.String())
			continue
		}
		out.Write(tokens.Raw())
	}
}

func sameSiteLink(site *url.URL, href string) bool {
	href = strings.TrimSpace(href)
	ref, err := url.Parse(href)
	if err != nil || ref.User != nil || ref.Opaque != "" || strings.Contains(href, `\`) ||
		(strings.HasPrefix(href, "//") && ref.Host == "") {
		return false
	}
	target := site.ResolveReference(ref)
	return strings.EqualFold(target.Scheme, site.Scheme) && strings.EqualFold(target.Host, site.Host)
}

func withoutNofollow(attrs []html.Attribute) []html.Attribute {
	for i := range attrs {
		if attrs[i].Key != "rel" {
			continue
		}
		kept := make([]string, 0)
		for rel := range strings.FieldsSeq(attrs[i].Val) {
			if !strings.EqualFold(rel, "nofollow") {
				kept = append(kept, rel)
			}
		}
		if len(kept) == 0 {
			return append(attrs[:i], attrs[i+1:]...)
		}
		attrs[i].Val = strings.Join(kept, " ")
		break
	}
	return attrs
}

// Retarget known editorial links without rewriting stored Markdown or changing
// unrelated URLs. These are the final destinations of the audited redirects.
func canonicalOutboundURL(href string) string {
	switch href {
	case "https://www.candy.io/":
		return "https://candy.io/"
	case "https://chatsheet.com", "https://chatsheet.com/":
		return "https://getspectra.com/"
	case "https://platform.openai.com/docs/guides/function-calling":
		return "https://developers.openai.com/api/docs/guides/function-calling"
	case "https://owasp.org/Top10/2025/":
		return "https://top10.owasp.org/2025/"
	case "https://owasp.org/www-project-application-security-verification-standard/":
		return "https://asvs.owasp.org/"
	case "https://owasp.org/API-Security/editions/2023/en/0x11-t10/":
		return "https://api-security.owasp.org/editions/2023/en/0x11-t10/"
	case "https://youtu.be/TLyp2O98nQs?si=Q1yeqenP9qh8aM2u&t=300":
		return "https://www.youtube.com/watch?v=TLyp2O98nQs&t=300"
	}
	return href
}
