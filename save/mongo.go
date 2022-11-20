package save

import (
	"context"
	"jasonLuFa/simpleLine-Webhook/util"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitMongo(config util.Config) (*mongo.Client, context.Context) {
	ctx := context.TODO()
	mongoConn := options.Client().ApplyURI(config.DBSource)
	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		log.Fatal("error while connecting with mongo", err)
	}

	// Know if a MongoDB server has been found and connected to
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("error while trying to ping mongo", err)
	}
	log.Println("mongo connection established")
	return mongoClient, ctx
}
