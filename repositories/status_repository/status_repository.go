package status_repository

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-ourproject/helpers"
	"go-ourproject/models/identities/statuses"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type StatusRepository struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client
}

func NewStatusRepository(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *StatusRepository {
	return &StatusRepository{
		db:        db,
		logLogrus: logLogrus,
		rdb:       rdb,
	}
}

func (r *StatusRepository) StatusMasterpieceRepository() (helpers.ApiResponse, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	masterpieceStatusJson, err := r.rdb.Get(ctx, "masterpiece_status").Result()
	if err != nil {
		responseRepo, code, err := r.statusMasterpieceDBMysql()
		return responseRepo, code, err
	}

	var masterpieceStatus []statuses.MasterpieceStatus
	err = json.Unmarshal([]byte(masterpieceStatusJson), &masterpieceStatus)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert json unmarshal masterpiece"}, http.StatusInternalServerError, err
	}

	return helpers.ApiResponse{Message: "Success Getting masterpiece status in redis", Data: fiber.Map{
		"masterpiece_status": masterpieceStatus,
	}}, http.StatusOK, nil
}

func (r *StatusRepository) statusMasterpieceDBMysql() (helpers.ApiResponse, int, error) {
	var masterpieceStatus []statuses.MasterpieceStatus
	if err := r.db.Find(&masterpieceStatus).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to get masterpiece Status"}, http.StatusInternalServerError, err
	}

	masterpieceStatusJson, err := json.Marshal(masterpieceStatus)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert marshal in masterpiece status"}, http.StatusInternalServerError, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	r.rdb.Set(ctx, "masterpiece_status", masterpieceStatusJson, 24*time.Hour)

	return helpers.ApiResponse{Message: "Successfully get masterpiece Status in database", Data: fiber.Map{
		"masterpiece_status": masterpieceStatus,
	}}, http.StatusOK, nil
}

func (r *StatusRepository) StatusThesisRepository() (helpers.ApiResponse, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	thesisStatusJson, err := r.rdb.Get(ctx, "thesis_status").Result()
	if err != nil {
		responseRepo, code, err := r.statusThesisDBMysql()
		return responseRepo, code, err
	}

	var thesisStatus []statuses.MasterpieceStatus
	err = json.Unmarshal([]byte(thesisStatusJson), &thesisStatus)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert json unmarshal thesis"}, http.StatusInternalServerError, err
	}

	return helpers.ApiResponse{Message: "Success Getting thesis status in redis", Data: fiber.Map{
		"thesis_status": thesisStatus,
	}}, http.StatusOK, nil
}

func (r *StatusRepository) statusThesisDBMysql() (helpers.ApiResponse, int, error) {
	var thesisStatus []statuses.ThesisStatus
	if err := r.db.Find(&thesisStatus).Error; err != nil {
		return helpers.ApiResponse{Message: "Failed to get thesis Status"}, http.StatusInternalServerError, err
	}

	thesisStatusJson, err := json.Marshal(thesisStatus)
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to convert marshal in thesis status"}, http.StatusInternalServerError, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err = r.rdb.Set(ctx, "thesis_status", thesisStatusJson, 24*time.Hour).Err()
	if err != nil {
		return helpers.ApiResponse{Message: "Failed to save data thesis status in redis"}, http.StatusInternalServerError, err
	}

	return helpers.ApiResponse{Message: "Successfully get thesis status in database", Data: fiber.Map{
		"thesis_status": thesisStatus,
	}}, http.StatusOK, nil
}
