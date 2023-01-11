package models

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	Database     *mongo.Database
	App          *fiber.App
	Http         *http.Server
	SessionStore *session.Store
}
