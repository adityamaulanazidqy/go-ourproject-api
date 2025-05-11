package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	AuthRoutes "go-ourproject/routes/auth_routes"
	"go-ourproject/routes/otp_email_routes"
	"gorm.io/gorm"
)

func Router(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) {
	AuthRoutes.AuthRoute(app, db, logLogrus)
	AuthRoutes.UpdatePasswordRoute(app, db, logLogrus)
	AuthRoutes.LogoutRoute(app, logLogrus, rdb)

	otp_email_routes.OtpEmailRoute(app, db, logLogrus, rdb)
}
