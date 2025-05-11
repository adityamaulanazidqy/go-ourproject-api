package otp_email_routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/controllers/otp_email_controller"
	"go-ourproject/helpers"
	"gorm.io/gorm"
)

func OtpEmailRoute(app *fiber.App, db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) {
	controller := otp_email_controller.NewOtpEmailController(db, logLogrus, rdb)

	otpGroup := app.Group("/otp")

	otpGroup.Post("/send-otp", func(ctx *fiber.Ctx) error {
		return controller.OtpEmail(ctx)
	})

	otpGroup.Post("/verify-otp", func(ctx *fiber.Ctx) error {
		return controller.VerifyOtp(ctx)
	})

	otpGroup.All("/*", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusMethodNotAllowed).JSON(helpers.ApiResponse{
			Message: "Method Not Allowed",
			Data:    nil,
		})
	})
}
