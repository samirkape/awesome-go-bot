// types is a part of mparser, responsible for maintaining types, variables and constants
package parser

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	FILE = "./awesome.md"
	URL  = "https://raw.githubusercontent.com/avelino/awesome-go/master/README.md"
)

var (
	DBURI       string
	Config      *config
	MongoClient *mongo.Client
)

func init() {
	// load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Print(".env file found")
	}

	DBURI = os.Getenv("ATLAS_URI")

	Config = &config{
		PackageDBName: "packagedb",
		UserDBName:    "usersdb",
		UserDBOldCtr:  "pkgcount",
		UserDBColName: "requestctr",
		MongoURL:      os.Getenv("ATLAS_URI"),
	}

	MongoClient = GetDbClient()
}

type Categories []Category

type olderCount struct {
	Old int `bson:"old"`
}

type config struct {
	PackageDBName string
	UserDBName    string
	UserDBColName string
	UserDBOldCtr  string
	MongoURL      string
}

type Category struct {
	Title          string
	PackageDetails []Package
	RawLines       []string // * [How To Code in Go eBook](https://www.digitalocean.com/community/books/how-to-code-in-go-ebook) - A 600 page introduction to Go aimed at first time developers.
	SubTitle       string
	Count          int
}

type Package struct {
	Name string `bson:"name" json:"name"`
	URL  string `bson:"url" json:"url"`
	Info string `bson:"info" json:"info"`
	// ID   primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}
