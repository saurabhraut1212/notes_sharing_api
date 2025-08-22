package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	DBName    string
	Port      string
	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		MongoURI:  mustEnv("Mongo_URI"),
		DBName:    getEnv("DBName", "dbNotes"),
		Port:      getEnv("Port", "8080"),
		JWTSecret: mustEnv("JWT_SECRET"),
	}
}

func getEnv(k, d string) string {
	v := os.Getenv(k)
	if v != "" {
		return v
	}
	return d

}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing env %s", k)
	}
	return v

}
