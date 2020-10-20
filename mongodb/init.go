package mongodb

import (
	"context"
	"crawler/util"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionString string = util.Config.MongoDB.ConnectionString
var dbname string = util.Config.MongoDB.DBName
var articleCollection string = util.Config.MongoDB.ArticleCollection
var linkCollection string = util.Config.MongoDB.LinkCollection
var client *mongo.Client

const publisherStr = "publisher"
const idStr = "id"
const linkStr = "link"
const uriString = "uri"

func init() {
	log.SetFlags(log.Lshortfile)
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(util.Config.MongoDB.ConnectionString))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
