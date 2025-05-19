package thesis_repository

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/identities/statuses"
	"go-ourproject/models/response_models"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"time"
)

type ThesisRepository struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client
}

func NewThesisRepository(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *ThesisRepository {
	return &ThesisRepository{
		db:        db,
		logLogrus: logLogrus,
		rdb:       rdb,
	}
}

func (r *ThesisRepository) CreateThesisRepo(thesisReq *identity.ThesisRequest) (response_models.ThesisResponse, int, string, string, error) {
	const op = "repository.ThesisRepository.CreateThesisRepo"

	var thesis = identity.Thesis{
		UserID:      thesisReq.UserID,
		Title:       thesisReq.Title,
		Description: thesisReq.Description,
		StatusID:    3,
	}

	var user identity.Users
	r.db.
		Preload("Role").
		Preload("Major").
		First(&user, thesisReq.UserID)

	user.RoleName = user.Role.Name
	user.MajorName = user.Major.Name

	var status statuses.ThesisStatus
	if err := r.db.Where("id = ?", thesis.StatusID).First(&status).Error; err != nil {
		return response_models.ThesisResponse{}, fiber.StatusNotFound, op, "Status not found", err
	}

	var thesisResp = response_models.ThesisResponse{
		Users:      user,
		Status:     status,
		StatusName: status.Name,
	}

	if err := r.db.Create(&thesis).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return response_models.ThesisResponse{}, fiber.StatusConflict, op, "Thesis repository already exists", err
		}
		return thesisResp, fiber.StatusInternalServerError, op, "Failed to create thesis", err
	}

	return thesisResp, fiber.StatusCreated, op, "Success created thesis", nil
}

func (r *ThesisRepository) CreateSupervisionRepo(supervisionReq *identity.SupervisionRequest) (identity.SupervisionResponse, int, string, string, error) {
	const op = "repository.ThesisRepository.TeacherSelectThesisRepo"

	var supervision = identity.Supervision{
		ThesisID:  supervisionReq.ThesisID,
		TeacherID: supervisionReq.TeacherID,
		Notes:     supervisionReq.Notes,
	}

	var status statuses.ThesisStatus
	if err := r.db.Where("id = ?", supervisionReq.StatusID).First(&status).Error; err != nil {
		return identity.SupervisionResponse{}, fiber.StatusNotFound, op, "Failed to create supervision. status not found", err
	}

	var thesis identity.Thesis
	if err := r.db.Where("id = ?", supervision.ThesisID).First(&thesis).Error; err != nil {
		return identity.SupervisionResponse{}, fiber.StatusNotFound, op, "Failed to create supervision, thesis not found", err
	}

	var teacher identity.Users
	if err := r.db.Where("id = ?", supervisionReq.TeacherID).First(&teacher).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return identity.SupervisionResponse{}, fiber.StatusNotFound, op, "Failed to create supervision, teacher not found", err
	}

	if err := r.db.Create(&supervision).Error; err != nil {
		return identity.SupervisionResponse{}, fiber.StatusInternalServerError, op, "Failed to create supervision", err
	}

	if err := r.db.Model(&thesis).Where("id = ?", thesis.Id).Updates(map[string]any{
		"teacher_id": supervisionReq.TeacherID,
		"status_id":  supervisionReq.StatusID,
	}).Error; err != nil {
		return identity.SupervisionResponse{}, fiber.StatusInternalServerError, op, "Failed to update thesis", err
	}

	var supervisionResp = identity.SupervisionResponse{
		Thesis:      thesis,
		Supervision: supervision,
		Status:      status.Name,
	}

	return supervisionResp, fiber.StatusCreated, op, "Success created supervision", nil
}

func (r *ThesisRepository) GetAllThesisRepo() ([]response_models.GetAllThesisResponse, int, string, string, error) {
	const op = "repository.ThesisRepository.GetAllThesisRepo"

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	thesisJson, err := r.rdb.Get(ctx, "all_thesis").Result()
	if err != nil {
		responseRepo, code, opRepo, msg, err := r.getAllThesisDBMysql()
		return responseRepo, code, opRepo, msg, err
	}

	var thesisResp []response_models.GetAllThesisResponse
	if err := json.Unmarshal([]byte(thesisJson), &thesisResp); err != nil {
		return nil, fiber.StatusInternalServerError, op, "Failed to convert data unmarshal", err
	}

	return thesisResp, fiber.StatusOK, op, "Success get all thesis in redis", nil
}

func (r *ThesisRepository) getAllThesisDBMysql() ([]response_models.GetAllThesisResponse, int, string, string, error) {
	const op = "repository.ThesisRepository.GetAllThesisDBMysql"

	var thesis []identity.Thesis
	if err := r.db.
		Preload("User.Role").
		Preload("User.Major").
		Preload("Status").
		Where("status_id = ?", 3).
		Find(&thesis).Error; err != nil {
		return nil, fiber.StatusInternalServerError, op, "Failed to get thesis", err
	}

	var thesisResp []response_models.GetAllThesisResponse
	for _, th := range thesis {
		thesisResp = append(thesisResp, response_models.GetAllThesisResponse{
			Thesis: th,
			User:   th.User,
			Status: th.Status.Name,
		})
	}

	for i := range thesisResp {
		thesisResp[i].User.RoleName = thesisResp[i].User.Role.Name
		thesisResp[i].User.MajorName = thesisResp[i].User.Major.Name
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	thesisJson, err := json.Marshal(thesisResp)
	if err != nil {
		return nil, fiber.StatusInternalServerError, op, "Failed to convert thesis in data marshal", err
	}

	err = r.rdb.Set(ctx, "all_thesis", thesisJson, 10*time.Second).Err()
	if err != nil {
		return nil, fiber.StatusInternalServerError, op, "Failed to save data thesis in redis", err
	}

	return thesisResp, fiber.StatusOK, op, "Success get all thesis", nil
}

func (r *ThesisRepository) GetThesisRepo(thesisID *int) (response_models.GetAllThesisResponse, int, string, string, error) {
	const op = "repository.ThesisRepository.GetThesis"

	var thesis identity.Thesis
	if err := r.db.
		Preload("User.Role").
		Preload("User.Major").
		Preload("Status").
		First(&thesis, *thesisID).Error; err != nil {
		return response_models.GetAllThesisResponse{}, fiber.StatusInternalServerError, op, "Failed to get thesis", err
	}

	var thesisResponse = response_models.GetAllThesisResponse{
		Thesis: thesis,
		User:   thesis.User,
		Status: thesis.Status.Name,
	}

	thesisResponse.User.RoleName = thesisResponse.User.Role.Name
	thesisResponse.User.MajorName = thesisResponse.User.Major.Name

	return thesisResponse, fiber.StatusOK, op, "Success get thesis", nil
}
