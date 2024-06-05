package router

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/nikhil478/function-wale-server/handlers"
	"gorm.io/gorm"
)

// SetupRoutes initializes the routes
func SetupRoutes(app *fiber.App, db *gorm.DB, s3Client *s3.S3) {

	app.Post("/api/signup", handlers.SignUp(db))
	app.Post("/api/login", handlers.Login(db))

	api := app.Group("/api")
	
	api.Get("/organization/:id", handlers.GetOrganizationDetailById(db))
	api.Post("/upload", handlers.UploadFileAmazonS3(s3Client))
	api.Get("/organizations", handlers.GetOrganizations(db))

	// Photo routes
    api.Get("/photos/:id", handlers.GetPhotosByOrganization(db))
    api.Post("/photos", handlers.CreatePhoto(db))

    // Plan routes
    api.Post("/plans", handlers.CreatePlan(db))
    api.Get("/plans", handlers.GetPlansByOrganization(db))
    api.Put("/plans", handlers.UpdatePlan(db))

	api.Use(handlers.JWTMiddleware())
	api.Post("/organization", handlers.CreateOrganization(db))
	api.Put("/organization", handlers.UpdateOrganization(db))
	api.Get("/myorganization", handlers.GetOrganizationByUserID(db))
}