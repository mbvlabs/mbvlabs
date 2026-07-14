package config

import "github.com/caarlos0/env/v11"

type firecrawl struct {
	APIKey string `env:"FIRECRAWL_API_KEY"`
}

func newFirecrawlConfig() firecrawl {
	cfg := firecrawl{}

	if err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	return cfg
}
