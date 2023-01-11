package queries

import (
	"context"
	"fmt"
	"time"

	b64 "encoding/base64"

	"github.com/TheLoGgI/database"
	"github.com/TheLoGgI/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUser(uid string) (models.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	var userUid = uuid.MustParse(uid)

	var foundUser models.User

	err := collection.FindOne(ctx, bson.D{
		{Key: "uid", Value: userUid},
	}).Decode(&foundUser)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Could not find user with user with uid: %s", uid)
		return foundUser, err
	}

	return foundUser, err
}

func GetUserWithEmail(email string) (models.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	var foundUser models.User
	cursor := collection.FindOne(ctx, bson.D{
		{Key: "email", Value: email},
	})
	cursor.Decode(&foundUser)
	err := cursor.Err()

	// if err != nil {
	// 	fmt.Println("testing for failed user fetch")
	// 	return models.User{}
	// }

	return foundUser, err
}

func CreateSessionCookieForUser(user models.User) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	// GenerateSessionId
	var expires = time.Now().Add(time.Minute*5).UnixNano() / int64(time.Millisecond)
	sessionId := fmt.Sprintf("%s:%d:%s", user.Uid.String(), expires, user.Username)
	fmt.Println("sessionID: " + sessionId)
	encryptedSession := b64.StdEncoding.EncodeToString([]byte(sessionId))

	var updatedDocument bson.M
	filter := bson.D{{Key: "_id", Value: user.EntryId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "sessionId", Value: encryptedSession}}}}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)

	return encryptedSession, err
}

func RemoveSessionCookieForUser(user models.User) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	// GenerateSessionId
	sessionId := fmt.Sprintf("%s:%d:%s", user.Uid, time.Now().Add(time.Minute*5), user.Username)
	encryptedSession := b64.StdEncoding.EncodeToString([]byte(sessionId))

	var updatedDocument bson.M
	filter := bson.D{{Key: "_id", Value: user.EntryId}}
	update := bson.D{{Key: "$unset", Value: bson.D{{Key: "sessionId", Value: ""}}}}
	err := collection.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)

	return encryptedSession, err
}
