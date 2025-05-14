package role_controller

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/jwt_models"
	"go-ourproject/repositories/role_repository"
	"gorm.io/gorm"
	"strconv"
)

type RoleController struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client
	roleRepo  *role_repository.RoleRepository
}

func NewRoleController(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *RoleController {
	return &RoleController{
		db:        db,
		logLogrus: logLogrus,
		rdb:       rdb,
		roleRepo:  role_repository.NewRoleRepository(db, logLogrus, rdb),
	}
}

func (c *RoleController) FindAllRoles(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.StatusMasterpiece"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	responseRepo, code, err := c.roleRepo.FindAllRoleRepository()
	if err != nil {
		c.logError(claims.Email, err, "Failed to get roles")
		return c.handleError(ctx, code, op, err, "Failed to get roles")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *RoleController) FindRoleById(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.FindRolesById"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	roleIDStr := ctx.Params("id", "")
	if roleIDStr == "" {
		err := errors.New("missing roleID")
		c.logError(claims.Email, err, "Failed to get roleID")
		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get roleID")
	}

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.logError(claims.Email, err, "Failed to convert roleID to int")
		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to convert roleID to int")
	}

	responseRepo, code, err := c.roleRepo.FindRolesByIdRepository(roleID)
	if err != nil {
		c.logError(claims.Email, err, "Failed to get roles by id")
		return c.handleError(ctx, code, op, err, "Failed to get roles by id")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *RoleController) CreateRole(ctx *fiber.Ctx) error {
	const op = "masterpiece.controller.CreateRole"

	claims, ok := ctx.Locals("user").(*jwt_models.JWTClaims)
	if !ok || claims == nil {
		err := errors.New("missing claims")
		c.logLogrus.WithFields(logrus.Fields{
			"err":     err,
			"message": "missing claims",
		}).Error("Failed to get claims")

		return c.handleError(ctx, fiber.StatusInternalServerError, op, err, "Failed to get claims")
	}

	var role identity.Roles
	err := ctx.BodyParser(&role)
	if err != nil {
		c.logError(claims.Email, err, "Failed to parse body roles")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Failed to parse body")
	}

	if role.Name == "" {
		c.logError(claims.Email, err, "Failed to parse body roles. role name is required")
		return c.handleError(ctx, fiber.StatusBadRequest, op, err, "Failed to parse body roles. role name is required")
	}

	responseRepo, code, err := c.roleRepo.CreateRoleRepository(&role)
	if err != nil {
		c.logError(claims.Email, err, "Failed to add roles")
		return c.handleError(ctx, code, op, err, "Failed to add roles")
	}

	return ctx.Status(code).JSON(responseRepo)
}

func (c *RoleController) logError(email string, err error, message string) {
	fields := logrus.Fields{"email": email, "message": message}
	if err != nil {
		fields["error"] = err.Error()
	}
	c.logLogrus.WithFields(fields).Error(message)
}

func (c *RoleController) handleError(ctx *fiber.Ctx, status int, op string, err error, message string) error {
	fields := logrus.Fields{
		"operation": op,
		"error":     err,
	}
	c.logLogrus.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(helpers.ApiResponse{
		Message: message,
	})
}
