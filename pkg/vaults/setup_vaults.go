package vaults

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvsFromVault() {
	envValues := GetAllEnvsFromRedis()
  if envValues == "" {
    return
  }
	envsMaped := serializeEnvs(envValues)
	setEnvs(envsMaped)
}

func serializeEnvs(envStr string) map[string]string {
	envMap, err := godotenv.Unmarshal(envStr)
	if err != nil {
		log.Panic("ERROR | Cannot serialize string from Env")
	}
	return envMap
}

func setEnvs(envsMapped map[string]string) {
	for key, value := range envsMapped {
		os.Setenv(key, value)
	}
}
