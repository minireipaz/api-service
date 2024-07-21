package config

type EnvZitadelConfig struct{}

type Config interface {
	GetZitadelURI() string
	GetZitadelKeyUserID() string
	GetZitadelKeyPrivate() string
	GetZitadelKeyID() string
	GetEnv(key, fallback string) string
}

func (e *EnvZitadelConfig) GetZitadelURI() string {
	return GetEnv("ZITADEL_URI", "")
}

func (e *EnvZitadelConfig) GetZitadelKeyUserID() string {
	return GetEnv("ZITADEL_KEY_USERID", "")
}

func (e *EnvZitadelConfig) GetZitadelKeyPrivate() string {
	return GetEnv("ZITADEL_KEY_PRIVATE", "")
}

func (e *EnvZitadelConfig) GetZitadelKeyID() string {
	return GetEnv("ZITADEL_KEY_KEYID", "")
}

func (e *EnvZitadelConfig) GetEnv(key, fallback string) string {
	return GetEnv(key, fallback)
}

func NewZitaldelEnvConfig() Config {
	return &EnvZitadelConfig{}
}
