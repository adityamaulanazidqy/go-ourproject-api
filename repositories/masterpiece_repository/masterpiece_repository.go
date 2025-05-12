package masterpiece_repository

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/identities/files"
	"go-ourproject/models/identities/statuses"
	"go-ourproject/models/response_models"
	"gorm.io/gorm"
)

type MasterpieceRepository struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
	rdb       *redis.Client
}

func NewMasterpieceRepository(db *gorm.DB, logLogrus *logrus.Logger, rdb *redis.Client) *MasterpieceRepository {
	return &MasterpieceRepository{
		db:        db,
		logLogrus: logLogrus,
		rdb:       rdb,
	}
}

func (r *MasterpieceRepository) CreateMasterpieceWithFiles(masterpiece *identity.Masterpiece, fileNames []string) (response_models.MasterpieceResponse, int, string, string, error) {
	const op = "masterpiece.repository.CreateMasterpieceWithFiles"

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var (
		user     identity.Users
		semester identity.Semesters
		status   statuses.MasterpieceStatus
		class    identity.Classes
	)

	if err := tx.First(&user, masterpiece.UserID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response_models.MasterpieceResponse{}, fiber.StatusNotFound, op, "User not found", err
		}
		return response_models.MasterpieceResponse{}, fiber.StatusInternalServerError, op, "Failed to verify user", err
	}

	if err := r.db.First(&status, masterpiece.StatusID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response_models.MasterpieceResponse{}, fiber.StatusNotFound, op, "status not found", err
		}

		return response_models.MasterpieceResponse{}, fiber.StatusInternalServerError, op, "Failed to get statuses", err
	}

	if err := r.db.First(&class, masterpiece.ClassID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response_models.MasterpieceResponse{}, fiber.StatusNotFound, op, "class not found", err
		}

		return response_models.MasterpieceResponse{}, fiber.StatusInternalServerError, op, "Failed to get classes", err
	}

	if err := r.db.First(&semester, masterpiece.SemesterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response_models.MasterpieceResponse{}, fiber.StatusNotFound, op, "semester not found", err
		}

		return response_models.MasterpieceResponse{}, fiber.StatusInternalServerError, op, "Failed to get semester", err
	}

	masterpieceEntity := identity.Masterpiece{
		UserID:          masterpiece.UserID,
		StatusID:        masterpiece.StatusID,
		ClassID:         masterpiece.ClassID,
		SemesterID:      masterpiece.SemesterID,
		PublicationDate: masterpiece.PublicationDate,
		LinkGithub:      masterpiece.LinkGithub,
		User:            user,
		Status:          status,
		Class:           class,
		Semester:        semester,
	}

	if err := tx.Create(&masterpieceEntity).Error; err != nil {
		tx.Rollback()
		return response_models.MasterpieceResponse{}, fiber.StatusInternalServerError, op, "Failed to create masterpiece", err
	}

	var filenameResponse []response_models.FileMasterpieceResponse

	for _, fileName := range fileNames {
		file := files.FileMasterpiece{
			MasterpieceID: masterpieceEntity.Id,
			FilePath:      fileName,
		}

		if err := tx.Create(&file).Error; err != nil {
			tx.Rollback()

			r.logLogrus.WithFields(logrus.Fields{
				"masterpiece_id": masterpieceEntity.Id,
				"fileName":       fileName,
				"filePath":       file.FilePath,
			}).Error("Failed to create file", err, nil)

			return response_models.MasterpieceResponse{}, fiber.StatusInternalServerError, op, "Failed to save file records", err
		}

		filenameResponse = append(filenameResponse, response_models.FileMasterpieceResponse{Name: fileName})
	}

	if err := tx.Commit().Error; err != nil {
		return response_models.MasterpieceResponse{}, fiber.StatusInternalServerError, op, "Failed to complete transaction", err
	}

	r.db.
		Preload("User.Role").
		Preload("User.Major").
		Preload("Status").
		Preload("Class").
		Preload("Semester").
		First(&masterpieceEntity, masterpieceEntity.Id)

	var userResponse = identity.UsersResponse{
		Username:  masterpieceEntity.User.Username,
		Email:     masterpieceEntity.User.Email,
		Batch:     masterpieceEntity.User.Batch,
		Photo:     masterpieceEntity.User.Photo,
		CreatedAt: masterpieceEntity.User.CreatedAt,
		UpdatedAt: masterpieceEntity.User.UpdatedAt,
		Role:      masterpieceEntity.User.Role,
		Major:     masterpieceEntity.User.Major,
	}

	return response_models.MasterpieceResponse{
		User:            userResponse,
		Status:          masterpieceEntity.Status,
		Class:           masterpieceEntity.Class,
		Semester:        masterpieceEntity.Semester,
		LinkGithub:      masterpiece.LinkGithub,
		Files:           filenameResponse,
		PublicationDate: masterpiece.PublicationDate,
	}, fiber.StatusCreated, op, "", nil
}
