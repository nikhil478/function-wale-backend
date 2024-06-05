package handlers

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nikhil478/function-wale-server/models"
	"gorm.io/gorm"
)

var jwtKey = []byte("your_secret_key")

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	UserID uint   `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// SignUp handles the user registration
// @Summary Sign up a new user
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User Credentials"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/signup [post]
func SignUp(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var creds Credentials
		if err := ctx.BodyParser(&creds); err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request"})
			return err
		}

		if creds.Email == "" || creds.Password == "" {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Email and Password required"})
			return nil
		}

		var existingUser models.User
		if err := db.Where("email = ?", creds.Email).First(&existingUser).Error; err == nil {
			ctx.Status(http.StatusConflict).JSON(&fiber.Map{"message": "User already exists"})
			return nil
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not hash password"})
			return err
		}

		user := models.User{
			Email:    creds.Email,
			Password: string(hashedPassword),
			Role:     "user",
		}

		if err := db.Create(&user).Error; err != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not create user"})
			return err
		}

		expirationTime := time.Now().Add(15 * time.Minute)
		claims := &Claims{
			UserID: user.ID,
			Email:  user.Email,
			Role:   user.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not create token"})
			return err
		}

		ctx.Status(http.StatusCreated).JSON(&fiber.Map{
			"message": "User created successfully",
			"token":   tokenString,
		})
		return nil
	}
}

// Login handles the user login and JWT token issuance
// @Summary Log in a user
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User Credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /api/login [post]
func Login(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var creds Credentials
		if err := ctx.BodyParser(&creds); err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Invalid request"})
			return err
		}

		var user models.User
		if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
			ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
			return err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid email or password"})
			return nil
		}

		expirationTime := time.Now().Add(15 * time.Minute)
		claims := &Claims{
			UserID: user.ID,
			Email:  user.Email,
			Role:   user.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not create token"})
			return err
		}

		ctx.JSON(&fiber.Map{
			"message": "Login successful",
			"token":   tokenString,
		})
		return nil
	}
}

// JWTMiddleware is the middleware for protecting routes
func JWTMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")
		if tokenString == "" {
			ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "Missing or invalid token"})
			return nil
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{"message": "Invalid token"})
			return nil
		}

		ctx.Locals("user", claims)
		return ctx.Next()
	}
}

// CreateOrganization handles the creation of a new organization
// @Summary Create a new organization
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body models.Organization true "Organization Data"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/organization [post]
func CreateOrganization(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("user").(*Claims)

		var count int64
		if err := db.Model(&models.Organization{}).Where("user_id = ?", claims.UserID).Count(&count).Error; err != nil {
			ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "Could not process request"})
			return err
		}

		if count > 0 {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "User already has an organization"})
			return nil
		}

		organization := models.Organization{}
		if err := ctx.BodyParser(&organization); err != nil {
			ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
			return err
		}

		organization.UserID = claims.UserID

		if err := db.Create(&organization).Error; err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not process request"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Organization has been added", "data": organization})
		return nil
	}
}

// UpdateOrganization handles updating an existing organization's details
// @Summary Update organization details
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body models.Organization true "Organization Data"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/organization [put]
func UpdateOrganization(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("user").(*Claims)

		organization := &models.Organization{}
		if err := ctx.BodyParser(organization); err != nil {
			ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed"})
			return err
		}

		existingOrganization := &models.Organization{}
		if err := db.Where("id = ? AND user_id = ?", organization.ID, claims.UserID).First(existingOrganization).Error; err != nil {
			ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "Organization not found"})
			return err
		}

		if err := db.Model(existingOrganization).Updates(organization).Error; err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not update organization details"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "Organization details updated successfully"})
		return nil
	}
}

// GetOrganizationDetailById handles retrieving organization details by ID
// @Summary Get organization details by ID
// @Tags organizations
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/organization/{id} [get]
func GetOrganizationDetailById(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		organization := &models.Organization{}
		id := ctx.Params("id")

		if err := db.Where("id = ?", id).First(organization).Error; err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get organization details"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Organization details fetched successfully",
			"data":    organization,
		})

		return nil
	}
}

// GetOrganizations handles retrieving a list of all organizations
// @Summary Get list of all organizations
// @Tags organizations
// @Produce json
// @Success 200 {array} models.Organization
// @Failure 400 {object} map[string]string
// @Router /api/organizations [get]
func GetOrganizations(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var organizations []models.Organization

		if err := db.Find(&organizations).Error; err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not get organizations"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Organizations fetched successfully",
			"data":    organizations,
		})

		return nil
	}
}

// GetOrganizationByUserID handles retrieving the organization associated with the authenticated user
// @Summary Get organization details by user ID
// @Tags organizations
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/myorganization [get]
func GetOrganizationByUserID(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("user").(*Claims)

		organization := &models.Organization{}
		if err := db.Where("user_id = ?", claims.UserID).First(organization).Error; err != nil {
			ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "No organization found for this user"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Organization details fetched successfully",
			"data":    organization,
		})

		return nil
	}
}
