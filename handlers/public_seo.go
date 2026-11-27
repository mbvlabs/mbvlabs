package handlers

import (
	"context"
	"encoding/xml"
	"strings"
	"time"

	"mbvlabs/domains"
	"mbvlabs/routes"

	"github.com/mbvlabs/andurel/pkg/anduin"
)

func (h PublicHandler) absolute(path string) string {
	return strings.TrimRight(h.cfg.BaseURL(), "/") + "/" + strings.TrimLeft(path, "/")
}

func (h PublicHandler) robots(ctx context.Context) error {
	return anduin.Text(ctx, 200, "User-agent: *\nAllow: /\nContent-Signal: ai-train=no, search=yes, ai-input=yes\nSitemap: "+h.absolute(routes.Sitemap.URL())+"\n")
}

type publicSitemapURL struct {
	Loc        string `xml:"loc"`
	ChangeFreq string `xml:"changefreq"`
	LastMod    string `xml:"lastmod,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}
type publicSitemap struct {
	XMLName xml.Name           `xml:"urlset"`
	XMLNS   string             `xml:"xmlns,attr"`
	URLs    []publicSitemapURL `xml:"url"`
}

func publicLastMod(updated time.Time, published domains.OptionalTime) string {
	if !updated.IsZero() {
		return updated.Format(time.RFC3339)
	}
	if published.Valid {
		return published.Time.Format(time.RFC3339)
	}
	return ""
}
func (h PublicHandler) sitemap(ctx context.Context) error {
	urls := []publicSitemapURL{{h.absolute("/"), "monthly", "2024-10-22T09:43:09+00:00", "1"}}
	for _, path := range []string{routes.ServiceOfferingIndex.URL(), routes.WorkIndex.URL(), routes.ProjectIndex.URL(), routes.BlogPostIndex.URL(), routes.AboutMe.URL(), routes.ProjectInquiryIndex.URL()} {
		urls = append(urls, publicSitemapURL{h.absolute(path), "monthly", "", "0.8"})
	}
	works, err := h.works.AllPublished(ctx)
	if err != nil {
		return err
	}
	for _, item := range works {
		urls = append(urls, publicSitemapURL{h.absolute(routes.WorkShow.URL(item.Slug)), "monthly", publicLastMod(item.UpdatedAt, item.PublishedAt), "0.6"})
	}
	projects, err := h.projects.AllPublished(ctx)
	if err != nil {
		return err
	}
	for _, item := range projects {
		urls = append(urls, publicSitemapURL{h.absolute(routes.ProjectShow.URL(item.Slug)), "monthly", publicLastMod(item.UpdatedAt, item.PublishedAt), "0.6"})
	}
	posts, err := h.posts.AllPublished(ctx)
	if err != nil {
		return err
	}
	for _, item := range posts {
		urls = append(urls, publicSitemapURL{h.absolute(routes.BlogPostShow.URL(item.Slug)), "monthly", publicLastMod(item.UpdatedAt, item.PublishedAt), "0.6"})
	}
	data, err := xml.MarshalIndent(publicSitemap{XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9", URLs: urls}, "", "  ")
	if err != nil {
		return err
	}
	return anduin.Blob(ctx, 200, "application/xml", append([]byte(xml.Header), data...))
}

type publicRSS struct {
	XMLName xml.Name         `xml:"rss"`
	Version string           `xml:"version,attr"`
	Channel publicRSSChannel `xml:"channel"`
}
type publicRSSChannel struct {
	Title       string          `xml:"title"`
	Link        string          `xml:"link"`
	Description string          `xml:"description"`
	Language    string          `xml:"language"`
	Items       []publicRSSItem `xml:"item"`
}
type publicRSSItem struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	GUID        publicRSSGUID `xml:"guid"`
	Description string        `xml:"description"`
	PubDate     string        `xml:"pubDate"`
}
type publicRSSGUID struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

func (h PublicHandler) feed(ctx context.Context) error {
	result, err := h.posts.PaginatePublished(ctx, 1, 25)
	if err != nil {
		return err
	}
	items := make([]publicRSSItem, 0, len(result.BlogPosts))
	for _, post := range result.BlogPosts {
		link := h.absolute(routes.BlogPostShow.URL(post.Slug))
		published := post.CreatedAt
		if post.PublishedAt.Valid {
			published = post.PublishedAt.Time
		}
		items = append(items, publicRSSItem{post.Title, link, publicRSSGUID{true, link}, post.Excerpt.String, published.UTC().Format(time.RFC1123Z)})
	}
	feed := publicRSS{Version: "2.0", Channel: publicRSSChannel{Title: "MBV Labs Blog", Link: h.absolute(routes.BlogPostIndex.URL()), Description: "Pragmatic engineering notes from MBV Labs on product delivery, software architecture, and building maintainable systems.", Language: "en-US", Items: items}}
	data, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return err
	}
	return anduin.Blob(ctx, 200, "application/rss+xml; charset=utf-8", append([]byte(xml.Header), data...))
}
