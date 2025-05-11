package auth_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/auth_controller"
	"go-ourproject/middlewares"
)

func LogoutRoute(app *fiber.App, logLogrus *logrus.Logger, rdb *redis.Client) {
	controller := auth_controller.NewLogoutController(rdb, logLogrus)

	app.Post("/logout", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.Logout(ctx)
	})
}
