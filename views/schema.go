package views

import (
	"maps"
	"strings"
	"time"

	"mbvlabs/models"
	"mbvlabs/router/routes"
)

type SchemaNode map[string]any
type SchemaBuilder func(HeadData) []SchemaNode

type SchemaBreadcrumb struct {
	Name string
	Path string
}

type schemaListItem struct {
	Name        string
	Path        string
	Description string
	Image       string
}

const (
	organizationDescription = "MBV Labs is a fractional tech lead and engineering studio for founders and lean teams that need senior engineering judgment and hands-on execution."
	defaultServiceAudience  = "Founders, CTOs, product leads, agency owners, and lean technical teams that need senior technical leadership and implementation support."
)

var organizationSameAs = []string{
	"https://www.youtube.com/@mbvlabs",
	"https://x.com/mbvlabs",
	"https://www.linkedin.com/in/mortenvistisen",
	"https://www.linkedin.com/company/89185792/",
	"https://github.com/mbvlabs",
	"https://github.com/mbvisti",
}

func SetSchema(builders ...SchemaBuilder) HeadDataOption {
	return func(hd *HeadData) {
		hd.SchemaBuilders = append(hd.SchemaBuilders, builders...)
	}
}

func buildStructuredData(data HeadData) map[string]any {
	if len(data.SchemaBuilders) == 0 {
		return nil
	}

	siteURL := absoluteSchemaURL(data.canonical, "/")
	orgID := siteURL + "#organization"
	websiteID := siteURL + "#website"
	webpageID := data.canonicalURL + "#webpage"

	graph := []SchemaNode{
		{
			"@type":         "Organization",
			"@id":           orgID,
			"name":          "MBV Labs",
			"alternateName": "mbvlabs",
			"url":           siteURL,
			"email":         "hello@mbvlabs.com",
			"description":   organizationDescription,
			"sameAs":        organizationSameAs,
		},
		{
			"@type":       "WebSite",
			"@id":         websiteID,
			"name":        "MBV Labs",
			"url":         siteURL,
			"description": "Practical engineering notes, consulting services, client work, projects, and technical writing from MBV Labs.",
			"inLanguage":  "en-US",
			"publisher":   schemaRef(orgID),
		},
		{
			"@type":       defaultPageType(data),
			"@id":         webpageID,
			"url":         data.canonicalURL,
			"name":        cleanPageTitle(data),
			"description": data.Description,
			"isPartOf":    schemaRef(websiteID),
			"about":       schemaRef(orgID),
			"publisher":   schemaRef(orgID),
			"inLanguage":  "en-US",
		},
	}

	if data.Image != "" {
		graph[2]["primaryImageOfPage"] = SchemaNode{
			"@type": "ImageObject",
			"url":   absoluteSchemaURL(data.canonical, data.Image),
		}
	}

	for _, builder := range data.SchemaBuilders {
		graph = mergeSchemaNodes(graph, builder(data)...)
	}

	return map[string]any{
		"@context": "https://schema.org",
		"@graph":   graph,
	}
}

func PageSchema(pageType string, breadcrumbs []SchemaBreadcrumb) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		nodes := []SchemaNode{}
		if pageType != "" {
			nodes = append(nodes, SchemaNode{
				"@id":   data.canonicalURL + "#webpage",
				"@type": pageType,
			})
		}
		if len(breadcrumbs) > 0 {
			nodes = append(nodes, breadcrumbSchema(data, breadcrumbs))
		}
		return nodes
	}
}

func HomePageSchema() SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		return []SchemaNode{
			serviceCatalogSchema(data),
			{
				"@type":       "Service",
				"@id":         data.canonicalURL + "#service",
				"name":        "Fractional tech lead services",
				"description": "Senior engineering judgment, hands-on implementation, and delivery tradeoffs for founders and lean teams.",
				"url":         data.canonicalURL,
				"provider":    schemaRef(organizationID(data)),
				"areaServed":  "Worldwide",
				"serviceType": []string{
					"Fractional tech lead",
					"Software engineering consulting",
					"Full-stack development",
					"Internal tools",
					"Workflow automation",
					"Developer infrastructure",
				},
				"audience": SchemaNode{
					"@type":        "BusinessAudience",
					"audienceType": defaultServiceAudience,
				},
			},
		}
	}
}

