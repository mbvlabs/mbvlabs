package handlers

import (
	"strings"
	"testing"
)

func TestMarkdownToHTMLKeepsFormattingWithoutActiveContent(t *testing.T) {
	html := renderPublicMarkdown("# Example\n\n**Bold**\n\n<script>bad()</script>\n\n<a href=\"javascript:bad()\" onclick=\"bad()\">Link</a>")
	for _, want := range []string{`id="example"`, "<strong>Bold</strong>", "Link"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q from %q", want, html)
		}
	}
	for _, unwanted := range []string{"<script", "javascript:", "onclick="} {
		if strings.Contains(html, unwanted) {
			t.Errorf("found %q in %q", unwanted, html)
		}
	}
}
