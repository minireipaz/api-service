package config

type EnvZitadelConfig struct{}

type ZitadelConfig interface {
	GetZitadelURI() string
	GetZitadelServiceUserID() string
	GetZitadelServiceUserKeyPrivate() string
	GetZitadelServiceUserKeyID() string
	GetZitadelServiceUserClientID() string
	GetZitadelBackendID() string
	GetZitadelBackendKeyPrivate() string
	GetZitadelBackendKeyID() string
	GetZitadelKeyClientID() string
	GetZitadelProjectID() string
	GetZitadelBackendClientID() string
	GetEnv(key, fallback string) string
}

func (e *EnvZitadelConfig) GetZitadelURI() string {
	return GetEnv("ZITADEL_URI", "")
}

func (e *EnvZitadelConfig) GetZitadelProjectID() string {
	return GetEnv("ZITADEL_PROJECTID", "")
}

func (e *EnvZitadelConfig) GetZitadelKeyClientID() string {
	return GetEnv("ZITADEL_KEY_CLIENTID", "")
}

func (e *EnvZitadelConfig) GetZitadelBackendClientID() string {
	return GetEnv("ZITADEL_KEY_CLIENTID", "")
}

func (e *EnvZitadelConfig) GetZitadelServiceUserID() string {
	return GetEnv("ZITADEL_KEY_USERID_SERVICE_ACCOUNT", "")
}

func (e *EnvZitadelConfig) GetZitadelServiceUserKeyPrivate() string {
	return GetEnv("ZITADEL_KEY_PRIVATE_SERVICE_ACCOUNT", "")
}

func (e *EnvZitadelConfig) GetZitadelServiceUserKeyID() string {
	return GetEnv("ZITADEL_KEY_KEYID_SERVICE_ACCOUNT", "")
}

func (e *EnvZitadelConfig) GetZitadelServiceUserClientID() string {
	return GetEnv("ZITADEL_KEY_CLIENTID_SERVICE_ACCOUNT", "")
}

func (e *EnvZitadelConfig) GetEnv(key, fallback string) string {
	return GetEnv(key, fallback)
}

func (e *EnvZitadelConfig) GetZitadelBackendID() string {
	return GetEnv("ZITADEL_KEY_APP_ID_BACKEND", "")
}

func (e *EnvZitadelConfig) GetZitadelBackendKeyPrivate() string {
	return GetEnv("ZITADEL_KEY_PRIVATE_BACKEND", "")
}

func (e *EnvZitadelConfig) GetZitadelBackendKeyID() string {
	return GetEnv("ZITADEL_KEY_KEYID_BACKEND", "")
}

func NewZitaldelEnvConfig() ZitadelConfig {
	return &EnvZitadelConfig{}
}
