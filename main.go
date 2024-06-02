package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/nikhil478/function-wale-server/models"
	"github.com/nikhil478/function-wale-server/storage"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB     *gorm.DB
	AWS_S3 *s3.S3
}

func (r *Repository) CreateOrganization(context *fiber.Ctx) error {
	organization := models.Organization{}
	err := context.BodyParser(&organization)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}
	err = r.DB.Create(&organization).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not process request"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "organization has been added"})
	return nil
}

func (r *Repository) UpdateOrganization(ctx *fiber.Ctx) error {
	organization := &models.Organization{}

	if err := ctx.BodyParser(organization); err != nil {
		ctx.Status(fiber.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}
	existingOrganization := &models.Organization{}
	if err := r.DB.First(existingOrganization, organization.ID).Error; err != nil {
		ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{"message": "organization not found"})
		return err
	}

	if err := r.DB.Model(existingOrganization).Updates(organization).Error; err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "could not update organization details"})
		return err
	}

	ctx.Status(fiber.StatusOK).JSON(&fiber.Map{"message": "organization details updated successfully"})
	return nil
}


func (r *Repository) GetOrganizationDetailById(ctx *fiber.Ctx) error {
	organization := &models.Organization{}
	id := ctx.Params("id")

	if err := r.DB.Where("id = ?", id).First(organization).Error; err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "could not get organization details"})
		return err
	}

	ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "organization details fetched successfully",
		"data":    organization,
	})

	return nil
}

func (r *Repository) UploadFileAmazonS3(ctx *fiber.Ctx) error {
	ctx.Accepts("multipart/form-data")

	form, err := ctx.MultipartForm()
	if err != nil {
		fmt.Println("Error parsing multipart form:", err)
		return err
	}

	files := form.File["files"]

	fmt.Println("Number of files received:", len(files))

	bucketName := "function-wale"
	src, err := files[0].Open()
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer src.Close()

	s3FileName := uuid.New().String() + filepath.Ext(files[0].Filename)
	s3Path := "uploads/" + s3FileName
	_, err = r.AWS_S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(s3Path),
		Body:   src,
	})
	if err != nil {
		fmt.Println("Error uploading file to S3:", err)
		return err
	}
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, s3Path)
	url := models.ResponseUploadFile{
		Url: fileURL,
	}
	return ctx.JSON(url)
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/organization", r.CreateOrganization)
	api.Post("/upload", r.UploadFileAmazonS3)
	api.Get("/organization", r.GetOrganizationDetailById)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),     // Added commas after each line
		Port:     os.Getenv("DB_PORT"),     // Added commas after each line
		User:     os.Getenv("DB_USER"),     // Added commas after each line
		Password: os.Getenv("DB_PASSWORD"), // Added commas after each line
		DBName:   os.Getenv("DB_NAME"),     // Added commas after each line
		SSLMode:  os.Getenv("DB_SSLMODE"),  // Added commas after each line
	}

	awsConfig := &storage.AwsS3Config{
		Region: os.Getenv("AWS_S3_REGION"), // Added commas after each line
		Id:     os.Getenv("AWS_S3_ID"),     // Added commas after each line
		Secret: os.Getenv("AWS_S3_SECRET"), // Added commas after each line
		Token:  os.Getenv("AWS_S3_TOKEN"),  // Added commas after each line
	}

	awsS3, err := storage.NewAwsConnection(awsConfig)

	if err != nil {
		log.Fatal("could not load the db")
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the db")
	}

	err = models.MigrateOrganization(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}
	models.MigrateFiles(db)

	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB:     db,
		AWS_S3: awsS3,
	}
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	r.SetupRoutes(app)
	app.Listen(":8000")
}
