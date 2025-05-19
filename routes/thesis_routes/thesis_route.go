package thesis_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/thesis_controller"
	"go-ourproject/middlewares"
	"gorm.io/gorm"
)

func ThesisRoutes(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) {
	controller := thesis_controller.NewThesisController(db, logLogrus, rdb)

	thesisGroup := app.Group("/thesis")
	thesisGroup.Get("", middlewares.JWTMiddleware("Pembimbing", "Guru"), controller.GetThesis)
	thesisGroup.Get("/all", middlewares.JWTMiddleware("Pembimbing", "Guru"), controller.GetAllThesis)
	thesisGroup.Post("/create", middlewares.JWTMiddleware("Siswa"), controller.CreateThesisTitle)

	supervisionGroup := app.Group("/supervision")
	supervisionGroup.Post("/create", middlewares.JWTMiddleware("Guru", "Pembimbing"), controller.CreateSupervision)
}
