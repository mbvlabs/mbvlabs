package handlers

import (
	"strings"
	"testing"
)

func TestMarkdownToHTMLKeepsFormattingWithoutActiveContent(t *testing.T) {
	html := renderPublicMarkdown("# Example\n\n**Bold**\n\n<script>bad()</script>\n\n<a href=\"javascript:bad()\" onclick=\"bad()\">Link</a>", "https://mbvlabs.com")
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

func TestMarkdownFollowsOnlySameOriginLinks(t *testing.T) {
	result := renderPublicMarkdown(
		`[Relative](/blog/rescue) [Absolute](https://mbvlabs.com/work/abc) [Protocol-relative](//mbvlabs.com/blog/guide) [External](https://external.example/guide) [Lookalike](https://mbvlabs.com.external.example/guide) [Other port](https://mbvlabs.com:8443/blog/guide)`,
		"https://mbvlabs.com",
	)
	for _, expected := range []string{
		`<a href="/blog/rescue">Relative</a>`,
		`<a href="https://mbvlabs.com/work/abc">Absolute</a>`,
		`<a href="//mbvlabs.com/blog/guide">Protocol-relative</a>`,
		`<a href="https://external.example/guide" rel="nofollow">External</a>`,
		`<a href="https://mbvlabs.com.external.example/guide" rel="nofollow">Lookalike</a>`,
		`<a href="https://mbvlabs.com:8443/blog/guide" rel="nofollow">Other port</a>`,
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("missing %q from %q", expected, result)
		}
	}
}

func TestMarkdownResolvesAuditedRedirectsWithoutDroppingVideoTime(t *testing.T) {
	result := renderPublicMarkdown(
		`[Video](https://youtu.be/TLyp2O98nQs?si=Q1yeqenP9qh8aM2u&t=300) [OWASP](https://owasp.org/Top10/2025/) [Docs](https://platform.openai.com/docs/guides/function-calling)`,
		"https://mbvlabs.com",
	)
	for _, expected := range []string{
		`href="https://www.youtube.com/watch?v=TLyp2O98nQs&amp;t=300" rel="nofollow"`,
		`href="https://top10.owasp.org/2025/" rel="nofollow"`,
		`href="https://developers.openai.com/api/docs/guides/function-calling" rel="nofollow"`,
	} {
		if !strings.Contains(result, expected) {
			t.Errorf("missing %q from %q", expected, result)
		}
	}
}
