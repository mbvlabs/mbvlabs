package handlers

import "slices"

// Group existing guides by the decision a reader is trying to make.
var relatedArticleGroups = [][]string{
	{"software-project-rescue", "codebase-audit", "software-project-handover-checklist", "should-i-rewrite-my-legacy-codebase"},
	{"software-project-discovery-phase", "how-to-scope-an-mvp", "how-to-choose-tech-stack-for-startup", "what-is-a-fractional-tech-lead"},
	{"agentic-engineering-is-not-about-writing-more-code", "ai-assisted-software-development-engineering-priorities", "api-first-ai-agents"},
	{"build-vs-buy-internal-tools", "what-business-process-to-automate-first", "custom-internal-tools-development"},
}

func relatedArticleSlugs(slug string) []string {
	for _, group := range relatedArticleGroups {
		if slices.Contains(group, slug) {
			var related []string
			for _, candidate := range group {
				if candidate != slug {
					related = append(related, candidate)
				}
			}
			return related
		}
	}
	return nil
}
