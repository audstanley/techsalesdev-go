package handlers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func checkForEmptyEnvironment(s string) bool {
	if s == "" {
		return true
	} else {
		return false
	}
}

func GetEnvironment() map[string]string {
	envs := map[string]string{
		"API_DEV_PORT":        goDotEnvVariable("API_DEV_PORT"),
		"ACCESS_TOKEN_SECRET": goDotEnvVariable("ACCESS_TOKEN_SECRET"),
		"SMTP_ACCOUNT":        goDotEnvVariable("SMTP_ACCOUNT"),
		"SMTP_PASS":           goDotEnvVariable("SMTP_PASS"),
		"REDIS_PASS":          goDotEnvVariable("REDIS_PASS"),
		"REDIS_ENDPOINT":      goDotEnvVariable("REDIS_ENDPOINT"),
	}
	for k, v := range envs {
		if checkForEmptyEnvironment(v) {
			if k == "SMTP_ACCOUNT" || k == "SMTP_PASS" {
				if !DisableSendingEmail {
					fmt.Sprintf("NEED %s Credentials, otherwise endpoints that envolve sending an email will crash the backend.\n", k)
					fmt.Println("If you don't set the SMTP_ACCOUNT and SMTP_PASS because you don't care, then ignore this error")
				}
			} else {
				log.Fatalf("Missing %s environment variable", k)
			}
		}
	}
	return envs
}

var Envs map[string]string = GetEnvironment()
