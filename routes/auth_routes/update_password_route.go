package auth_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/auth_controller"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func UpdatePasswordRoute(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger) {
	controller := auth_controller.NewUpdatePasswordController(db, logLogrus)

	app.Post("/update_password", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.UpdatePassword(ctx)
	})
}
