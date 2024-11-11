package config

type Config struct {
	DSN         string `json:"dsn"`
	Query       string `json:"query"`
	DurationMS  int    `json:"duration_ms"`
	Concurrency int    `json:"concurrency"`
}
