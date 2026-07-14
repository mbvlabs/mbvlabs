package config

import "github.com/caarlos0/env/v11"

type serper struct {
	APIKey string `env:"SERPER_API_KEY"`
}

func newSerperConfig() serper {
	cfg := serper{}

	if err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		panic(err)
	}

	return cfg
}
