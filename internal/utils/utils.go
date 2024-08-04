package utils

import (
	"os"
)

type Config struct {
	DocDBEndpoint string
	DocDBUser     string
	DocDBPassword string
	DocDBName     string
}

func LoadConfig() *Config {
	return &Config{
		DocDBEndpoint: os.Getenv("DOCDB_ENDPOINT"),
		DocDBUser:     os.Getenv("DOCDB_USER"),
		DocDBPassword: os.Getenv("DOCDB_PASSWORD"),
		DocDBName:     os.Getenv("DOCDB_NAME"),
	}
}
