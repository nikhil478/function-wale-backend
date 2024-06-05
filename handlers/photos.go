package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nikhil478/function-wale-server/models"
	"gorm.io/gorm"
)

// CreatePhoto handles the creation of a new photo entry
// @Summary Create a new photo entry
// @Tags photos
// @Accept json
// @Produce json
// @Param photo body models.Photo true "Photo Data"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/photos [post]
func CreatePhoto(db *gorm.DB) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        photo := models.Photo{}
        if err := ctx.BodyParser(&photo); err != nil {
            ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
            return err
        }

        if err := db.Create(&photo).Error; err != nil {
            ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not process request"})
            return err
        }

        ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "photo has been added", "data": photo})
        return nil
    }
}

// GetPhotosByOrganization handles retrieving photos by organization ID and optionally by tag
// @Summary Get photos by organization ID and optionally by tag
// @Tags photos
// @Produce json
// @Param id path string true "Organization ID"
// @Param tag query string false "Photo Tag"
// @Success 200 {array} models.Photo
// @Failure 400 {object} map[string]string
// @Router /api/photos [get]
func GetPhotosByOrganization(db *gorm.DB) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        var photos []models.Photo
        organizationID := ctx.Params("id")
        tag := ctx.Query("tag")

        query := db.Where("organization_id = ?", organizationID)
        if tag != "" {
            query = query.Where("tag = ?", tag)
        }

        if err := query.Find(&photos).Error; err != nil {
            ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get photos"})
            return err
        }

        ctx.Status(http.StatusOK).JSON(&fiber.Map{
            "message": "photos fetched successfully",
            "data":    photos,
        })

        return nil
    }
}


