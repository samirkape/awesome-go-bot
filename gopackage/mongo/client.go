package mongo

import (
	"awesome-go-bot-refactored/gopackage"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Client struct {
	client *mongo.Client
}

func GetClient() *Client {
	config := newDefaultConfig()

	// Get database handle
	clientOptions := options.Client().ApplyURI(config.GetURL())
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{client: client}
}

func (m *Client) GetAllPackages() (map[gopackage.CategoryName][]gopackage.Package, error) {
	collections := m.listCollections(TABLENAME)
	packages := make(map[gopackage.CategoryName][]gopackage.Package)
	for _, collectionName := range collections {
		p, err := m.getPackagesByCollectionName(TABLENAME, collectionName)
		if err != nil {
			return nil, err
		}
		packages[gopackage.CategoryName(collectionName)] = p
	}
	return packages, nil
}

func (m *Client) listCollections(databaseName string) (collections []string) {
	database := m.client.Database(databaseName)

	// List the collections
	collections, err := database.ListCollectionNames(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return collections
}

func (m *Client) getPackagesByCollectionName(databaseName string, collectionName string) ([]gopackage.Package, error) {
	cursor := &mongo.Cursor{}
	filter := bson.D{}
	ctx := context.TODO()
	packages := []gopackage.Package{}

	database := m.client.Database(databaseName).Collection(collectionName)

	cursor, err := database.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {

		err := cur.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var p gopackage.Package
		err := cursor.Decode(&p)
		if err != nil {
			return nil, err
		} else {
			p.Category = collectionName
			packages = append(packages, p)
		}
	}
	return packages, nil
}