func ServiceOfferingIndexSchema() SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{
				{Name: "Home", Path: routes.HomePage.URL()},
				{Name: "Services", Path: routes.ServiceOfferingIndex.URL()},
			}),
			serviceCatalogSchema(data),
		}
	}
}

func ContactPageSchema() SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{
				{Name: "Home", Path: routes.HomePage.URL()},
				{Name: "Contact", Path: routes.ProjectInquiryIndex.URL()},
			}),
			{
				"@type":      "ContactPage",
				"@id":        data.canonicalURL + "#webpage",
				"name":       cleanPageTitle(data),
				"url":        data.canonicalURL,
				"about":      schemaRef(organizationID(data)),
				"mainEntity": schemaRef(organizationID(data)),
			},
			{
				"@type": "Organization",
				"@id":   organizationID(data),
				"contactPoint": SchemaNode{
					"@type":             "ContactPoint",
					"email":             "hello@mbvlabs.com",
					"contactType":       "project inquiries",
					"areaServed":        "Worldwide",
					"availableLanguage": "English",
				},
			},
		}
	}
}

func WorkIndexSchema(items []models.WorkEntity) SchemaBuilder {
	return collectionSchema(
		"Work",
		"Client work and engineering experience from MBV Labs.",
		[]SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Work", Path: routes.WorkIndex.URL()}},
		workListItems(items),
	)
}

func WorkShowSchema(work WorkData) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		node := SchemaNode{
			"@type":            "Article",
			"@id":              data.canonicalURL + "#article",
			"headline":         work.Title,
			"description":      work.Summary,
			"url":              data.canonicalURL,
			"mainEntityOfPage": schemaRef(data.canonicalURL + "#webpage"),
			"author":           schemaRef(organizationID(data)),
			"publisher":        schemaRef(organizationID(data)),
			"about":            work.CombinedTags(),
		}
		addImage(node, data.canonical, firstNonEmpty(work.CoverImageUrl, work.ClientLogoUrl))
		addDates(node, work.PublishedAt, work.UpdatedAt, work.CreatedAt)
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Work", Path: routes.WorkIndex.URL()}, {Name: work.Title, Path: routes.WorkShow.URL(work.Slug)}}),
			node,
		}
	}
}

func ProjectIndexSchema(items []models.ProjectEntity) SchemaBuilder {
	return collectionSchema(
		"Projects",
		"Open source and commercial projects built in and through MBV Labs.",
		[]SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Projects", Path: routes.ProjectIndex.URL()}},
		projectListItems(items),
	)
}

func ProjectShowSchema(project ProjectData, technologies []string) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		schemaType := "CreativeWork"
		if project.RepositoryURL != "" {
			schemaType = "SoftwareSourceCode"
		}
		node := SchemaNode{
			"@type":       schemaType,
			"@id":         data.canonicalURL + "#project",
			"name":        project.Name,
			"description": firstNonEmpty(project.Tagline, project.Description),
			"url":         data.canonicalURL,
			"creator":     schemaRef(organizationID(data)),
			"publisher":   schemaRef(organizationID(data)),
			"keywords":    technologies,
		}
		if project.RepositoryURL != "" {
			node["codeRepository"] = project.RepositoryURL
		}
		if project.LiveURL != "" {
			node["sameAs"] = project.LiveURL
		}
		addImage(node, data.canonical, project.ImageURL)
		addDates(node, project.PublishedAt, project.UpdatedAt, project.CreatedAt)
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Projects", Path: routes.ProjectIndex.URL()}, {Name: project.Name, Path: routes.ProjectShow.URL(project.Slug)}}),
			node,
		}
	}
}

