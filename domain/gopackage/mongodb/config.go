package mongodb

import (
	"log"
	"os"
)

const TABLENAME = "packagedb"

type Config struct {
	PackageTableName string
	MongoURL         string
}

func NewConfig(tableName, URL string) *Config {
	return &Config{
		PackageTableName: tableName,
		MongoURL:         URL,
	}
}

func WithDefaultConfig() *Config {
	URL, found := os.LookupEnv("ATLAS_URI")
	if !found {
		log.Fatal("MONGO_URL environment variable is not set")
	}
	return &Config{
		PackageTableName: TABLENAME,
		MongoURL:         URL,
	}
}

func (c *Config) GetPackageTableName() string {
	return c.PackageTableName
}

func (c *Config) GetURL() string {
	return c.MongoURL
}
