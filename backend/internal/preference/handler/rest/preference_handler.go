package handler

import (
	"fmt"

	"Notification_Preferences/internal/entities"
	prefRepo "Notification_Preferences/internal/preference/repository"
	"Notification_Preferences/pkg/apperror"
	"Notification_Preferences/pkg/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HttpPreferenceHandler struct {
	prefRepo prefRepo.PreferenceRepository
}

func NewHttpPreferenceHandler(repo prefRepo.PreferenceRepository) *HttpPreferenceHandler {
	return &HttpPreferenceHandler{prefRepo: repo}
}

// GET /api/v1/preferences — fetch the authenticated user's notification preferences
func (h *HttpPreferenceHandler) GetPreferences(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.Error(c, apperror.ErrUnauthorized)
	}

	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	prefs, err := h.prefRepo.GetPreferenceByUserID(ctx, userID)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(fiber.Map{
		"user_id":     userID,
		"preferences": prefs,
	})
}

// PUT /api/v1/preferences — update the authenticated user's notification preferences
func (h *HttpPreferenceHandler) UpdatePreferences(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.Error(c, apperror.ErrUnauthorized)
	}

	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	var prefs entities.NotificationPreferences
	if err := c.BodyParser(&prefs); err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid preferences body")
	}

	if err := h.prefRepo.UpsertPreference(ctx, userID, prefs); err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(fiber.Map{
		"message":     "preferences updated",
		"user_id":     userID,
		"preferences": prefs,
	})
}
