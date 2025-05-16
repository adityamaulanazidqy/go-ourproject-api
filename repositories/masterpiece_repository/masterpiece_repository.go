package masterpiece_repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	identity "go-ourproject/models/identities"
	"go-ourproject/models/identities/statuses"
	"go-ourproject/models/response_models"
	"go-ourproject/models/websocket_models"
	"gorm.io/gorm"
	"strings"
	"sync"
	"sync/atomic"
	"time"
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
		file := identity.FileMasterpiece{
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
		Id:        masterpieceEntity.Id,
		Username:  masterpieceEntity.User.Username,
		Email:     masterpieceEntity.User.Email,
		Batch:     masterpieceEntity.User.Batch,
		Photo:     masterpieceEntity.User.Photo,
		CreatedAt: masterpieceEntity.User.CreatedAt,
		UpdatedAt: masterpieceEntity.User.UpdatedAt,
		RoleName:  masterpieceEntity.User.Role.Name,
		MajorName: masterpieceEntity.User.Major.Name,
	}

	var files []string
	for _, file := range filenameResponse {
		files = append(files, file.Name)
	}

	if err := r.rdb.Publish(context.Background(), "masterpiece_updates", "update").Err(); err != nil {
		r.logLogrus.WithFields(logrus.Fields{
			"error":   err,
			"message": "Failed to update masterpiece",
		}).Error("Failed to update masterpiece", err, nil)
	}

	return response_models.MasterpieceResponse{
		User:            userResponse,
		StatusName:      masterpieceEntity.Status.Name,
		ClassName:       masterpieceEntity.Class.Class,
		SemesterName:    masterpieceEntity.Semester.Name,
		LinkGithub:      masterpiece.LinkGithub,
		Files:           files,
		PublicationDate: masterpiece.PublicationDate.Format("2006-01-02"),
	}, fiber.StatusCreated, op, "", nil
}

func (r *MasterpieceRepository) GetMasterpiecesRepository() ([]identity.Masterpiece, int, string, string, error) {
	const op = "repository.Masterpieces.GetMasterpiecesRepository"

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	masterpiecesJson, err := r.rdb.Get(ctx, "masterpieces").Result()
	if err != nil {
		responseRepo, code, opRepo, msg, err := r.getMasterpiecesDBMysql()
		return responseRepo, code, opRepo, msg, err
	}

	var masterpieces []identity.Masterpiece
	err = json.Unmarshal([]byte(masterpiecesJson), &masterpieces)
	if err != nil {
		return masterpieces, fiber.StatusInternalServerError, op, "Failed to convert unmarshal", err
	}

	return masterpieces, fiber.StatusOK, op, "Success Get data masterpieces in redis", nil
}

func (r *MasterpieceRepository) getMasterpiecesDBMysql() ([]identity.Masterpiece, int, string, string, error) {
	const op = "repository.Masterpieces.GetMasterpieces"

	var masterpieces []identity.Masterpiece

	err := r.db.
		Preload("User.Role").
		Preload("User.Major").
		Preload("Status").
		Preload("Class").
		Preload("Semester").
		Preload("Files").
		Preload("Like").
		Preload("Dislike").
		Preload("Comments.User.Role").
		Preload("Comments.User.Major").
		Preload("Comments.Masterpiece").
		Find(&masterpieces).Error

	for i := range masterpieces {
		for j := range masterpieces[i].Comments {
			masterpieces[i].Comments[j].User.RoleName = masterpieces[i].Comments[j].User.Role.Name
			masterpieces[i].Comments[j].User.MajorName = masterpieces[i].Comments[j].User.Major.Name
		}

		var files []string
		for _, file := range masterpieces[i].Files {
			files = append(files, file.FilePath)
		}

		var comments []string
		for _, comment := range masterpieces[i].Comments {
			comments = append(comments, comment.Message)
		}

		masterpieces[i].FilesNames = files
		masterpieces[i].ClassName = masterpieces[i].Class.Class
		masterpieces[i].StatusName = masterpieces[i].Status.Name

		masterpieces[i].User.RoleName = masterpieces[i].User.Role.Name
		masterpieces[i].User.MajorName = masterpieces[i].User.Major.Name

		masterpieces[i].LikeCount = masterpieces[i].Like.Count
		masterpieces[i].DislikesCount = masterpieces[i].Dislike.Count

		masterpieces[i].CommentsArray = comments
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masterpieces, fiber.StatusNotFound, op, "Masterpieces not found", err
		}

		return masterpieces, fiber.StatusInternalServerError, op, "Failed to get masterpieces", err
	}

	if len(masterpieces) == 0 {
		err = errors.New("masterpieces not found")
		return masterpieces, fiber.StatusNotFound, op, "Masterpieces not found", err
	}

	masterpiecesJson, err := json.Marshal(masterpieces)
	if err != nil {
		return masterpieces, fiber.StatusInternalServerError, op, "Failed to convert json marshal", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err = r.rdb.Set(ctx, "masterpieces", masterpiecesJson, 15*time.Second).Err()
	if err != nil {
		return masterpieces, fiber.StatusInternalServerError, op, "Failed to save data in redis", err
	}

	return masterpieces, fiber.StatusOK, op, "Success Get data masterpieces in database", nil
}

