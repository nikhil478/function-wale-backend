package handlers

import (
    "net/http"

    "github.com/gofiber/fiber/v2"
    "github.com/nikhil478/function-wale-server/models"
    "gorm.io/gorm"
)

// CreatePlan handles the creation of a new plan entry
// @Summary Create a new plan entry
// @Tags plans
// @Accept json
// @Produce json
// @Param plan body models.Plan true "Plan Data"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/plans [post]
func CreatePlan(db *gorm.DB) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        plan := models.Plan{}
        if err := ctx.BodyParser(&plan); err != nil {
            ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
            return err
        }

        if err := db.Create(&plan).Error; err != nil {
            ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not process request"})
            return err
        }

        ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "plan has been added", "data": plan})
        return nil
    }
}

// GetPlansByOrganization handles retrieving plans by organization ID and optionally by tag
// @Summary Get plans by organization ID and optionally by tag
// @Tags plans
// @Produce json
// @Param id path string true "Organization ID"
// @Param tag query string false "Plan Tag"
// @Success 200 {array} models.Plan
// @Failure 400 {object} map[string]string
// @Router /api/plans [get]
func GetPlansByOrganization(db *gorm.DB) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        var plans []models.Plan
        organizationID := ctx.Params("id")
        tag := ctx.Query("tag")

        query := db.Where("organization_id = ?", organizationID)
        if tag != "" {
            query = query.Where("tag = ?", tag)
        }

        if err := query.Find(&plans).Error; err != nil {
            ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get plans"})
            return err
        }

        ctx.Status(http.StatusOK).JSON(&fiber.Map{
            "message": "plans fetched successfully",
            "data":    plans,
        })

        return nil
    }
}

// UpdatePlan handles updating an existing plan's details
// @Summary Update plan details
// @Tags plans
// @Accept json
// @Produce json
// @Param plan body models.Plan true "Plan Data"
// @Success 200 {object} map[string]interface{}
// @Failure 422 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/plans [put]
func UpdatePlan(db *gorm.DB) fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        plan := &models.Plan{}

        if err := ctx.BodyParser(plan); err != nil {
            ctx.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
            return err
        }

        existingPlan := &models.Plan{}
        if err := db.First(existingPlan, plan.ID).Error; err != nil {
            ctx.Status(http.StatusNotFound).JSON(&fiber.Map{"message": "plan not found"})
            return err
        }

        if err := db.Model(existingPlan).Updates(plan).Error; err != nil {
            ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not update plan details"})
            return err
        }

        ctx.Status(http.StatusOK).JSON(&fiber.Map{"message": "plan details updated successfully"})
        return nil
    }
}
