package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDb() (*mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI(DbAddress)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	return client.Database("urlShortener").Collection("url"), nil
}

type Url struct {
	ShortUrl string `bson:"_id"`
	LongUrl  string `bson:"longUrl"`
}

func StoreUrl(collection *mongo.Collection, shortUrl string, longUrl string) error {
	_, err := collection.InsertOne(context.TODO(), Url{ShortUrl: shortUrl, LongUrl: longUrl})

	return err
}

func CheckUrl(collection *mongo.Collection, shortUrl string) bool {
	var result Url

	err := collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: shortUrl}}).Decode(&result)

	if err != nil && err == mongo.ErrNoDocuments {
		return false
	}

	return true
}

func GetLongUrl(collection *mongo.Collection, shortUrl string) string {
	var result Url

	err := collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: shortUrl}}).Decode(&result)

	if err != nil && err == mongo.ErrNoDocuments {
		return ""
	}

	return result.LongUrl
}
