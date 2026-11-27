package config

import "github.com/mbvlabs/andurel/pkg/env"

type SiteAPI struct {
	Username     string
	Password     string
	SerperKey    string
	FirecrawlKey string
}

func newSiteAPI() SiteAPI {
	e := env.New()
	return SiteAPI{
		Username:     e.String("API_BASIC_AUTH_USERNAME", ""),
		Password:     e.String("API_BASIC_AUTH_PASSWORD", ""),
		SerperKey:    e.String("SERPER_API_KEY", ""),
		FirecrawlKey: e.String("FIRECRAWL_API_KEY", ""),
	}
}
