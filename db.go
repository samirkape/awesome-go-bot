// mparser-db is responsible for handling database related operation
// which may include connect, write, query
package mybot

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	CategoryList []string
	Packages     []Package

	// define DB config
	dbconfig struct {
		PackageDBName string
		UserDBName    string
		UserDBColName string
		UserDBKey     string
		MongoURL      string
	}
)

func init() {
	// Initialize database related data structures.
	DBConfig = &dbconfig{
		PackageDBName: "packagedb",
		UserDBName:    "usersdb",
		UserDBColName: "requestctr",
		UserDBKey:     "count",
		MongoURL:      os.Getenv("ATLAS_URI"),
	}

	// Get database handle ( MongoDB ).
	DBClient = InitDbClient()
}

func init() {
	RequestCounter = GetRequestCount()
}

// ListCategories returns a list of categories from database.
// categories are stored as a collections in the database.
func ListCategories() CategoryList {
	c := listCollections(DBClient, DBConfig.PackageDBName)
	return c
}

// PackageByIndex returns a n number of packages stored inside the
// collection as a []Package. index int will be used to map
// name from category slice returned by ListCategories() and
// we get all the documents that belongs to the particular category
// by using a find query with empty bson object.
func PackageByIndex(index int, colls []string) Packages {
	p, _ := findPackages(colls[index])
	return p
}

// listCollections returns database collections as a string slice.
func listCollections(client *mongo.Client, DB string) []string {
	collections, _ := client.Database(DB).ListCollectionNames(context.TODO(), bson.D{})
	return collections
}

// findPackages internally calls Find() on the collection name given by
// colName. The query parameter to Find() is left empty, hence it returns all
// the documents.
func findPackages(colName string) ([]Package, error) {
	// packageList will contains packages that are
	// requested by user by providing category number
	var packageList []Package
	var cur *mongo.Cursor
	var findError error

	// Get database name and client from config
	client := GetDbClient()
	DB := GetPackageDbName()

	// Get collection handle
	collection := client.Database(DB).Collection(colName)

	// bson.D{} specifies 'all documents'
	filter := bson.D{}

	// Find  all documents in the "Collection"
	cur, findError = collection.Find(context.TODO(), filter)

	if findError != nil {
		return nil, findError
	}

	defer cur.Close(context.TODO())

	//Map result to slice
	for cur.Next(context.TODO()) {
		p := Package{}
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		} else {
			packageList = append(packageList, p)
		}
	}
	return packageList, nil
}

// InitDbClient establishes connection to mongodb cloud database for a given URI and
// returns *mongo.Client which needs to be used for further operations on database
// such as finding packages.
func InitDbClient() *mongo.Client {
	mongoURI := GetMongoURI()

	// Get database handle
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
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

// update request counter
func UpdateQueryCount(client *mongo.Client, DbName string, CollectionName string, data interface{}) *mongo.Collection {
	//Create a handle to the respective collection in the database.
	collection := client.Database(DbName).Collection(CollectionName)
	//Perform InsertMany operation & validate against the error.
	_, err := collection.ReplaceOne(context.TODO(), bson.D{}, data)
	if err != nil {
		log.Fatal(err)
	}
	return collection
}

func GetRequestCount() int {
	// init DS
	var result UserRequestCounter
	client := GetDbClient()
	userDBName := GetUserDbName()
	userDBColName := GetUserDbColName()

	// Get user DB collection handle
	collection := client.Database(userDBName).Collection(userDBColName)

	// There is only one document in the user db. FindOne returns that
	err := collection.FindOne(context.TODO(), bson.D{}).Decode(&result)
	if err != nil {
		return -1
	}

	return result.Count
}

func GetMongoURI() string {
	return DBConfig.MongoURL
}

func GetDbClient() *mongo.Client {
	return DBClient
}

func GetPackageDbName() string {
	return DBConfig.PackageDBName
}

func GetUserDbName() string {
	return DBConfig.PackageDBName
}

func GetUserDbColName() string {
	return DBConfig.UserDBColName
}
