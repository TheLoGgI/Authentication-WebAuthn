package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/TheLoGgI/database"
	"github.com/TheLoGgI/models"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteUser(user models.User) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	var deleteUserFilter = bson.D{
		{Key: "_id", Value: user.Uid},
	}

	cursor := collection.FindOneAndDelete(ctx, deleteUserFilter)

	err := cursor.Err()

	if err != nil {

		fmt.Printf("Delete of User %s, failed", user.Uid)
		panic(err)
	}

}
