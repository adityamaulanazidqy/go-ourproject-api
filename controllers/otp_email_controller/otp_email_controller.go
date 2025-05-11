package otp_email_controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	"go-ourproject/models/request_models"
	"go-ourproject/repositories/otp_email_repository"
	"golang.org/x/net/context"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type OtpEmailController struct {
	Db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client

	otpEmailRepo *otp_email_repository.OtpEmailRepository
}

var (
	smtpUser string
	smtpPass string
)

func NewOtpEmailController(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *OtpEmailController {
	return &OtpEmailController{db, logLogrus, rdb, otp_email_repository.NewOtpEmailRepository(db, logLogrus)}
}

func SetOtpEmail() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	smtpUser = os.Getenv("SMTP_USER")
	smtpPass = os.Getenv("SMTP_PASSWORD")
}

func (controller *OtpEmailController) generateOtpEmail() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // OTP 6 digit
}

func (controller *OtpEmailController) sendEmail(to, otp string) (helpers.ApiResponse, int, error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "adityamaullana234@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Kode OTP - Go-ourprojects")

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<title>Kode OTP Anda</title>
			<style>
				body {
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
					background-color: #f9fafb;
					margin: 0;
					padding: 0;
					color: #333;
				}
				.container {
					max-width: 600px;
					margin: 40px auto;
					background-color: #ffffff;
					padding: 30px;
					border-radius: 10px;
					box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
				}
				.header {
					text-align: center;
					border-bottom: 1px solid #eeeeee;
					padding-bottom: 20px;
					margin-bottom: 20px;
				}
				.header h1 {
					color: #2563eb;
					font-size: 24px;
					margin: 0;
				}
				.content p {
					font-size: 16px;
					line-height: 1.6;
				}
				.otp {
					font-size: 28px;
					font-weight: bold;
					color: #10b981;
					text-align: center;
					margin: 20px 0;
				}
				.footer {
					text-align: center;
					font-size: 12px;
					color: #999999;
					margin-top: 30px;
				}
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>Go-ourprojects</h1>
				</div>
				<div class="content">
					<p>Halo,</p>
					<p>Berikut adalah <strong>Kode OTP</strong> Anda untuk melanjutkan proses verifikasi akun di <strong>Go-ourprojects</strong>:</p>
					<div class="otp">%s</div>
					<p>Jangan berikan kode ini kepada siapa pun. Kode ini hanya berlaku untuk beberapa menit ke depan.</p>
				</div>
				<div class="footer">
					<p>Jika Anda tidak meminta kode ini, Anda bisa mengabaikan email ini.</p>
					<p>&copy; 2025 Go-ourprojects</p>
				</div>
			</div>
		</body>
		</html>
	`, otp)

	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Error:", err)
		fmt.Println("OTP:", otp, "To:", to)
		return helpers.ApiResponse{Message: "Failed to send email", Data: nil}, http.StatusInternalServerError, err
	}

	return helpers.ApiResponse{Message: "Successfully sent email"}, http.StatusOK, nil
}

func (controller *OtpEmailController) OtpEmail(c *fiber.Ctx) error {
	var request request_models.RequestOtpEmail
	if err := c.BodyParser(&request); err != nil {
		controller.logLogrus.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ApiResponse{
			Message: "Failed to parse request body",
			Data:    nil,
		})
	}

	responseRepo, code, err := controller.otpEmailRepo.VerificationEmail(request.Email)
	if err != nil {
		return c.Status(code).JSON(responseRepo)
	}

	otp := controller.generateOtpEmail()
	otpJson, err := json.Marshal(otp)
	if err != nil {
		controller.logLogrus.WithError(err).Error("Failed to marshal otp")
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ApiResponse{
			Message: "Failed to marshal otp",
			Data:    nil,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err = controller.rdb.Set(ctx, fmt.Sprintf("Otp %s:", request.Email), otpJson, 2*time.Minute).Err()
	if err != nil {
		controller.logLogrus.WithError(err).Error("Failed to set otp store")
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ApiResponse{
			Message: "Failed to store OTP",
			Data:    nil,
		})
	}

	responseRepo, code, err = controller.sendEmail(request.Email, otp)
	if err != nil {
		return c.Status(code).JSON(responseRepo)
	}

	return c.Status(code).JSON(responseRepo)
}

func (controller *OtpEmailController) VerifyOtp(c *fiber.Ctx) error {
	var request request_models.VerificationOtpEmail
	if err := c.BodyParser(&request); err != nil {
		controller.logLogrus.WithError(err).Error("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ApiResponse{
			Message: "Failed to parse request body",
			Data:    nil,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	otpJson, err := controller.rdb.Get(ctx, fmt.Sprintf("Otp %s:", request.Email)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			controller.logLogrus.WithError(err).Error("OTP does not exist")
			return c.Status(fiber.StatusNotFound).JSON(helpers.ApiResponse{
				Message: "OTP not found or expired",
				Data:    nil,
			})
		}

		controller.logLogrus.WithError(err).Error("Failed to get otp from Redis")
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ApiResponse{
			Message: "Failed to get OTP",
			Data:    nil,
		})
	}

	var otp string
	if err := json.Unmarshal([]byte(otpJson), &otp); err != nil {
		controller.logLogrus.WithError(err).Error("Failed to unmarshal OTP")
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ApiResponse{
			Message: "Failed to parse stored OTP",
			Data:    nil,
		})
	}

	if otp != request.Otp {
		controller.logLogrus.Warn("Invalid OTP provided")
		return c.Status(fiber.StatusUnauthorized).JSON(helpers.ApiResponse{
			Message: "Invalid OTP",
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ApiResponse{
		Message: "OTP verified successfully",
		Data: fiber.Map{
			"email": request.Email,
			"otp":   otp,
		},
	})
}
