package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/TheLoGgI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection URI
// var mongoPassword = os.Getenv("MONGOPASS") // OomqdcOZ5HiNGhlW

var globalClient *mongo.Client

// func InsertOne(ctx *mongo.Collection, bson bson.D) *mongo.InsertOneResult {
// 	// Insert document
// 	// res, err := ctx.InsertOne(context.Background(), bson)
// 	result, err := coll.InsertOne(
// 		context.TODO(),
// 		bson.D{
// 			{"type", "Masala"},
// 			{"rating", 10},
// 			{"vendor", bson.A{"A", "C"}}
// 		}
// 	)

// 	if err != nil {
// 		// return err
// 		fmt.Println(err)
// 	}

// 	return result
// 	// id := res.InsertedID
// 	// fmt.Printf("InsertId: %s", id)
// 	// return id
// }

func PatchUserSignInToken(ctx *mongo.Collection, bson primitive.M) {

}

func GetAuthToken(token string) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var collection = globalClient.Database("salvare").Collection("users")

	filter := bson.D{
		{Key: "token", Value: bson.D{{Key: "$in", Value: bson.A{token}}}},
	}

	result := collection.FindOne(ctx, filter)

	fmt.Println("GetAuthToken Result: ")
	fmt.Println(result)
}

func MongoCollection() *mongo.Collection {
	// ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// Create a new client and connect to the server

	if globalClient == nil {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
		globalClient = client

		if err != nil {
			fmt.Println("Mongo connection failed")
			panic(err)
		}

		// IIFE
		// defer func() {
		// 	if err = client.Disconnect(context.Background()); err != nil {
		// 		panic(err)
		// 	}
		// }()
	}

	return globalClient.Database("salvare").Collection("users")
}

func GetMongoDatabase() *mongo.Database {

	if globalClient == nil {
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("ATLAS_URI")))
		globalClient = client

		if err != nil {
			fmt.Println("Mongo connection failed")
			panic(err)
		}

		// IIFE
		// defer func() {
		// 	if err = client.Disconnect(context.Background()); err != nil {
		// 		panic(err)
		// 	}
		// }()
	}

	return globalClient.Database("salvare")
}

// func HashPassword(password string) (string, error) {
//     bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
//     return string(bytes), err
// }

// func CheckPasswordHash(password, hash string) bool {
//     err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//     return err == nil
// }

func ValidateDatabaseUser(username string, hashedPassword []byte) bool {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	db := globalClient.Database("salvare").Collection("users")

	var foundUser models.User
	err := db.FindOne(ctx, bson.D{
		{Key: "username", Value: username},
	}).Decode(&foundUser)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not find user with username: %s", username)
		return false
	}

	return true
	// var passwordFailed = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), hashedPassword)

	// return passwordFailed == nil

}
