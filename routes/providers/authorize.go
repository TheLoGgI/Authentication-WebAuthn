package providers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

// ClientID
// The client_id is a public identifier for apps. Even though it’s public, it’s best that it isn’t guessable by third parties, so many implementations use something like a 32-character hex string.
// - Foursquare: ZYDPLLBWSK3MVQJSIYHB1OR2JXCY0X2C5UJ2QAR2MAAIT5Q
// - Github: 6779ef20e75817b79602
// - Google: 292085223830.apps.googleusercontent.com
// - Instagram: f2a1ed52710d4533bde25be6da03b6e3
// - SoundCloud: 269d98e4922fb3895e9ae2108cbb5064
// - Windows Live: 00000000400ECB04
// - Okta: 0oa2hl2inow5Uqc6c357

// Client Secret
// The client_secret is a secret known only to the application and the authorization server. It is essential the application’s own password.
// - A great way to generate a secure secret is to use a cryptographically-secure library to generate a 256-bit value and then convert it to a hexadecimal representation.

// type Token struct {
// 	user string
// 	type string
// 	exp int
// 	iat int
// }

// Open ID, ID token /authentication token  -> authentication
// {
// 	"sub": "612f5efd4c1cedb60f46baf1", For the iss and sub fields, specify your service account's email address.
// 	"tid": "2835540a-cd00-420d-9a27-f361ced1e2fa",
//  "aud": For the aud field, specify the API endpoint. For example: https://SERVICE.googleapis.com/.
// 	"type": "refresh",
// 	"iat": 1665183098 For the iat field, specify the current Unix time,
// 	"exp": 1672959099, and for the exp field, specify the time exactly 3600 seconds later, when the JWT will expire.
//   }

func Authorize(w http.ResponseWriter, r *http.Request) {

	// Generate JWT Token

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authenticaton", "application/json")

	// w.Header().Set("Set-Cookie", )

	w.Write([]byte("Hello wold"))
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtPrivateKey = []byte("my_secret_key")

// func TokenString() (string, error) {
// Declare the expiration time of the token
// here, we have kept it as 5 minutes
// expirationTime := time.Now().Add(5 * time.Minute)
// Create the JWT claims, which includes the username and expiry time
// claims := &Claims{
// 	Username: "testing",
// 	StandardClaims: jwt.StandardClaims{
// 		// In JWT, the expiry time is expressed as unix milliseconds
// 		ExpiresAt: expirationTime.Unix(),
// 	},
// }

// token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
// 	"username":  "testing",
// 	"ExpiresAt": expirationTime,
// })
// // Create the JWT string
// // tokenString, err := token.Method.Sign(jwtPrivateKey, "Token")
// // (jwtKey)

// return tokenString, err
// }

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j JWT) Create(claims jwt.Claims) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (j JWT) Validate(token string) (interface{}, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(j.publicKey)
	if err != nil {
		return "", fmt.Errorf("validate: parse key: %w", err)
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims["dat"], nil
}

const TokenMaxAge = 60 * 60 * 10 /* 10 mins */

func getPublicPrivetKeys() ([]byte, []byte) {
	privateKey, privateErr := ioutil.ReadFile("cert/rsa.private")
	if privateErr != nil {
		log.Fatalln(privateErr)
	}
	publicKey, publicErr := ioutil.ReadFile("cert/rsa.public")
	if publicErr != nil {
		log.Fatalln(publicErr)
	}
	return privateKey, publicKey
}

// http://www.inanzzz.com/index.php/post/kdl9/creating-and-validating-a-jwt-rsa-token-in-golang
func Login(w http.ResponseWriter, r *http.Request) {

	// Generate JWT Token
	// var tokenString, _ = TokenString()
	// fmt.Println("Toiken: " + tokenString)

	// CORS Headers
	// w.Header().Set("Access-Control-Allow-Origin", "https://foo.example")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	// w.Header().Set("Access-Control-Allow-Headers", "X-PINGOTHER, Content-Type")
	w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%v", TokenMaxAge))
	w.Header().Set("Content-Type", "application/json")

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)

	// prvKey, privateErr := ioutil.ReadFile("cert/rsa.private")
	// if privateErr != nil {
	// 	log.Fatalln(privateErr)
	// }
	// pubKey, publicErr := ioutil.ReadFile("cert/rsa.public")
	// if publicErr != nil {
	// 	log.Fatalln(publicErr)
	// }

	// Declare the token with the algorithm used for signing, and the claims
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// token := NewJWT(prvKey, pubKey)
	// var tokenString, err = token.Create(TokenMaxAge)
	if err != nil {
		log.Fatalln(err)
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
	}

	// // 2. Validate an existing JWT token.
	// content, err := jwtToken.Validate(tok)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("CONTENT:", content)

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "tokenString",
		Expires:  expirationTime,
		HttpOnly: true,
		// Path:     "/",
		MaxAge: TokenMaxAge,
	})

	w.Write([]byte("Hello wold"))
}
