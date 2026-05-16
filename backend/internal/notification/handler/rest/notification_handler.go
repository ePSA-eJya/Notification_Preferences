package handler

import (
	"fmt"
	"strconv"

	notifRepo "Notification_Preferences/internal/notification/repository"
	"Notification_Preferences/pkg/apperror"
	"Notification_Preferences/pkg/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HttpNotificationHandler struct {
	notifRepo notifRepo.NotificationRepository
}

func NewHttpNotificationHandler(repo notifRepo.NotificationRepository) *HttpNotificationHandler {
	return &HttpNotificationHandler{notifRepo: repo}
}

// POST /api/v1/notifications/read — mark the authenticated user's in-app notifications as read
func (h *HttpNotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.Error(c, apperror.ErrUnauthorized)
	}

	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	if err := h.notifRepo.MarkAllAsRead(ctx, userID); err != nil {
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "notifications marked as read")
}

// GET /api/v1/notifications — fetch the authenticated user's notifications
func (h *HttpNotificationHandler) GetNotifications(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.Error(c, apperror.ErrUnauthorized)
	}

	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	notifications, err := h.notifRepo.GetByRecipientID(ctx, userID, limit)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(fiber.Map{
		"count":         len(notifications),
		"notifications": notifications,
	})
}
