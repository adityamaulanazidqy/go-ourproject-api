package auth_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/auth_controller"
	"go-ourproject/helpers"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func AuthRoute(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger) {
	controller := auth_controller.NewAuthController(db, logLogrus)

	authGroup := app.Group("/auth")
	authGroup.Post("/login", middlewares.RateLimiter(), func(ctx *fiber.Ctx) error {
		return controller.Login(ctx)
	})
	authGroup.Post("/register", middlewares.RateLimiter(), func(ctx *fiber.Ctx) error {
		return controller.Register(ctx)
	})

	authGroup.All("/*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(helpers.ApiResponse{
			Message: "Method Not Allowed",
			Data:    nil,
		})
	})
}
