package packages

import (
	"context"
	"github.com/samirkape/awesome-go-bot/domain/gopackage/mongodb"
	"github.com/samirkape/awesome-go-bot/internal/services/packages/analytics/inmemory"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sort"
)

var cache inmemory.AllPackages

type Service interface {
	GetAllPackages() inmemory.AllPackages
}

type Client struct {
	client *mongo.Client
}

func NewService(client *mongo.Client) Service {
	return &Client{client: client}
}

func (c *Client) GetAllPackages() inmemory.AllPackages {
	if cache == nil {
		collections := c.listCollections(mongodb.TABLENAME)
		packages := make(map[inmemory.CategoryName][]inmemory.Package)
		for _, collectionName := range collections {
			p, err := c.getPackagesByCollectionName(mongodb.TABLENAME, collectionName)
			if err != nil {
				return nil
			}
			// sort before adding to map
			sort.Slice(p, func(i, j int) bool {
				return p[i].Stars > p[j].Stars
			})
			packages[inmemory.CategoryName(collectionName)] = p
		}
		cache = packages
	}
	return cache
}

func (c *Client) listCollections(databaseName string) (collections []string) {
	database := c.client.Database(databaseName)

	// List the collections
	collections, err := database.ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	return collections
}

func (c *Client) getPackagesByCollectionName(databaseName string, collectionName string) ([]inmemory.Package, error) {
	cursor := &mongo.Cursor{}
	filter := bson.D{}
	ctx := context.TODO()
	packages := []inmemory.Package{}

	database := c.client.Database(databaseName).Collection(collectionName)

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
		var p inmemory.Package
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
