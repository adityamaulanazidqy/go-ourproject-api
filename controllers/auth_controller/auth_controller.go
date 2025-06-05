package auth_controller

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"go-ourproject/helpers"
	"go-ourproject/models/auth_models"
	identity "go-ourproject/models/identities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type AuthController struct {
	db        *gorm.DB
	logLogrus *logrus.Logger
}

func NewAuthController(
	db *gorm.DB,
	logLogrus *logrus.Logger,
) *AuthController {
	return &AuthController{
		db:        db,
		logLogrus: logLogrus,
	}
}

func (c *AuthController) openFileExcel() (
	*excelize.File,
	error,
) {
	excelFilePath := filepath.Join(
		"assets",
		"data",
		"excel",
		"data_siswa.xlsx",
	)

	if _, err := os.Stat(excelFilePath); os.IsNotExist(err) {
		return nil, errors.New("file does not exist in directory: " + excelFilePath)
	}

	f, err := excelize.OpenFile(excelFilePath)
	if err != nil {
		return nil, fmt.Errorf(
			"gagal membuka file excel: %w",
			err,
		)
	}
	defer f.Close()

	fmt.Printf(
		"Berhasil membuka file Excel: %s\n",
		excelFilePath,
	)

	return f, nil
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	const op = "auth.controller.Login"

	const (
		msgNotFillField  = "Please fill all fields!"
		msgFormatEmail   = "Please use your school email account (@siswa.smktiannajiyah.sch.id)"
		msgEmailNotFound = "Email not found!"
		msgRoleNotFound  = "Role not found!"
		msgServerError   = "Internal server error!"
		msgInvalidPass   = "Invalid password!"
	)

	var req auth_models.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		c.logError(
			req.Email,
			err,
			msgNotFillField,
		)
		return c.handleError(
			ctx,
			fiber.StatusBadRequest,
			op,
			err,
			msgNotFillField,
		)
	}

	err := helpers.ValidateLoginRequest(
		req,
		"@siswa.smktiannajiyah.sch.id",
	)
	if err != nil {
		c.logError(
			req.Email,
			err,
			msgFormatEmail,
		)
		return c.handleError(
			ctx,
			fiber.StatusBadRequest,
			op,
			err,
			msgFormatEmail,
		)
	}

	var user identity.Users
	if err = c.db.Where(&identity.Users{Email: req.Email}).First(&user).Error; err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			c.logError(
				req.Email,
				err,
				msgEmailNotFound,
			)
			return c.handleError(
				ctx,
				fiber.StatusNotFound,
				op,
				err,
				msgEmailNotFound,
			)
		}

		c.logError(
			req.Email,
			err,
			msgServerError,
		)
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			msgServerError,
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		c.logError(
			req.Email,
			err,
			msgInvalidPass,
		)
		return c.handleError(
			ctx,
			fiber.StatusBadRequest,
			op,
			err,
			msgInvalidPass,
		)
	}

	if err = c.db.Preload("Role").Preload("Major").Where(&identity.Users{Email: req.Email}).Find(&user).Error; err != nil {
		if errors.Is(
			err,
			gorm.ErrRecordNotFound,
		) {
			c.logError(
				req.Email,
				err,
				msgRoleNotFound,
			)
			return c.handleError(
				ctx,
				fiber.StatusInternalServerError,
				op,
				err,
				msgRoleNotFound,
			)
		}

		c.logError(
			req.Email,
			err,
			msgServerError,
		)
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			msgServerError,
		)
	}

	token, err := helpers.GenerateToken(
		user.Id,
		user.Email,
		user.Role.Name,
	)
	if err != nil {
		c.logError(
			req.Email,
			err,
			msgServerError,
		)
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			msgServerError,
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		helpers.ApiResponseAuthorization{
			Message: "Login Success",
			Data: auth_models.LoginResponse{
				Username:  user.Username,
				Email:     user.Email,
				RoleName:  user.Role.Name,
				MajorName: user.Major.Name,
			},
			Token: token,
		},
	)
}

