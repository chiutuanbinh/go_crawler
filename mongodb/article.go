package mongodb

import (
	"binhct/common/xtype"
	"log"

	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func Add(id string, article *xtype.Article) error {
	collection := client.Database(dbname).Collection(articleCollection)
	_, err := collection.InsertOne(context.TODO(), article)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func Exist(uri string) bool {
	collection := client.Database(dbname).Collection(articleCollection)
	filter := bson.M{uriString: uri}
	cursor, err := collection.Find(nil, filter)
	if err != nil || !cursor.Next(nil) {
		return false
	}
	return true
}

func GetByID(id string) (*xtype.Article, error) {
	ctx := context.TODO()
	collection := client.Database(dbname).Collection(articleCollection)
	filter := bson.M{idStr: id}
	res := xtype.Article{}
	err := collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}
	return &res, err
}

func Delete(id string) error {
	ctx := context.TODO()
	collection := client.Database(dbname).Collection(articleCollection)
	filter := bson.M{idStr: id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func GetByPublisher(publisher string) ([]xtype.Article, error) {
	ctx := context.TODO()
	collection := client.Database(dbname).Collection(articleCollection)
	filter := bson.M{publisherStr: publisher}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	res := make([]xtype.Article, 0)
	for cursor.Next(ctx) {
		tmp := xtype.Article{}
		err = cursor.Decode(&tmp)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		res = append(res, tmp)
	}
	return res, nil
}
