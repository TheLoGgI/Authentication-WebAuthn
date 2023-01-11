package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TheLoGgI/database"
	"github.com/TheLoGgI/models"
	"go.mongodb.org/mongo-driver/bson"
)

type defaultUser struct {
	id []byte
}

func CreateUser(newUser models.NewUserAccountRequest) models.User {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	// Convert to BSON and insert new user to database
	newUserBson, _ := bson.Marshal(newUser)
	cursor, err := collection.InsertOne(ctx, newUserBson)

	if err != nil {
		panic(err)
	}

	// New user created
	userUid := cursor.InsertedID
	log.Printf("New User created with %s \n", userUid)

	// Find new created user
	var newCreatedUser models.User
	newCreatedErr := collection.FindOne(ctx, bson.D{
		{Key: "uid", Value: newUser.Uid},
	}).Decode(&newCreatedUser)

	if newCreatedErr != nil {
		fmt.Println(newCreatedErr)
		fmt.Printf("Could not find user with username: %s", newUser.Email)
	}

	log.Println(newCreatedUser)

	return newCreatedUser
}
