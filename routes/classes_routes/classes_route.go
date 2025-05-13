package classes_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/classes_controller"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func ClassesRoute(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) {
	controller := classes_controller.NewClassesController(db, logLogrus, rdb)

	classesGroup := app.Group("/classes")
	classesGroup.Get("", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.FindAllClasses(ctx)
	})
	classesGroup.Post("", middlewares.JWTMiddleware("Guru"), func(ctx *fiber.Ctx) error {
		return controller.CreateClasses(ctx)
	})
	classesGroup.Get("/:id", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.FindClassesById(ctx)
	})
}
