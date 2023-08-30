package database

import (
	"context"
	"log"
	"fmt"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

 func DBSet() *mongo.Client{

/* 
	Deprecated
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil{
         log.Fatal(err)
	}
*/


	bsonOpts := &options.BSONOptions{
		UseJSONStructTags: true,
		NilMapAsEmpty:     true,
		NilSliceAsEmpty:   true,
	}

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017").SetBSONOptions(bsonOpts)
    ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	client.Ping(context.TODO(),nil)
	if err != nil {
		log.Println("failed to connect to mongodb :(")
		return nil
	}

	fmt.Println("Successfully connected to mongodb")
	return client

 }

     var Client *mongo.Client = DBSet()
	 
 func UserData(client *mongo.Client, collectionName string) *mongo.Collection{
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return collection
 }

 func ProductData(client *mongo.Client, collectionName string) *mongo.Collection{
	var productcollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return productcollection
 }

