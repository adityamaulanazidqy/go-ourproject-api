package semesters_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/semesters_controller"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func SemesterRoutes(app *fiber.App, db *gorm.DB, logrusLogger *logrus.Logger, rdb *redis.Client) {
	controller := semesters_controller.NewSemestersController(db, logrusLogger, rdb)

	semesterGroup := app.Group("/semesters")
	semesterGroup.Get("", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.FindAllSemesters(ctx)
	})
	semesterGroup.Post("", middlewares.JWTMiddleware("Guru"), func(ctx *fiber.Ctx) error {
		return controller.CreateSemesters(ctx)
	})
	semesterGroup.Get("/:id", middlewares.JWTMiddleware("Siswa", "Guru", "Pembimbing"), func(ctx *fiber.Ctx) error {
		return controller.FindSemestersById(ctx)
	})
}
