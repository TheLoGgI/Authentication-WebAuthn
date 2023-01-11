package providers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TheLoGgI/models"
	"github.com/TheLoGgI/queries"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func deleteSessionId(user models.User) {
	time.Sleep(time.Minute * 2)
	fmt.Println("Go routine is awake")

	queries.RemoveSessionCookieForUser(user)
}

func CookieAuthLogin(c *fiber.Ctx) error {
	fmt.Println("Method: " + c.Method())

	// Content Request Headers

	// CORS Headers

	// Preflight Request Headers (CORS)

	fmt.Println(c.ClientHelloInfo())

	// HTTPS is encrypted
	password := c.FormValue("password")
	email := c.FormValue("email")

	// Find user
	user, userError := queries.GetUserWithEmail(email)

	if userError != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	hashedPasswordsDidNotMatch := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if hashedPasswordsDidNotMatch != nil {
		c.SendStatus(http.StatusBadRequest)
		return c.SendString("Password and username din't match\n")
	}

	sessionId, sessionIdError := queries.CreateSessionCookieForUser(user)

	fmt.Println("sessionId: " + sessionId)

	if sessionIdError != nil {
		c.SendStatus(http.StatusBadRequest)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "CookieAuthToken",
		Value:    sessionId,
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		MaxAge:   60 * 30, /* 30 secs */
	})

	c.Cookie(&fiber.Cookie{
		Name:     "userUid",
		Value:    user.Uid.String(),
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: false,
		Secure:   true,
		MaxAge:   60 * 30, /* 30 secs */
	})

	// Deletes sessionId of user after 10 mins, e.g User can no longer validate
	go deleteSessionId(user)

	c.SendStatus(http.StatusOK)

	return c.JSON(fiber.Map{
		"sessionId": sessionId,
	})

}
