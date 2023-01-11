package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/TheLoGgI/database"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUser(userUid uuid.UUID, updatedUserModel bson.D) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	var updatedDocument bson.M
	filter := bson.D{{Key: "uid", Value: userUid}}

	err := collection.FindOneAndUpdate(ctx, filter, updatedUserModel).Decode(&updatedDocument)

	// cursor, err := collection.UpdateByID(ctx, userUid.String(), updatedUserModel)

	if err != nil {
		fmt.Println("Update User Failed with userUid: " + userUid.String())
		panic(err)
	}

	fmt.Println(updatedDocument)

}
