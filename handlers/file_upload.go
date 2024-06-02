package handlers

import (
	"fmt"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nikhil478/function-wale-server/models"
)

// UploadFileAmazonS3 handles file uploads to Amazon S3
// @Summary Upload file to Amazon S3
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} models.ResponseUploadFile
// @Failure 400 {object} map[string]string
// @Router /api/upload [post]
func UploadFileAmazonS3(s3Client *s3.S3) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Accepts("multipart/form-data")

		form, err := ctx.MultipartForm()
		if err != nil {
			fmt.Println("Error parsing multipart form:", err)
			return err
		}

		files := form.File["file"]

		fmt.Println("Number of files received:", len(files))

		if len(files) == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{"message": "no files uploaded"})
		}

		bucketName := "function-wale"
		src, err := files[0].Open()
		if err != nil {
			fmt.Println("Error opening file:", err)
			return err
		}
		defer src.Close()

		s3FileName := uuid.New().String() + filepath.Ext(files[0].Filename)
		s3Path := "uploads/" + s3FileName
		_, err = s3Client.PutObject(&s3.PutObjectInput{
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
}




