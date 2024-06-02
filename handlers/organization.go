package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/nikhil478/function-wale-server/models"
	"gorm.io/gorm"
)

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
		organization := models.Organization{}
		if err := ctx.BodyParser(&organization); err != nil {
			ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
			return err
		}

		if err := db.Create(&organization).Error; err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not process request"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "organization has been added", "data": organization})
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
		organization := &models.Organization{}

		if err := ctx.BodyParser(organization); err != nil {
			ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
			return err
		}

		existingOrganization := &models.Organization{}
		if err := db.First(existingOrganization, organization.ID).Error; err != nil {
			ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "organization not found"})
			return err
		}

		if err := db.Model(existingOrganization).Updates(organization).Error; err != nil {
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not update organization details"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "organization details updated successfully"})
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
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get organization details"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "organization details fetched successfully",
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
			ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get organizations"})
			return err
		}

		ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "organizations fetched successfully",
			"data":    organizations,
		})

		return nil
	}
}

