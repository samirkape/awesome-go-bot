package mongo

import (
	"context"
	"log"
	"time"
)

func InitDbClient() *mongo.Client {
	config := newDefaultConfig()

	// Get database handle
	client, err := mongo.NewClient(options.Client().ApplyURI(config.GetURL()))
	if err != nil {
		log.Fatal(err)
	}

	// Set connection deadline to 10 sec
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to database
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
