package handlers

import (
	"encoding/json"
	"net/url"
	"strings"

	"mbvlabs/payloads"
)

type HeadData struct {
	siteName       string
	canonical      string
	canonicalURL   string
	structuredData map[string]any
	payloads.PublicSEO
}

func buildCanonicalURL(base, path string) string {
	u, err := url.Parse(base)
	if err == nil {
		ref, err := url.Parse(path)
		if err == nil {
			return u.ResolveReference(ref).String()
		}
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(path, "/")
}

func (h PublicHandler) metadata(title, description, path, image, metaType string, builders ...SchemaBuilder) HeadData {
	titleText := h.cfg.SiteName
	if title != "" {
		titleText = title + " - " + titleText
	}
	if description == "" && title == "" {
		description = "MBV Labs is a fractional tech lead and engineering studio for founders and lean teams."
	}
	imageAlt := title
	if image == "" {
		image = "https://media.mbvlabs.com/covers/mbvlabs-fractional-tech-lead.png"
		imageAlt = "MBV Labs fractional tech lead and engineering studio"
	}
	if metaType == "" {
		metaType = "website"
	}
	base := h.cfg.BaseURL()
	canonical := buildCanonicalURL(base, path)
	data := HeadData{
		siteName: h.cfg.SiteName, canonical: base, canonicalURL: canonical,
		Title: titleText, Description: description, Image: image, ImageAlt: imageAlt,
		Canonical: canonical, MetaType: metaType,
	}
	data.structuredData = buildStructuredData(data, builders...)
	return data
}

func projectSearchDescription(slug, description string) string {
	if slug == "deploycrate-community-edition" {
		return "DeployCrate Community Edition is an open-source platform for automated deployments, application monitoring, and managing your own servers."
	}
	return description
}

func publicStringList(raw json.RawMessage) []string {
	result := make([]string, 0)
	_ = json.Unmarshal(raw, &result)
	return result
}