func (r *MasterpieceRepository) GetMasterpieceById(masterpieceID string) (identity.Masterpiece, int, string, string, error) {
	const op = "repository.Masterpieces.GetMasterpieceById"

	var masterpiece identity.Masterpiece
	if err := r.db.Where("id = ?", masterpieceID).First(&masterpiece).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return identity.Masterpiece{}, fiber.StatusNotFound, op, "Masterpiece not found", err
		}

		return identity.Masterpiece{}, fiber.StatusInternalServerError, op, "Failed to get data masterpiece", err
	}

	err := r.db.
		Preload("User.Role").
		Preload("User.Major").
		Preload("Status").
		Preload("Class").
		Preload("Semester").
		Preload("Files").
		Preload("Like").
		Preload("Dislike").
		Preload("Comments.User.Role").
		Preload("Comments.User.Major").
		Preload("Comments.Masterpiece").
		Find(&masterpiece).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masterpiece, fiber.StatusNotFound, op, "Masterpiece not found", err
		}

		return masterpiece, fiber.StatusInternalServerError, op, "Failed to get data masterpiece", err
	}

	for i := range masterpiece.Files {
		var files []string
		files = append(files, masterpiece.Files[i].FilePath)
		masterpiece.FilesNames = files
	}

	for j := range masterpiece.Comments {
		masterpiece.Comments[j].User.RoleName = masterpiece.Comments[j].User.Role.Name
		masterpiece.Comments[j].User.MajorName = masterpiece.Comments[j].User.Major.Name
	}

	for c := range masterpiece.Comments {
		var comments []string
		comments = append(comments, masterpiece.Comments[c].Message)
		masterpiece.CommentsArray = comments
	}

	masterpiece.ClassName = masterpiece.Class.Class
	masterpiece.StatusName = masterpiece.Status.Name

	masterpiece.User.RoleName = masterpiece.User.Role.Name
	masterpiece.User.MajorName = masterpiece.User.Major.Name

	masterpiece.LikeCount = masterpiece.Like.Count
	masterpiece.DislikesCount = masterpiece.Dislike.Count

	var masterpiecesEmpety identity.Masterpiece
	err = r.db.Model(&masterpiecesEmpety).Where("id = ?", masterpieceID).UpdateColumn("viewer_count", gorm.Expr("viewer_count + ?", 1)).Error
	if err != nil {
		return masterpiece, fiber.StatusInternalServerError, op, "Failed to update viewer_count", err
	}

	var viewerCount int
	if err = r.db.Table("masterpieces").Where("id = ?", masterpieceID).Select("viewer_count").Scan(&viewerCount).Error; err != nil {
		return masterpiece, fiber.StatusInternalServerError, op, "Failed to get viewer_count", err
	}

	masterpiece.ViewerCount = viewerCount

	return masterpiece, fiber.StatusOK, op, "Success Get data masterpieces", nil
}

func (r *MasterpieceRepository) GetMasterpiecesByStatusId(statusID string) ([]identity.Masterpiece, int, string, string, error) {
	const op = "repository.Masterpieces.GetMasterpiecesByStatusId"

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	masterpiecesJson, err := r.rdb.Get(ctx, fmt.Sprintf("masterpieces_statusID: %s", statusID)).Result()
	if err != nil {
		responseRepo, code, opRepo, msg, err := r.getMasterpiecesByStatusIdDBMysql(statusID)
		return responseRepo, code, opRepo, msg, err
	}

	var masterpieces []identity.Masterpiece
	err = json.Unmarshal([]byte(masterpiecesJson), &masterpieces)
	if err != nil {
		return nil, fiber.StatusInternalServerError, op, "Failed to json unmarshal masterpieces", err
	}

	return masterpieces, fiber.StatusOK, op, fmt.Sprintf("Success Get data masterpieces in redis with status %s", masterpieces[0].StatusName), nil
}

