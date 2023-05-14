package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func SetUrl(path string) string {
	baseUrl := "http://localhost:" + os.Getenv("PORT")
	fmt.Println(baseUrl + path)
	return baseUrl + path
}
