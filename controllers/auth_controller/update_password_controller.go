package auth_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/jwt_models"
	"go-ourproject/models/request_models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UpdatePasswordController struct {
	Db        *gorm.DB
	logLogrus *logrus.Logger
}

func NewUpdatePasswordController(db *gorm.DB, logLogrus *logrus.Logger) *UpdatePasswordController {
	return &UpdatePasswordController{
		Db:        db,
		logLogrus: logLogrus,
	}
}

func (controller *UpdatePasswordController) UpdatePassword(c *fiber.Ctx) error {
	var updatePassword request_models.UpdatePasswordRequest

	if err := c.BodyParser(&updatePassword); err != nil {
		controller.logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Invalid request body",
		}).Error("Error in Update Password")

		return c.Status(fiber.StatusBadRequest).JSON(helpers.ApiResponse{
			Message: "Request body not compatible with expected format.",
		})
	}

	claims, ok := c.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims.UserID <= 0 {
		err := errors.New("invalid user ID")
		controller.logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Invalid user ID",
		}).Error("Error in Update Password")

		return c.Status(fiber.StatusBadRequest).JSON(helpers.ApiResponse{
			Message: "Invalid user ID.",
		})
	}

	if updatePassword.Password == "" {
		err := errors.New("invalid user password")
		controller.logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Empty password",
		}).Error("Error in Update Password")

		return c.Status(fiber.StatusBadRequest).JSON(helpers.ApiResponse{
			Message: "Password cannot be empty.",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatePassword.Password), bcrypt.DefaultCost)
	if err != nil {
		controller.logLogrus.WithError(err).Error("Failed to hash password")

		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ApiResponse{
			Message: "Failed to hash password.",
		})
	}

	var user identity.Users
	if err := controller.Db.Model(&user).Where("id = ?", claims.UserID).Update("password", string(hashedPassword)).Error; err != nil {
		controller.logLogrus.WithError(err).Error("Failed to update password")

		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ApiResponse{
			Message: "Failed to update password.",
		})
	}

	controller.logLogrus.WithFields(logrus.Fields{
		"userID": claims.UserID,
	}).Info("Successfully updated password")

	return c.Status(fiber.StatusOK).JSON(helpers.ApiResponse{
		Message: "Password updated successfully.",
	})
}