func (r *MasterpieceRepository) getMasterpiecesByStatusIdDBMysql(statusID string) ([]identity.Masterpiece, int, string, string, error) {
	const op = "repository.Masterpieces.GetMasterpiecesByStatusIdDBMysql"

	var masterpieces []identity.Masterpiece

	err := r.db.
		Preload("User.Role").
		Preload("User.Major").
		Preload("Status").
		Preload("Class").
		Preload("Semester").
		Preload("Files").Where("status_id = ?", statusID).
		Preload("Like").
		Preload("Dislike").
		Preload("Comments.User.Role").
		Preload("Comments.User.Major").
		Preload("Comments.Masterpiece").
		Find(&masterpieces).Error

	for i := range masterpieces {
		for j := range masterpieces[i].Comments {
			masterpieces[i].Comments[j].User.RoleName = masterpieces[i].Comments[j].User.Role.Name
			masterpieces[i].Comments[j].User.MajorName = masterpieces[i].Comments[j].User.Major.Name
		}

		var files []string
		for _, file := range masterpieces[i].Files {
			files = append(files, file.FilePath)
		}

		var comments []string
		for _, comment := range masterpieces[i].Comments {
			comments = append(comments, comment.Message)
		}

		masterpieces[i].FilesNames = files
		masterpieces[i].ClassName = masterpieces[i].Class.Class
		masterpieces[i].StatusName = masterpieces[i].Status.Name

		masterpieces[i].User.RoleName = masterpieces[i].User.Role.Name
		masterpieces[i].User.MajorName = masterpieces[i].User.Major.Name

		masterpieces[i].LikeCount = masterpieces[i].Like.Count
		masterpieces[i].DislikesCount = masterpieces[i].Dislike.Count

		masterpieces[i].CommentsArray = comments
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masterpieces, fiber.StatusNotFound, op, "Masterpieces not found", err
		}

		return masterpieces, fiber.StatusInternalServerError, op, "Failed to get data masterpieces", err
	}

	if len(masterpieces) == 0 {
		err = errors.New("no masterpieces found")
		return masterpieces, fiber.StatusNotFound, op, "Masterpieces not found", err
	}

	masterpiecesJson, err := json.Marshal(masterpieces)
	if err != nil {
		return nil, fiber.StatusInternalServerError, op, "Failed to convert json marshal", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err = r.rdb.Set(ctx, fmt.Sprintf("masterpieces_statusID: %s", statusID), masterpiecesJson, 15*time.Second).Err()
	if err != nil {
		return masterpieces, fiber.StatusInternalServerError, op, "Failed to save data in redis", err
	}

	return masterpieces, fiber.StatusOK, op, fmt.Sprintf("Success Get data masterpieces with status %s", masterpieces[0].StatusName), nil
}

var (
	semaphore = make(chan struct{}, 100)
	connCount int32
)

func (r *MasterpieceRepository) SearchMasterpiecesSocket(conn *websocket.Conn) {
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	atomic.AddInt32(&connCount, 1)
	defer atomic.AddInt32(&connCount, -1)

	const op = "repository.Masterpieces.SearchMasterpiecesSocket"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		mu           sync.Mutex
		masterpieces = r.loadInitialData()
		debounce     = 300 * time.Millisecond
		timer        *time.Timer
	)

	go func() {
		pubsub := r.rdb.Subscribe(ctx, "masterpiece_updates")
		defer pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					if !errors.Is(err, context.Canceled) && !errors.Is(err, redis.ErrClosed) {
						r.logLogrus.WithFields(logrus.Fields{
							"operation": op,
							"error":     err,
						}).Warn("PubSub receive error")
					}
					return
				}

				newData := r.loadInitialData()
				mu.Lock()
				masterpieces = newData
				mu.Unlock()
			}
		}
	}()

	conn.SetReadLimit(1024)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				r.logLogrus.WithFields(logrus.Fields{
					"operation": op,
					"error":     err,
				}).Warn("WebSocket unexpected close")
			}
			break
		}

		var searchMsg websocket_models.SearchMessage
		if err := json.Unmarshal(msg, &searchMsg); err != nil {
			r.logLogrus.WithFields(logrus.Fields{
				"operation": op,
				"error":     err,
			}).Warn("Invalid search message")
			continue
		}

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(debounce, func() {
			mu.Lock()
			results := r.searchInMemory(masterpieces, searchMsg.Query)
			mu.Unlock()

			if err := conn.WriteJSON(results); err != nil {
				r.logLogrus.WithFields(logrus.Fields{
					"operation": op,
					"error":     err,
				}).Warn("Failed to send results")
			}
		})
	}
}

