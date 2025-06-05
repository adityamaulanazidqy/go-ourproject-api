package main

import (
	"github.com/gofiber/fiber/v2"
	"go-ourproject/config"
	"go-ourproject/controllers/otp_email_controller"
	"go-ourproject/middlewares"
	"go-ourproject/routes"
)

func main() {
	logLogrus := config.LogrusLogger()

	db := config.ConnDB(logLogrus)

	rdb := config.ConnRedis(logLogrus)

	c := fiber.New()
	routes.Router(
		c,
		db,
		logLogrus,
		rdb,
	)

	middlewares.SetRedisClientMiddleware(rdb)

	otp_email_controller.SetOtpEmail()

	logLogrus.Fatal(c.Listen(":8673"))
}
