package role_repository

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type RoleRepository struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client
}

func NewRoleRepository(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *RoleRepository {
	return &RoleRepository{
		db:        db,
		logLogrus: logLogrus,
		rdb:       rdb,
	}
}

func (r *RoleRepository) FindAllRoleRepository() (helpers.ApiResponse, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	rolesJson, err := r.rdb.Get(ctx, "roles").Result()
	if err != nil {
		responseRepo, code, err := r.findAllRolesDBMysql()
		return responseRepo, code, err
	}

	var roles []identity.Roles
	err = json.Unmarshal([]byte(rolesJson), &roles)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert unmarshal roles"}, fiber.StatusInternalServerError, err
	}

	return helpers.ApiResponse{Message: "Successfully find all roles in redis", Data: roles}, http.StatusOK, nil
}

func (r *RoleRepository) findAllRolesDBMysql() (helpers.ApiResponse, int, error) {
	var roles []identity.Roles
	if err := r.db.Find(&roles).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to find all roles"}, http.StatusInternalServerError, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	rolesJson, err := json.Marshal(roles)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert marshal roles"}, http.StatusInternalServerError, nil
	}

	err = r.rdb.Set(ctx, "roles", rolesJson, 24*time.Hour).Err()
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to set roles in redis"}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully find all roles in database", Data: roles}, http.StatusOK, nil
}

func (r *RoleRepository) FindRolesByIdRepository(roleID int) (helpers.ApiResponse, int, error) {
	var roles identity.Roles
	if err := r.db.Where("id = ?", roleID).First(&roles).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to find roles"}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully find roles by id", Data: roles}, http.StatusOK, nil
}

func (r *RoleRepository) CreateRoleRepository(roles *identity.Roles) (helpers.ApiResponse, int, error) {
	if err := r.db.Create(&roles).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to add roles", Data: roles}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully add roles", Data: roles}, http.StatusOK, nil
}
