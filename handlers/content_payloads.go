package handlers

import (
	"mbvlabs/domains"
	"mbvlabs/payloads"
)

func workItem(e domains.WorkEntity) payloads.WorkItem {
	return payloads.WorkItem{
		ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, Title: e.Title, Slug: e.Slug,
		ClientName: payloads.ContentString(e.ClientName), ClientIndustry: payloads.ContentString(e.ClientIndustry),
		ClientURL: payloads.ContentString(e.ClientURL), ClientLogoURL: payloads.ContentString(e.ClientLogoURL),
		Summary: e.Summary, CoverImageURL: payloads.ContentString(e.CoverImageURL),
		Specialisms: e.Specialisms, Platforms: e.Platforms, Technologies: e.Technologies,
		Challenge: e.Challenge, Approach: e.Approach, Deliverables: e.Deliverables, Outcome: e.Outcome, Content: e.Content,
		StartedAt: payloads.ContentTime(e.StartedAt), CompletedAt: payloads.ContentTime(e.CompletedAt),
		Status: string(e.Status), PublishedAt: payloads.ContentTime(e.PublishedAt), IsFeatured: e.IsFeatured,
	}
}

func projectItem(e domains.ProjectEntity) payloads.ProjectItem {
	return payloads.ProjectItem{
		ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, Name: e.Name, Slug: e.Slug, Tagline: e.Tagline,
		Description: payloads.ContentString(e.Description), ProjectType: e.ProjectType, RepositoryURL: payloads.ContentString(e.RepositoryURL),
		LiveURL: payloads.ContentString(e.LiveURL), ImageURL: payloads.ContentString(e.ImageURL), Technologies: e.Technologies,
		StartedAt: payloads.ContentTime(e.StartedAt), LaunchedAt: payloads.ContentTime(e.LaunchedAt),
		PublishedAt: payloads.ContentTime(e.PublishedAt), IsFeatured: e.IsFeatured,
	}
}

func articleItem(e domains.BlogPostEntity) payloads.ArticleItem {
	return payloads.ArticleItem{
		ID: e.ID, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt, Title: e.Title, Slug: e.Slug,
		Excerpt: payloads.ContentString(e.Excerpt), Body: e.Body, Status: e.Status, CoverImageUrl: payloads.ContentString(e.CoverImageUrl),
		Tags: e.Tags, PublishedAt: payloads.ContentTime(e.PublishedAt), PublicationSchedule: e.PublicationSchedule,
	}
}

func mapPayloads[A, B any](items []A, convert func(A) B) []B {
	if items == nil {
		return nil
	}
	result := make([]B, len(items))
	for i, item := range items {
		result[i] = convert(item)
	}
	return result
}