func BlogPostIndexSchema(items []models.BlogPostEntity) SchemaBuilder {
	return collectionSchema(
		"Writing",
		"Pragmatic engineering notes from MBV Labs.",
		[]SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Writing", Path: routes.BlogPostIndex.URL()}},
		blogPostListItems(items),
	)
}

func BlogPostShowSchema(post BlogPostData, tags []string) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		node := SchemaNode{
			"@type":            "BlogPosting",
			"@id":              data.canonicalURL + "#article",
			"headline":         post.Title,
			"description":      post.Excerpt,
			"url":              data.canonicalURL,
			"mainEntityOfPage": schemaRef(data.canonicalURL + "#webpage"),
			"author":           schemaRef(organizationID(data)),
			"publisher":        schemaRef(organizationID(data)),
			"keywords":         tags,
		}
		addImage(node, data.canonical, post.CoverImageUrl)
		addDates(node, post.PublishedAt, post.UpdatedAt, post.CreatedAt)
		return []SchemaNode{
			breadcrumbSchema(data, []SchemaBreadcrumb{{Name: "Home", Path: routes.HomePage.URL()}, {Name: "Writing", Path: routes.BlogPostIndex.URL()}, {Name: post.Title, Path: routes.BlogPostShow.URL(post.Slug)}}),
			node,
		}
	}
}

func collectionSchema(name, description string, breadcrumbs []SchemaBreadcrumb, items []schemaListItem) SchemaBuilder {
	return func(data HeadData) []SchemaNode {
		nodes := []SchemaNode{
			breadcrumbSchema(data, breadcrumbs),
			{
				"@type":       "CollectionPage",
				"@id":         data.canonicalURL + "#webpage",
				"name":        name,
				"description": description,
				"url":         data.canonicalURL,
			},
		}
		if len(items) > 0 {
			nodes = append(nodes, itemListSchema(data, name, items))
		}
		return nodes
	}
}

func serviceCatalogSchema(data HeadData) SchemaNode {
	return SchemaNode{
		"@type": "OfferCatalog",
		"@id":   absoluteSchemaURL(data.canonical, routes.ServiceOfferingIndex.URL()) + "#services",
		"name":  "MBV Labs services",
		"url":   absoluteSchemaURL(data.canonical, routes.ServiceOfferingIndex.URL()),
		"itemListElement": []SchemaNode{
			serviceOffer("Shape the technical plan", "Turn fuzzy product or business goals into a practical build plan, architecture choices, delivery risks, and clear tradeoffs before implementation starts.", data),
			serviceOffer("Build production software", "Ship web apps, internal tools, workflow automation, integrations, backend systems, and developer infrastructure with maintainability in view.", data),
			serviceOffer("Improve existing systems", "Audit, modernize, refactor, or extend codebases so the team can move faster without compounding technical debt.", data),
		},
	}
}

func serviceOffer(name, description string, data HeadData) SchemaNode {
	return SchemaNode{
		"@type": "Offer",
		"itemOffered": SchemaNode{
			"@type":       "Service",
			"name":        name,
			"description": description,
			"provider":    schemaRef(organizationID(data)),
			"areaServed":  "Worldwide",
			"audience": SchemaNode{
				"@type":        "BusinessAudience",
				"audienceType": defaultServiceAudience,
			},
		},
	}
}

func breadcrumbSchema(data HeadData, breadcrumbs []SchemaBreadcrumb) SchemaNode {
	items := make([]SchemaNode, 0, len(breadcrumbs))
	for index, breadcrumb := range breadcrumbs {
		items = append(items, SchemaNode{
			"@type":    "ListItem",
			"position": index + 1,
			"name":     breadcrumb.Name,
			"item":     absoluteSchemaURL(data.canonical, breadcrumb.Path),
		})
	}
	return SchemaNode{
		"@type":           "BreadcrumbList",
		"@id":             data.canonicalURL + "#breadcrumb",
		"itemListElement": items,
	}
}

