package mongo

import (
	"log"
	"os"
)

const TABLENAME = "packagedb"

type config struct {
	PackageTableName string
	MongoURL         string
}

func newDefaultConfig() *config {
	URL, found := os.LookupEnv("MONGO_URL")
	if !found {
		log.Fatal("MONGO_URL environment variable is not set")
	}
	return &config{
		PackageTableName: TABLENAME,
		MongoURL:         URL,
	}
}

func (c *config) GetPackageTableName() string {
	return c.PackageTableName
}

func (c *config) GetURL() string {
	return c.MongoURL
}
