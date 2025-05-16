package masterpiece_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

	masterpieceGroup.Post("", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), controller.PostMasterpiece)

	masterpieceGroup.Get("", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), controller.GetMasterpieces)

	masterpieceGroup.Get("/:id", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), controller.GetMasterpieceById)

	masterpieceGroup.Get("/status/:status_id", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), controller.GetMasterpiecesByStatusId)

	masterpieceGroup.Get("/ws/search-masterpieces", websocket.New(func(conn *websocket.Conn) {
		err := controller.SearchMasterpiecesSocket(conn)
		if err != nil {
			return
		}
	}))

	masterpieceGroup.All("/*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(helpers.ApiResponse{
			Message: "Method Not Allowed",
			Data:    nil,
		})
	})
}
