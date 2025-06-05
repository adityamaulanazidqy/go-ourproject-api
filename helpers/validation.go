package helpers

import (
	"errors"
	"go-ourproject/models/auth_models"
	"strings"
)

func ValidateLoginRequest(
	req auth_models.LoginRequest,
	allowedDomain string,
) error {
	if req.Email == "" || req.Password == "" {
		return errors.New("please fill all fields")
	}

	if !strings.HasSuffix(
		req.Email,
		allowedDomain,
	) {
		return errors.New("invalid email domain")
	}

	return nil
}

func ValidateRegisterRequest(
	req auth_models.RegisterRequest,
	allowedDomain string,
) error {
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return errors.New("please fill all fields")
	}

	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	email := strings.ToLower(req.Email)
	domain := strings.ToLower(allowedDomain)

	if !strings.HasSuffix(
		email,
		domain,
	) {
		return errors.New("invalid email domain")
	}

	return nil
}
