package handler

import (
	"fmt"
	"strconv"

	"Notification_Preferences/internal/entities"
	"Notification_Preferences/internal/feed/usecase"
	"Notification_Preferences/pkg/apperror"
	"Notification_Preferences/pkg/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HttpFeedHandler struct {
	feedUsecase usecase.FeedUseCase
}

func NewHttpFeedHandler(feedUsecase usecase.FeedUseCase) *HttpFeedHandler {
	return &HttpFeedHandler{feedUsecase: feedUsecase}
}

type createPostRequest struct {
	Content string `json:"content"`
}

func (h *HttpFeedHandler) CreatePost(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}

	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	req := new(createPostRequest)
	if err := c.BodyParser(req); err != nil {
		return responses.Error(c, apperror.ErrInvalidData)
	}

	post := &entities.Post{
		UserID:  userID,
		Content: req.Content,
	}

	if err := h.feedUsecase.CreatePost(ctx, post); err != nil {
		return responses.Error(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

func (h *HttpFeedHandler) LikePost(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}
	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	postIDStr := c.Params("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid post id")
	}

	like := &entities.Like{
		PostID: postID,
		UserID: userID,
	}

	if err := h.feedUsecase.LikePost(ctx, like); err != nil {
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "liked")
}

func (h *HttpFeedHandler) UnlikePost(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}
	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	postIDStr := c.Params("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid post id")
	}

	if err := h.feedUsecase.UnlikePost(ctx, postID, userID); err != nil {
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "unliked")
}

func (h *HttpFeedHandler) IsPostLiked(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}
	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	postIDStr := c.Params("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid post id")
	}

	liked, err := h.feedUsecase.IsPostLikedByUser(ctx, postID, userID)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(fiber.Map{"liked": liked})
}

type commentRequest struct {
	Text string `json:"text"`
}

func (h *HttpFeedHandler) CommentOnPost(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}
	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	postIDStr := c.Params("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid post id")
	}

	req := new(commentRequest)
	if err := c.BodyParser(req); err != nil {
		return responses.Error(c, apperror.ErrInvalidData)
	}

	comment := &entities.Comment{
		PostID: postID,
		UserID: userID,
		Text:   req.Text,
	}

	if err := h.feedUsecase.CommentOnPost(ctx, comment); err != nil {
		return responses.Error(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(comment)
}

func (h *HttpFeedHandler) GetPostComments(c *fiber.Ctx) error {
	ctx := c.UserContext()

	postIDStr := c.Params("id")
	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid post id")
	}

	limit := 50
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	comments, err := h.feedUsecase.GetPostComments(ctx, postID, limit, offset)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(fiber.Map{"comments": comments})
}

func (h *HttpFeedHandler) GetFeed(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}
	userID, err := uuid.Parse(fmt.Sprint(userIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid user id")
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			limit = v
		}
	}
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil {
			offset = v
		}
	}

	posts, err := h.feedUsecase.GetFeed(ctx, userID, limit, offset)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(posts)
}
