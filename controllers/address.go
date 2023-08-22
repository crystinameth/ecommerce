package controllers

import (
	"context"
	"net/http"
	"time"
	"fmt"
	"github.com/crystinameth/ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func AddAddress() gin.HandlerFunc{

	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return
		}

		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500,"Internal Server Error")
		}

		var addresses models.Address

		addresses.Address_id = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		
		match_filter := bson.D{{key:"$match", Value: bson.D{primitive.E{key:"_id", Value: address}}}}
		unwind := bson.D{{key:"$unwind", Value: bson.D{primitive.E{key:"path", Value: "$address"}}}}
		group := bson.D{{key:"$group", Value: bson.D{primitive.E{key:"_id", Value: "$address_id"},{key:"count", Value:bson.D{primitive.E{Key:"$sum", Value: 1}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx,mongo.Pipeline{match_filter, unwind, group})
	
		if err != nil{
			c.IndentedJSON(500,"Internal server error")
		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil{
			panic(err)
		}

		var size int32
		for _, address_no := range addressinfo {
		count := address_no["count"]
		size =count.(int32)
		}
		if size < 2{
			filter := bson.D{primitive.E{key: "_id", Value: address}}
			update := bson.D{{key:"$push", Value: bson.D{primitive.E{key:"address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			}
		} else{
			c.IndentedJSON(400, "Not Allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc{

	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{key:"_id", Value: usert_id}}
		update := bson.D{{key: "$set", Value: bson.D{primitive.E{key:"address.0.house_name", Value: editaddress.House},{key:"address.0.street_name", Value: editaddress.Street},{key:"address.0.city_name", Value: editaddress.City},{key:"address.0.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the home address")
	}	
}

func EditWorkAddress() gin.HandlerFunc{
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}
		var editaddress models.Address
		if err := c.BindJSON(&editaddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{key:"_id", Value: usert_id}}
		update := bson.D{{key: "$set", Value: bson.D{primitive.E{key:"address.1.house_name", Value: editaddress.House},{key:"address.1.street_name", Value: editaddress.Street},{key:"address.1.city_name", Value: editaddress.City},{key:"address.1.pin_code", Value: editaddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the work address")
	}
}

func DeleteAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid Search Index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal server error")
		}

		var ctx, cancel = context.WithTimeOut(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{primitive.E{key:"_id", Value: usert_id}}
		update := bson.D{{key:"$set", Value: bson.D{primitive.E{key:"address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "Wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted")
	}
}
