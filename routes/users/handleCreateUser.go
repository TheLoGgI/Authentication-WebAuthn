package users

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TheLoGgI/commands"
	"github.com/TheLoGgI/database"
	"github.com/TheLoGgI/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type PublicKeyCredentialSource struct {
	Type       []byte // whose value is of PublicKeyCredentialType, defaulting to public-key.
	Id         string //A Credential ID.
	PrivateKey []byte // The credential private key.
	RpId       string // The Relying Party Identifier, for the Relying Party this public key credential source is scoped to.
	UserHandle func() // The user handle associated when this public key credential source was created. This item is nullable.
}

type CreateUserParams struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email" bson:"email"`
}

const emptyUserString = "00000000-0000-0000-0000-000000000000"

func CreateUser(c *fiber.Ctx) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := database.MongoCollection()

	// Check body for password
	formData := new(CreateUserParams)
	c.BodyParser(&formData)

	fmt.Println(formData)

	// Check email registration
	var foundEmailUser models.User
	collection.FindOne(ctx, bson.D{
		{Key: "email", Value: strings.TrimSpace(formData.Email)},
	}).Decode(&foundEmailUser)

	if foundEmailUser.Uid.String() != emptyUserString {
		errMsg := fmt.Sprintf("User with email already exists: %s", foundEmailUser.Username)

		c.SendStatus(http.StatusBadRequest)
		return c.JSON(fiber.Map{
			"error": errMsg,
		})
	}

	// pubKey, publicErr := ioutil.ReadFile("cert/rsa.public")
	// if publicErr != nil {
	// 	log.Fatalln(publicErr)
	// }

	// Create hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(formData.Password), bcrypt.MinCost)
	newUser := models.NewUserAccountRequest{
		Username: formData.Username,
		Email:    formData.Email,
		Password: hashedPassword,
		Uid:      uuid.New(),
		// Credential: webauthn.Credential{
		// 	ID: "",
		// 	PublicKey: pubKey,
		// 	AttestationType: "",
		// 	Authenticator: ""
		// },
	}

	// Create User in database
	commands.CreateUser(newUser)

	c.SendStatus(http.StatusOK)
	return c.JSON(fiber.Map{
		"message":    "User Created",
		"status":     200,
		"statusText": "OK",
	})

}

// encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))

// responseBody := r.Body

// headers := r.Header
// fmt.Println(responseBody)
// fmt.Printf("AuthToken from client: %s created with password %s \n", headers.Get("Auth-Token"), password)
