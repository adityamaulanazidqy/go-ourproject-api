package middlewares

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	"go-ourproject/models/jwt_models"
	"os"
	"strings"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	rdb       *redis.Client
	logLogrus logrus.Logger
)

func SetRedisClientMiddleware(redisClient *redis.Client) {
	rdb = redisClient
}

func JWTMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ApiResponse{
				Message: "Missing token",
				Data:    nil,
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ApiResponse{
				Message: "Invalid token format",
				Data:    nil,
			})
		}

		tokenStr := parts[1]
		claims := &jwt_models.JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(helpers.ApiResponse{
				Message: "Invalid or expired token",
				Data:    nil,
			})
		}

		if rdb != nil {
			ctxRedis := context.Background()
			blacklisted, err := rdb.Get(ctxRedis, "blacklist:"+tokenStr).Result()
			if err == nil && blacklisted == "true" {
				return c.Status(fiber.StatusUnauthorized).JSON(helpers.ApiResponse{
					Message: "Token has been logged out",
					Data:    nil,
				})
			}
		}

		if len(allowedRoles) > 0 {
			roleMatch := false
			for _, role := range allowedRoles {
				if claims.Roles == role {
					roleMatch = true
					break
				}
			}
			if !roleMatch {
				return c.Status(fiber.StatusForbidden).JSON(helpers.ApiResponse{
					Message: "Forbidden",
					Data:    nil,
				})
			}
		}

		c.Locals("user", claims)
		return c.Next()
	}
}

func ExtractTokenFromHeader(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")
	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1], nil
	}

	return "", errors.New("invalid token format")
}

func VerifyToken(tokenStr string) (*jwt_models.JWTClaims, error) {
	claims := &jwt_models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "invalid token or expired token",
		}).Error("invalid token or expired token")

		return nil, errors.New("invalid token or expired token")
	}

	return claims, nil
}
