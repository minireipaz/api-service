package config

type EnvClickhouseConfig struct{}

type ClickhouseConfig interface {
	GetClickhouseURI() string
	GetClickhouseToken() string
	GetEnv(key, fallback string) string
}

func (e *EnvClickhouseConfig) GetClickhouseURI() string {
	return GetEnv("CLICKHOUSE_API_URI", "")
}

func (e *EnvClickhouseConfig) GetClickhouseToken() string {
	return GetEnv("CLICKHOUSE_TOKEN_PIPES", "")
}

func (e *EnvClickhouseConfig) GetEnv(key, fallback string) string {
	return GetEnv(key, fallback)
}

func NewClickhouseEnvConfig() ClickhouseConfig {
	return &EnvClickhouseConfig{}
}
