package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	AuthRoutes "go-ourproject/routes/auth_routes"
	"go-ourproject/routes/classes_routes"
	"go-ourproject/routes/masterpiece_routes"
	"go-ourproject/routes/otp_email_routes"
	"go-ourproject/routes/role_routes"
	"go-ourproject/routes/semesters_routes"
	"go-ourproject/routes/status_routes"
	"gorm.io/gorm"
)

func Router(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) {
	AuthRoutes.AuthRoute(app, db, logLogrus)
	AuthRoutes.UpdatePasswordRoute(app, db, logLogrus)
	AuthRoutes.LogoutRoute(app, logLogrus, rdb)

	otp_email_routes.OtpEmailRoute(app, db, logLogrus, rdb)

	status_routes.StatusRoutes(app, db, logLogrus, rdb)

	classes_routes.ClassesRoute(app, db, logLogrus, rdb)

	role_routes.RoleRoutes(app, db, logLogrus, rdb)

	masterpiece_routes.MasterpieceRoute(app, db, logLogrus, rdb)

	semesters_routes.SemesterRoutes(app, db, logLogrus, rdb)
}