func (r *MasterpieceRepository) searchInMemory(masterpieces []identity.Masterpiece, query string) []identity.Masterpiece {
	if query == "" {
		return masterpieces
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var results []identity.Masterpiece

	for _, mp := range masterpieces {
		if r.matchMasterpiece(mp, query) {
			results = append(results, mp)
		}
	}

	return results
}

func (r *MasterpieceRepository) matchMasterpiece(mp identity.Masterpiece, query string) bool {
	fields := []string{
		mp.User.Username,
		mp.User.Email,
		mp.User.Role.Name,
		mp.User.Major.Name,
		mp.StatusName,
		mp.ClassName,
		mp.LinkGithub,
		mp.Semester.Name,
		mp.PublicationDate.Format("2006-01-02"),
		strings.Join(mp.FilesNames, " "),
	}

	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
}

func (r *MasterpieceRepository) loadInitialData() []identity.Masterpiece {
	masterpieces, _, _, _, err := r.GetMasterpiecesRepository()
	if err != nil {
		return []identity.Masterpiece{}
	}
	return masterpieces
}

func (r *MasterpieceRepository) CreateCommentRepository(comments identity.Comments) (identity.Masterpiece, int, string, string, error) {
	const op = "repository.CreateCommentRepository"

	var user identity.Users
	if err := r.db.Where("id = ?", comments.UserId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return identity.Masterpiece{}, fiber.StatusNotFound, op, "User not found", err
		}

		return identity.Masterpiece{}, fiber.StatusInternalServerError, op, "Internal server error", err
	}

	var masterpiece identity.Masterpiece
	if err := r.db.Where("id = ?", comments.MasterpieceID).First(&masterpiece).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return masterpiece, fiber.StatusNotFound, op, "Masterpiece not found", err
		}

		return masterpiece, fiber.StatusInternalServerError, op, "Internal server error", err
	}

	if err := r.db.Create(&comments).Error; err != nil {
		return masterpiece, fiber.StatusInternalServerError, op, "Failed to create comment", err
	}

	err := r.db.
		Preload("User.Role").
		Preload("User.Major").
		Preload("Status").
		Preload("Class").
		Preload("Semester").
		Preload("Files").
		Preload("Like").
		Preload("Dislike").
		Preload("Comments.User.Role").
		Preload("Comments.User.Major").
		Preload("Comments.Masterpiece").
		First(&masterpiece).Error

	if err != nil {
		return masterpiece, fiber.StatusInternalServerError, op, "Failed to create comment", err
	}

	for i := range masterpiece.Files {
		var files []string
		files = append(files, masterpiece.Files[i].FilePath)
		masterpiece.FilesNames = files
	}

	for c := range masterpiece.Comments {
		var comment []string
		comment = append(comment, masterpiece.Comments[c].Message)
		masterpiece.CommentsArray = comment
	}

	for j := range masterpiece.Comments {
		masterpiece.Comments[j].User.RoleName = masterpiece.Comments[j].User.Role.Name
		masterpiece.Comments[j].User.MajorName = masterpiece.Comments[j].User.Major.Name
	}

	masterpiece.ClassName = masterpiece.Class.Class
	masterpiece.StatusName = masterpiece.Status.Name

	masterpiece.User.RoleName = masterpiece.User.Role.Name
	masterpiece.User.MajorName = masterpiece.User.Major.Name

	masterpiece.LikeCount = masterpiece.Like.Count
	masterpiece.DislikesCount = masterpiece.Dislike.Count

	var viewerCount int
	if err = r.db.Table("masterpieces").Where("id = ?", comments.MasterpieceID).Select("viewer_count").Scan(&viewerCount).Error; err != nil {
		return masterpiece, fiber.StatusInternalServerError, op, "Failed to get viewer_count", err
	}

	masterpiece.ViewerCount = viewerCount

	return masterpiece, fiber.StatusCreated, op, "Success create comment", nil
}
