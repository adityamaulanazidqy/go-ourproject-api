package semesters_repository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	identity "go-ourproject/models/identities"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type SemestersRepository struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client
}

func NewSemestersRepository(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *SemestersRepository {
	return &SemestersRepository{
		db:        db,
		logLogrus: logLogrus,
		rdb:       rdb,
	}
}

func (r *SemestersRepository) FindAllSemestersRepository() (helpers.ApiResponse, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	semestersJson, err := r.rdb.Get(ctx, "semesters").Result()
	if err != nil {
		responseRepo, code, err := r.findAllSemestersDBMysql()
		return responseRepo, code, err
	}

	var masterpiece []identity.Semesters
	err = json.Unmarshal([]byte(semestersJson), &masterpiece)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert json unmarshal masterpiece"}, http.StatusInternalServerError, err
	}

	return helpers.ApiResponse{Message: "Success Getting masterpiece status in redis", Data: masterpiece}, http.StatusOK, nil
}

func (r *SemestersRepository) findAllSemestersDBMysql() (helpers.ApiResponse, int, error) {
	var masterpieceStatus []identity.Semesters
	if err := r.db.Find(&masterpieceStatus).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to get masterpiece Status"}, http.StatusInternalServerError, err
	}

	semestersJson, err := json.Marshal(masterpieceStatus)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert marshal in masterpiece status"}, http.StatusInternalServerError, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r.rdb.Set(ctx, "semesters", semestersJson, 24*time.Hour)

	return helpers.ApiResponse{Message: "Successfully get masterpiece Status in database", Data: masterpieceStatus}, http.StatusOK, nil
}

func (r *SemestersRepository) FindSemestersByIdRepository(semesterID int) (helpers.ApiResponse, int, error) {
	var semesters identity.Semesters
	if err := r.db.Where("id = ?", semesterID).First(&semesters).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to find semesters"}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully find semesters by id", Data: semesters}, http.StatusOK, nil
}

func (r *SemestersRepository) CreateSemestersRepository(semesters *identity.Semesters) (helpers.ApiResponse, int, error) {
	if err := r.db.Create(&semesters).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to add semesters", Data: semesters}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully add semesters", Data: semesters}, http.StatusOK, nil
}
