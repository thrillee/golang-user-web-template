package apps

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func MountSwaggerAPIDocs(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