func itemListSchema(data HeadData, name string, items []schemaListItem) SchemaNode {
	elements := make([]SchemaNode, 0, len(items))
	for index, item := range items {
		url := absoluteSchemaURL(data.canonical, item.Path)
		node := SchemaNode{
			"@type":    "ListItem",
			"position": index + 1,
			"url":      url,
			"name":     item.Name,
		}
		if item.Description != "" {
			node["description"] = item.Description
		}
		if item.Image != "" {
			node["image"] = absoluteSchemaURL(data.canonical, item.Image)
		}
		elements = append(elements, node)
	}
	return SchemaNode{
		"@type":           "ItemList",
		"@id":             data.canonicalURL + "#item-list",
		"name":            name,
		"itemListElement": elements,
	}
}

func workListItems(items []models.WorkEntity) []schemaListItem {
	listItems := make([]schemaListItem, 0, len(items))
	for _, item := range items {
		work := newWorkData(item)
		listItems = append(listItems, schemaListItem{
			Name:        work.Title,
			Path:        routes.WorkShow.URL(work.Slug),
			Description: work.Summary,
			Image:       firstNonEmpty(work.CoverImageUrl, work.ClientLogoUrl),
		})
	}
	return listItems
}

func projectListItems(items []models.ProjectEntity) []schemaListItem {
	listItems := make([]schemaListItem, 0, len(items))
	for _, item := range items {
		project := newProjectData(item)
		listItems = append(listItems, schemaListItem{
			Name:        project.Name,
			Path:        routes.ProjectShow.URL(project.Slug),
			Description: project.Tagline,
			Image:       project.ImageURL,
		})
	}
	return listItems
}

func blogPostListItems(items []models.BlogPostEntity) []schemaListItem {
	listItems := make([]schemaListItem, 0, len(items))
	for _, item := range items {
		post := newBlogPostData(item)
		listItems = append(listItems, schemaListItem{
			Name:        post.Title,
			Path:        routes.BlogPostShow.URL(post.Slug),
			Description: post.Excerpt,
			Image:       post.CoverImageUrl,
		})
	}
	return listItems
}

func defaultPageType(data HeadData) string {
	switch data.MetaType {
	case "article":
		return "Article"
	default:
		return "WebPage"
	}
}

func cleanPageTitle(data HeadData) string {
	suffix := " - " + data.siteName
	return strings.TrimSuffix(data.Title, suffix)
}

func organizationID(data HeadData) string {
	return absoluteSchemaURL(data.canonical, "/") + "#organization"
}

func schemaRef(id string) SchemaNode {
	return SchemaNode{"@id": id}
}

func absoluteSchemaURL(base, ref string) string {
	if ref == "" {
		return ""
	}
	return buildCanonicalURL(base, ref)
}

func addImage(node SchemaNode, base, image string) {
	if image != "" {
		node["image"] = absoluteSchemaURL(base, image)
	}
}

func addDates(node SchemaNode, publishedAt, updatedAt, createdAt time.Time) {
	if publishedAt.IsZero() {
		publishedAt = createdAt
	}
	if !publishedAt.IsZero() {
		node["datePublished"] = publishedAt.Format(time.RFC3339)
	}
	if updatedAt.IsZero() {
		updatedAt = publishedAt
	}
	if !updatedAt.IsZero() {
		node["dateModified"] = updatedAt.Format(time.RFC3339)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func mergeSchemaNodes(graph []SchemaNode, nodes ...SchemaNode) []SchemaNode {
	for _, node := range nodes {
		id, ok := node["@id"].(string)
		if !ok || id == "" {
			graph = append(graph, node)
			continue
		}

		merged := false
		for i := range graph {
			if graphID, ok := graph[i]["@id"].(string); ok && graphID == id {
				maps.Copy(graph[i], node)
				merged = true
				break
			}
		}
		if !merged {
			graph = append(graph, node)
		}
	}
	return graph
}
