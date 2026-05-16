package handler

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	Content string `form:"content"`
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

	// Parse form data
	content := c.FormValue("content")
	if content == "" {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "content is required")
	}

	// Create uploads directory with relative path (will be relative to where the app runs from)
	uploadsDir := "uploads/posts"
	if err := os.MkdirAll(uploadsDir, os.ModePerm); err != nil {
		log.Printf("Failed to create uploads directory: %v", err)
		return responses.Error(c, err)
	}

	// Process media files
	mediaURLs := []string{}
	form, err := c.MultipartForm()
	if err != nil && err != io.EOF {
		// If there are no files, that's OK
	}

	if form != nil && form.File["media"] != nil {
		for _, file := range form.File["media"] {
			// Validate file type (only images and videos)
			contentType := file.Header.Get("Content-Type")
			if !isAllowedMediaType(contentType) {
				return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid media type. only images and videos are allowed")
			}

			// Validate file size (max 50MB)
			if file.Size > 50*1024*1024 {
				return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "file too large. max size is 50MB")
			}

			// Generate unique filename
			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("%s_%s%s", userID.String(), uuid.New().String(), ext)
			filePath := filepath.Join(uploadsDir, filename)

			// Save file
			src, err := file.Open()
			if err != nil {
				log.Printf("Failed to open file: %v", err)
				return responses.Error(c, err)
			}
			defer src.Close()

			dst, err := os.Create(filePath)
			if err != nil {
				log.Printf("Failed to create file at %s: %v", filePath, err)
				return responses.Error(c, err)
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				log.Printf("Failed to copy file: %v", err)
				return responses.Error(c, err)
			}
			
			log.Printf("File saved successfully at: %s", filePath)

			// Store relative URL - this URL will be requested from the browser
			mediaURLs = append(mediaURLs, "/uploads/posts/"+filename)
		}
	}

	post := &entities.Post{
		UserID:    userID,
		Content:   content,
		MediaURLs: mediaURLs,
	}

	if err := h.feedUsecase.CreatePost(ctx, post); err != nil {
		return responses.Error(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// isAllowedMediaType checks if the content type is an allowed image or video
func isAllowedMediaType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg":       true,
		"image/jpg":        true,
		"image/png":        true,
		"image/gif":        true,
		"image/webp":       true,
		"video/mp4":        true,
		"video/webm":       true,
		"video/quicktime":  true,
		"video/x-msvideo": true,
	}
	
	// Extract the base content type without parameters
	baseType := strings.Split(contentType, ";")[0]
	return allowedTypes[strings.TrimSpace(baseType)]
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
