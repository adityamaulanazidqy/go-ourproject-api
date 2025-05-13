package status_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/status_controller"
	"go-ourproject/helpers"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func StatusRoutes(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) {
	controller := status_controller.NewStatusController(db, logLogrus, rdb)

	statusGroup := app.Group("/status")

	statusGroup.Get("/masterpiece", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.StatusMasterpiece(ctx)
	})

	statusGroup.All("/*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(helpers.ApiResponse{
			Message: "Method Not Allowed",
			Data:    nil,
		})
	})
}
