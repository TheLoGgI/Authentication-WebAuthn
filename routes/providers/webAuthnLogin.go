package providers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/TheLoGgI/commands"
	"github.com/TheLoGgI/queries"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// var store = session.New()

// // Logging into an account
func BeginLogin(c *fiber.Ctx) error {

	web, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Lasse Aakjær",          // Display Name for your site
		RPID:          "localhost",             // Generally the FQDN for your site
		RPOrigin:      "http://localhost:3001", // The origin URL for WebAuthn requests
		// RPIcon:        "http://localhost/logo.png", // Optional icon URL for your site
	})
	if err != nil {
		fmt.Println(err)
	}

	email := c.FormValue("email")

	// Find user trying to login
	user, userError := queries.GetUserWithEmail(email)

	if userError != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	options, sessionData, loginError := web.BeginLogin(user)

	if loginError != nil {

		c.SendStatus(http.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Trying to login failed",
		})
	}

	commands.UpdateUser(user.Uid, bson.D{
		{Key: "$set", Value: bson.D{{Key: "session", Value: sessionData}}},
	})

	c.SendStatus(http.StatusOK)
	return c.JSON(fiber.Map{
		"options": options,
		"userUid": user.Uid,
		"status":  http.StatusOK,
	})
	// options.publicKey contain our registration options
}

type LoginResponse struct {
	Username  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Uid       uuid.UUID `json:"uid"`
	Exp       time.Time `json:"exp"`
	// SessionId string    `json:"sessionId"`
	// CSRFToken uuid.UUID `json:"csrfToken"`
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {

	var splits = strings.Split(r.RequestURI, "/")
	var userUid = splits[3]

	// get user
	user, getUserError := queries.GetUser(userUid)

	// user doesn't exist
	if getUserError != nil {
		log.Println(getUserError)
		JSONResponse(w, getUserError.Error(), http.StatusBadRequest)
		return
	}

	// load the session data
	sessionData := user.Session

	var _, verificationError = web.FinishLogin(user, sessionData, r)
	if verificationError != nil {
		log.Println(verificationError)
		JSONResponse(w, verificationError.Error(), http.StatusBadRequest)
		return
	}

	// Sent (Refresh-token) and Access-token to client
	privateKey, publicKey := getPublicPrivetKeys()
	const tokenMaxAge = 60 * 60 * 60 /* 60 mins */
	now := time.Now()

	// refreshToken := NewJWT(privateKey, publicKey)
	// var refreshTokenString, err = refreshToken.Create(jwt.MapClaims{
	// 	"sub": user.Uid,                    // Subject of the JWT (the user)
	// 	"iss": "https://lasseaakjaer.com",  // Who created and signed this token
	// 	"exp": now.Add(tokenMaxAge).Unix(), // The expiration time after which the token must be disregarded.
	// 	"iat": now.Unix(),                  // The time at which the token was issued.
	// 	"nbf": now.Unix(),                  // The time before which the token must be disregarded.
	// 	"ver": "1.0",                       // Version
	// 	"typ": "refresh",
	// })

	accessToken := NewJWT(privateKey, publicKey)

	var accessTokenString, accessErr = accessToken.Create(jwt.MapClaims{
		"name":  user.Username,
		"email": user.Email,
		"uid":   user.Uid,
		"sub":   user.Uid,                    // Subject of the JWT (the user)
		"iss":   "https://lasseaakjaer.com",  // Who created and signed this token
		"exp":   now.Add(tokenMaxAge).Unix(), // The expiration time after which the token must be disregarded.
		"iat":   now.Unix(),                  // The time at which the token was issued.
		"nbf":   now.Unix(),                  // The time before which the token must be disregarded.
		"ver":   "1.0",                       // Version
		"typ":   "auth",
	})
	// When validating ->  accessToken.Validate(accessTokenString)

	if accessErr != nil {
		log.Fatalln(err)
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
	}

	// sessionId, sessionIdError := queries.CreateSessionCookieForUser(user)

	// Me endpoint / ID token
	var response = LoginResponse{
		// SessionId: sessionId,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Uid:       user.Uid,
		Exp:       time.Now().Add(time.Hour).UTC(),
		// CSRFToken: uuid.New(),
	}

	// w.Header().Set("Authorization", "Bearer "+tokenString)

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "true",
		Path:     "/",
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh-token",
	// 	Value:    refreshTokenString,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(time.Hour),
	// 	HttpOnly: true,
	// 	Secure:   true,
	// 	SameSite: http.SameSiteLaxMode,
	// })

	http.SetCookie(w, &http.Cookie{
		Name:     "access-token",
		Value:    accessTokenString,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * 30),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// handle successful login
	JSONResponse(w, response, http.StatusOK)
}
