package database

import (
	"context"
	"github.com/orewaee/go-auth/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func GetCollection(collection string) *mongo.Collection {
	return database.Collection(collection)
}

func Load() error {
	opts := options.Client().ApplyURI(config.MongoUri)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	database = client.Database("go-auth")

	if err := client.Ping(context.Background(), nil); err != nil {
		return err
	}

	return nil
}

func Unload() error {
	return database.Client().Disconnect(context.Background())
}
