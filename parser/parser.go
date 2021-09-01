package parser

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

var PackageCounter int

func downloadMd() *os.File {
	b, err := http.Get(URL)
	if err != nil {
		log.Println("unable to get md file from github")
		return nil
	}
	defer b.Body.Close()
	file, err := os.Create(FILE)
	if err != nil {
		log.Println("unable to get md file from github")
		return nil
	}
	io.Copy(file, b.Body)
	return file
}

func SyncReq(newCount int) bool {
	var result olderCount
	collection := MongoClient.Database(Config.UserDBName).Collection(Config.UserDBOldCtr)
	collection.FindOne(context.TODO(), bson.D{}).Decode(&result)
	return newCount > result.Old
}

func Sync() {
	defer MongoClient.Disconnect(context.TODO())
	downloadMd()
	f := FileHandle(FILE)
	m, count := GetSlice(f)
	final := SplitLinks(m)
	check := SyncReq(count)
	if check {
		DBWrite(MongoClient, final)
	}
}

// Open file specified in  filename and return its handle
func FileHandle(filename string) *os.File {
	awsm, err := os.Open(filename)
	if err != nil {
		fmt.Println("cannot read file")
		os.Exit(-1)
	}
	return awsm
}

// GetSlice is a driver function that gets filehandler as an input,
// reads file line-by-line and store slice of raw links into their particular map key
func GetSlice(f *os.File) (map[string][]string, int) {
	rd := bufio.NewReader(f)
	m := make(map[string][]string)
	var links []string
	var title string
	counter := 0
	for {
		line, err := rd.ReadString('\n')
		if strings.HasPrefix(line, "#") || err == io.EOF {
			if links != nil {
				if len(links) < 3 {
					continue
				}
				counter += len(links)
				m[title] = links
				links = nil
				title = line
				if err == io.EOF {
					break
				}
			} else {
				title = line
			}
		} else if strings.HasPrefix(line, "*") {
			links = append(links, line)
		}
	}
	return m, counter
}

// TrimString is a post-processing function that divides an input strings into,
// name -- package name
// url  -- package url
// info  -- a short info about the package
func TrimString(raw string) (string, string, string) {
	sre := regexp.MustCompile(`\[(.*?)\]`)
	rre := regexp.MustCompile(`\((.*?)\)`)
	_name := sre.FindAllString(raw, -1)
	_url := rre.FindAllString(raw, -1)
	if _name == nil || _url == nil {
		return "", "", ""
	}
	name := strings.Trim(_name[0], "[")
	name = strings.Trim(name, "]")
	url := strings.Trim(_url[0], "(")
	url = strings.Trim(url, ")")
	info := strings.Split(raw, "- ")
	if len(info) <= 1 {
		return name, url, ""
	}
	PackageCounter++
	return name, url, info[1]
}

// Split is a driver function for splitting the Line from []Package
// it calls TrimString for splitting and handles a creation and appending of
// a result into a LinkDetails struct.
func SplitLinks(m map[string][]string) Categories {
	categories := make(Categories, len(m))
	i := 0
	for key, value := range m {
		var TmpLinks []Package
		token := strings.IndexByte(key, ' ')
		categories[i].Title = key[token+1:]
		for _, e := range value {
			name, url, info := TrimString(e)
			LD := Package{Name: name, URL: url, Info: info}
			if reflect.ValueOf(LD).IsZero() {
				continue
			}
			TmpLinks = append(TmpLinks, LD)
		}
		categories[i].PackageDetails = append(categories[i].PackageDetails, TmpLinks...)
		i++
	}
	return categories
}
