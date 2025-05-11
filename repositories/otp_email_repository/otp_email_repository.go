package otp_email_repository

import (
	"errors"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"gorm.io/gorm"
	"net/http"
)

type OtpEmailRepository struct {
	Db        *gorm.DB
	logLogrus *logrus.Logger
}

func NewOtpEmailRepository(db *gorm.DB, logLogrus *logrus.Logger) *OtpEmailRepository {
	return &OtpEmailRepository{
		Db:        db,
		logLogrus: logLogrus,
	}
}

func (repository *OtpEmailRepository) VerificationEmail(email string) (resp helpers.ApiResponse, statusCode int, err error) {
	const (
		msgEmailNotFound      = "Email not found"
		msgVerificationFailed = "Failed to execute verification email"
		msgEmailVerified      = "Email verified"
	)

	var user identity.Users
	result := repository.Db.Table("users").Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			repository.logError(email, result.Error, msgEmailNotFound)
			return helpers.ApiResponse{Message: msgEmailNotFound}, http.StatusOK, nil
		}

		repository.logError(email, result.Error, msgVerificationFailed)
		return helpers.ApiResponse{Message: msgVerificationFailed}, http.StatusInternalServerError, result.Error
	}

	resp = helpers.ApiResponse{
		Message: msgEmailVerified,
	}
	return resp, http.StatusOK, nil
}

func (repository *OtpEmailRepository) logError(email string, err error, message string) {
	repository.logLogrus.WithFields(logrus.Fields{
		"email":   email,
		"error":   err,
		"message": message,
	}).Error(message)
}
