package config

import (
	"log"
	"minireipaz/pkg/vaults"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvs() {
	if err := godotenv.Load(); err != nil {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Panic("ERROR | Initial LoadEnvs FAILED")
		}
	}
	LoadEnvsFromVault()
}

func LoadEnvsFromVault() {
	vaults.GetEnvsFromVault()
}

func GetEnv(key, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}
