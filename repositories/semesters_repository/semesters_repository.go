package semesters_repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
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

	var semesters []identity.Semesters
	err = json.Unmarshal([]byte(semestersJson), &semesters)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert json unmarshal masterpiece"}, http.StatusInternalServerError, err
	}

	var semestersArray []string
	for _, semester := range semesters {
		semestersArray = append(semestersArray, semester.Name)
	}

	return helpers.ApiResponse{Message: "Success Getting masterpiece status in redis", Data: semestersArray}, http.StatusOK, nil
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

	var semestersArray []string
	for _, semester := range masterpieceStatus {
		semestersArray = append(semestersArray, semester.Name)
	}

	return helpers.ApiResponse{Message: "Successfully get masterpiece Status in database", Data: fiber.Map{
		"semesters": semestersArray,
	}}, http.StatusOK, nil
}

func (r *SemestersRepository) FindSemestersByIdRepository(semesterID int) (helpers.ApiResponse, int, error) {
	var semesters identity.Semesters
	if err := r.db.Where("id = ?", semesterID).First(&semesters).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to find semesters"}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully find semesters by id", Data: fiber.Map{
		"semester": semesters.Name,
	}}, http.StatusOK, nil
}

func (r *SemestersRepository) CreateSemestersRepository(semester *identity.Semesters) (helpers.ApiResponse, int, error) {
	var existingSemester identity.Semesters
	if err := r.db.Where("name = ?", semester.Name).First(&existingSemester).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ApiResponse{Message: "Failed to add semester. name already exists"}, http.StatusConflict, nil
		}

		return helpers.ApiResponse{Message: "Failed to add semester"}, http.StatusInternalServerError, nil
	}

	if err := r.db.Create(&semester).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to add semesters"}, http.StatusInternalServerError, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	semestersJson, err := r.rdb.Get(ctx, "semesters").Result()
	if err != nil {
		return helpers.ApiResponse{}, 0, err
	}

	var semesters []identity.Semesters
	err = json.Unmarshal([]byte(semestersJson), &semesters)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert unmarshal semesters"}, http.StatusInternalServerError, nil
	}

	semesters = append(semesters, *semester)

	updateJson, err := json.Marshal(semesters)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert marshal semesters"}, http.StatusInternalServerError, nil
	}

	err = r.rdb.Set(ctx, "semesters", updateJson, 24*time.Hour).Err()
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to set semesters in redis"}, http.StatusInternalServerError, nil
	}

	return helpers.ApiResponse{Message: "Successfully add semesters", Data: fiber.Map{
		"semester": semester,
	}}, http.StatusOK, nil
}
