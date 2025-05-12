package masterpiece_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/masterpiece_controller"
	"go-ourproject/helpers"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func MasterpieceRoute(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) {
	controller := masterpiece_controller.NewMasterpieceController(db, logLogrus, rdb)

	masterpieceGroup := app.Group("/masterpiece")

	masterpieceGroup.Post("/post-masterpiece", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), controller.PostMasterpiece)

	masterpieceGroup.All("/*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(helpers.ApiResponse{
			Message: "Method Not Allowed",
			Data:    nil,
		})
	})
}
