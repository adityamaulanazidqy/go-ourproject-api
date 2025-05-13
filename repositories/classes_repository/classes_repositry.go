package classes_repository

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

type ClassesRepository struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client
}

func NewClassesRepository(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *ClassesRepository {
	return &ClassesRepository{
		db:        db,
		logLogrus: logLogrus,
		rdb:       rdb,
	}
}

func (r *ClassesRepository) FindAllClassesRepository() (helpers.ApiResponse, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	classesJson, err := r.rdb.Get(ctx, "classes").Result()
	if err != nil {
		responseRepo, code, err := r.findAllClassesDBMysql()
		return responseRepo, code, err
	}

	var classes []identity.Classes
	err = json.Unmarshal([]byte(classesJson), &classes)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert unmarshal classes"}, fiber.StatusInternalServerError, err
	}

	return helpers.ApiResponse{Message: "Successfully find all classes in redis", Data: classes}, http.StatusOK, nil
}

func (r *ClassesRepository) findAllClassesDBMysql() (helpers.ApiResponse, int, error) {
	var classes []identity.Classes
	if err := r.db.Find(&classes).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to find all classes"}, http.StatusInternalServerError, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	classesJson, err := json.Marshal(classes)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert marshal classes"}, http.StatusInternalServerError, nil
	}

	err = r.rdb.Set(ctx, "classes", classesJson, 24*time.Hour).Err()
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to set classes in redis"}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully find all classes in database", Data: classes}, http.StatusOK, nil
}

func (r *ClassesRepository) FindClassesByIdRepository(classID int) (helpers.ApiResponse, int, error) {
	var classes identity.Classes
	if err := r.db.Where("id = ?", classID).First(&classes).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to find classes"}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully find classes by id", Data: classes}, http.StatusOK, nil
}

func (r *ClassesRepository) CreateClassesRepository(classes *identity.Classes) (helpers.ApiResponse, int, error) {
	if err := r.db.Create(&classes).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to add classes", Data: classes}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully add classes", Data: classes}, http.StatusOK, nil
}