//func (c *AuthController) Register(ctx *fiber.Ctx) error {
//	const op = "auth.controller.Register"
//
//	const (
//		msgNotFillField          = "Please fill all fields!"
//		msgFormatEmailAndPass    = "Please use your school email account (@siswa.smktiannajiyah.sch.id) and password must be at least 6 characters"
//		msgRoleNotFound          = "Role not found!"
//		msgMajorNotFound         = "Major not found!"
//		msgEmailAlreadyExists    = "Email Already Exists!"
//		msgUsernameAlreadyExists = "Username Already Exists!"
//		msgHashingPassError      = "Failed to hash password!"
//		msgConvertInt            = "Failed to convert int!"
//		msgProfileImage          = "Failed to get profile image!"
//		msgServerError           = "Internal server error!"
//	)
//
//	var (
//		username = ctx.FormValue("username")
//		email    = ctx.FormValue("email")
//		password = ctx.FormValue("password")
//	)
//
//	roleID, err := strconv.Atoi(ctx.FormValue("role_id"))
//	if err != nil {
//		c.logError(
//			email,
//			err,
//			msgConvertInt,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusInternalServerError,
//			op,
//			err,
//			msgConvertInt,
//		)
//	}
//	majorID, err := strconv.Atoi(ctx.FormValue("major_id"))
//	if err != nil {
//		c.logError(
//			email,
//			err,
//			msgConvertInt,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusInternalServerError,
//			op,
//			err,
//			msgConvertInt,
//		)
//	}
//	batch, err := strconv.Atoi(ctx.FormValue("batch"))
//	if err != nil {
//		c.logError(
//			email,
//			err,
//			msgConvertInt,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusInternalServerError,
//			op,
//			err,
//			msgConvertInt,
//		)
//	}
//
//	if username == "" || email == "" || password == "" || roleID == 0 || majorID == 0 || batch == 0 {
//		c.logError(
//			email,
//			err,
//			msgNotFillField,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusBadRequest,
//			op,
//			nil,
//			msgNotFillField,
//		)
//	}
//
//	var existingUser identity.Users
//	if err := c.db.Where(
//		"username = ?",
//		username,
//	).First(&existingUser).Error; !errors.Is(
//		err,
//		gorm.ErrRecordNotFound,
//	) {
//		c.logError(
//			email,
//			err,
//			msgUsernameAlreadyExists,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusConflict,
//			op,
//			err,
//			msgUsernameAlreadyExists,
//		)
//	}
//
//	if err := c.db.Where(
//		"email = ?",
//		email,
//	).First(&existingUser).Error; !errors.Is(
//		err,
//		gorm.ErrRecordNotFound,
//	) {
//		c.logError(
//			email,
//			err,
//			msgEmailAlreadyExists,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusConflict,
//			op,
//			err,
//			msgEmailAlreadyExists,
//		)
//	}
//
//	var role identity.Roles
//	var major identity.Majors
//
//	if err := c.db.First(
//		&role,
//		roleID,
//	).Error; err != nil {
//		c.logError(
//			email,
//			err,
//			msgRoleNotFound,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusBadRequest,
//			op,
//			err,
//			msgRoleNotFound,
//		)
//	}
//	if err := c.db.First(
//		&major,
//		majorID,
//	).Error; err != nil {
//		c.logError(
//			email,
//			err,
//			msgMajorNotFound,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusBadRequest,
//			op,
//			err,
//			msgMajorNotFound,
//		)
//	}
//
//	file, err := ctx.FormFile("profile")
//	if err != nil {
//		c.logError(
//			email,
//			err,
//			msgProfileImage,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusInternalServerError,
//			op,
//			err,
//			msgProfileImage,
//		)
//	}
//
//	openedFile, err := file.Open()
//	if err != nil {
//		c.logError(
//			email,
//			err,
//			"failed to open uploaded file",
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusInternalServerError,
//			op,
//			err,
//			"failed to open uploaded file",
//		)
//	}
//	defer openedFile.Close()
//
//	var req = auth_models.RegisterRequest{
//		Username: username,
//		Email:    email,
//		Password: password,
//		RoleId:   roleID,
//		MajorId:  majorID,
//		Batch:    batch,
//		Photo:    "",
//	}
//
//	if err = helpers.ValidateRegisterRequest(
//		req,
//		"@siswa.smktiannajiyah.sch.id",
//	); err != nil {
//		c.logError(
//			req.Email,
//			err,
//			msgFormatEmailAndPass,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusBadRequest,
//			op,
//			err,
//			msgFormatEmailAndPass,
//		)
//	}
//
//	hashPass, err := bcrypt.GenerateFromPassword(
//		[]byte(req.Password),
//		bcrypt.DefaultCost,
//	)
//	if err != nil {
//		c.logError(
//			email,
//			err,
//			msgHashingPassError,
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusInternalServerError,
//			op,
//			err,
//			msgHashingPassError,
//		)
//	}
//
//	filename, err := helpers.SaveImages().Profile(
//		openedFile,
//		file,
//		"_",
//	)
//	if err != nil {
//		c.logError(
//			email,
//			err,
//			"failed to save profile image",
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusBadRequest,
//			op,
//			err,
//			"invalid image file",
//		)
//	}
//	req.Photo = filename
//
//	user := identity.Users{
//		Username: req.Username,
//		Email:    req.Email,
//		Password: string(hashPass),
//		RoleID:   uint8(req.RoleId),
//		MajorID:  uint8(req.MajorId),
//		Role:     role,
//		Major:    major,
//		Batch:    req.Batch,
//		Photo:    req.Photo,
//	}
//
//	if err = c.db.Create(&user).Error; err != nil {
//		_ = os.Remove("assets/images/" + filename)
//		c.logError(
//			email,
//			err,
//			"failed to save user profile",
//		)
//		return c.handleError(
//			ctx,
//			fiber.StatusInternalServerError,
//			op,
//			err,
//			msgServerError,
//		)
//	}
//
//	c.db.Preload("Role").Preload("Major").First(
//		&user,
//		user.Id,
//	)
//
//	user.RoleName = role.Name
//	user.MajorName = major.Name
//
//	return ctx.Status(fiber.StatusCreated).JSON(
//		helpers.ApiResponse{
//			Message: "User Created Success",
//			Data:    user,
//		},
//	)
//}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	const op = "auth.controller.Register"

	var req auth_models.RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		c.logError(
			"-",
			err,
			"Failed to parse body",
		)
		return c.handleError(
			ctx,
			fiber.StatusBadRequest,
			op,
			err,
			"Failed to parse body",
		)
	}

	err := helpers.ValidateRegisterRequest(
		req,
		"@siswa.smktiannajiyah.sch.id",
	)
	if err != nil {
		c.logError(
			"-",
			err,
			"Failed to validate email format",
		)

		return c.handleError(
			ctx,
			fiber.StatusBadRequest,
			op,
			err,
			"Failed to validate email format",
		)
	}

	var user identity.Users
	if err := c.db.Where(&identity.Users{Email: req.Email}).First(&user).Error; !errors.Is(
		err,
		gorm.ErrRecordNotFound,
	) {
		c.logError(
			"-",
			err,
			"Email Already exists",
		)
		return c.handleError(
			ctx,
			fiber.StatusConflict,
			op,
			err,
			"Email Already exists",
		)
	}

	f, err := c.openFileExcel()
	if err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			"Failed to open file excel",
		)
	}

	shetname := "student"
	rows, err := f.GetRows(shetname)
	if err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			"Failed to get rows",
		)
	}

	var (
		majorName string
		className string
		batch     string
	)

	for _, row := range rows {
		if len(row) >= 6 && strings.EqualFold(
			row[2],
			req.Email,
		) {
			className = row[3]
			majorName = row[4]
			batch = row[5]
			break
		}
	}

	var class identity.Classes
	if err := c.db.Where(
		"class = ?",
		className,
	).First(&class).Error; err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			fmt.Sprintf(
				"Failed to get class : %s",
				className,
			),
		)
	}

	var major identity.Majors
	if err := c.db.Where(
		"name = ?",
		majorName,
	).First(&major).Error; err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			"Failed to get majors",
		)
	}

	hashPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			"Failed to hash password",
		)
	}

	batchInt, err := strconv.Atoi(batch)
	if err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			"Failed to parse batch",
		)
	}

	var role identity.Roles
	if err := c.db.Where(
		"id = ?",
		1,
	).First(&role).Error; err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			"Failed to get roles",
		)
	}

	user = identity.Users{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashPassword),
		RoleID:   1,
		MajorID:  uint8(major.Id),
		Batch:    batchInt,
	}

	if err := c.db.Create(&user).Error; err != nil {
		return c.handleError(
			ctx,
			fiber.StatusInternalServerError,
			op,
			err,
			"Failed to create user",
		)
	}

	return ctx.Status(fiber.StatusOK).JSON(
		helpers.ApiResponse{
			Message: "Successfully Registered",
			Data: auth_models.RegisterResponse{
				Username: user.Username,
				Email:    user.Email,
				Role:     role.Name,
				Batch:    batchInt,
				Major:    major.Name,
				Photo:    user.Photo,
			},
		},
	)
}

func (c *AuthController) logError(
	email string,
	err error,
	message string,
) {
	fields := logrus.Fields{
		"email":   email,
		"message": message,
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	c.logLogrus.WithFields(fields).Error(message)
}

func (c *AuthController) handleError(
	ctx *fiber.Ctx,
	status int,
	op string,
	err error,
	message string,
) error {
	fields := logrus.Fields{
		"operation": op,
		"error":     err,
	}
	c.logLogrus.WithFields(fields).Error(message)

	return ctx.Status(status).JSON(
		helpers.ApiResponse{
			Message: message,
			Data:    nil,
		},
	)
}
