package role_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/role_controller"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func RoleRoutes(app *fiber.App, db *gorm.DB, logrusLogger *logrus.Logger, rdb *redis.Client) {
	controller := role_controller.NewRoleController(db, logrusLogger, rdb)

	rolesGroup := app.Group("/roles")

	rolesGroup.Get("", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.FindAllRoles(ctx)
	})
	rolesGroup.Post("", middlewares.JWTMiddleware("Guru"), func(ctx *fiber.Ctx) error {
		return controller.CreateRole(ctx)
	})
	rolesGroup.Get("/:id", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.FindRoleById(ctx)
	})
}
