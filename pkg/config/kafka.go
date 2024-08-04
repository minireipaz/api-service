package config

type EnvKafkaConfig struct{}

type KafkaConfig interface {
	GetServersURI() string
	GetProtocol() string
	GetMechanisms() string
	GetUsername() string
	GetPassword() string
	GetTimeout() string
	GetEnv(key, fallback string) string
}

func (e *EnvKafkaConfig) GetServersURI() string {
	return GetEnv("bootstrap.servers", "localhost:9092")
}

func (e *EnvKafkaConfig) GetProtocol() string {
	return GetEnv("security.protocol", "SASL_SSL")
}

func (e *EnvKafkaConfig) GetMechanisms() string {
	return GetEnv("sasl.mechanisms", "PLAIN")
}

func (e *EnvKafkaConfig) GetUsername() string {
	return GetEnv("sasl.username", "")
}

func (e *EnvKafkaConfig) GetPassword() string {
	return GetEnv("sasl.password", "")
}

func (e *EnvKafkaConfig) GetTimeout() string {
	return GetEnv("session.timeout.ms", "45000")
}

func (e *EnvKafkaConfig) GetEnv(key, fallback string) string {
	return GetEnv(key, fallback)
}

func NewKafkaEnvConfig() KafkaConfig {
	return &EnvKafkaConfig{}
}
