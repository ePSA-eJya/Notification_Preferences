package handler

import (
	"fmt"

	"Notification_Preferences/internal/entities"
	"Notification_Preferences/internal/user/dto"
	"Notification_Preferences/internal/user/usecase"
	"Notification_Preferences/pkg/apperror"
	"Notification_Preferences/pkg/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type HttpUserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewHttpUserHandler(useCase usecase.UserUseCase) *HttpUserHandler {
	return &HttpUserHandler{userUseCase: useCase}
}

func (h *HttpUserHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	req := new(dto.RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return responses.Error(c, apperror.ErrInvalidData)
	}

	userEntity := dto.ToUserEntity(req)
	if err := h.userUseCase.Register(ctx, userEntity); err != nil {
		return responses.Error(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToUserResponse(userEntity))
}

func (h *HttpUserHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	loginReq := new(dto.LoginRequest)
	if err := c.BodyParser(loginReq); err != nil {
		return responses.Error(c, apperror.ErrInvalidData)
	}

	token, userEntity, err := h.userUseCase.Login(ctx, loginReq.Email, loginReq.Password)
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "invalid email or password")
	}

	return c.JSON(fiber.Map{
		"user":  dto.ToUserResponse(userEntity),
		"token": token,
	})
}

func (h *HttpUserHandler) GetUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userID := c.Locals("user_id")
	if userID == nil {
		return responses.Error(c, apperror.ErrInvalidData)
	}

	userEntity, err := h.userUseCase.FindUserByID(ctx, fmt.Sprint(userID))
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponse(userEntity))
}

func (h *HttpUserHandler) FindUserByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")
	if id == "" {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "id is required")
	}

	userEntity, err := h.userUseCase.FindUserByID(ctx, id)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponse(userEntity))
}

func (h *HttpUserHandler) FindAllUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()

	users, err := h.userUseCase.FindAllUsers(ctx)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponseList(users))
}

func (h *HttpUserHandler) PatchUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")

	var req dto.PatchUserRequest
	if err := c.BodyParser(&req); err != nil {
		return responses.ErrorWithMessage(c, err, "invalid request")
	}

	user := &entities.User{UserHandle: req.UserHandle}

	msg, err := validatePatchUser(user)
	if err != nil {
		return responses.ErrorWithMessage(c, err, msg)
	}

	updatedUser, err := h.userUseCase.PatchUser(ctx, id, user)
	if err != nil {
		return responses.Error(c, err)
	}

	return c.JSON(dto.ToUserResponse(updatedUser))
}

func (h *HttpUserHandler) DeleteUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id := c.Params("id")

	if err := h.userUseCase.DeleteUser(ctx, id); err != nil {
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "user deleted")
}

func (h *HttpUserHandler) FollowUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	followerIDRaw := c.Locals("user_id")
	if followerIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}

	followerID, err := uuid.Parse(fmt.Sprint(followerIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid follower id")
	}

	followeeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid followee id")
	}

	if err := h.userUseCase.FollowUser(ctx, followerID, followeeID); err != nil {
		if err == apperror.ErrInvalidData {
			return responses.ErrorWithMessage(c, err, "you cannot follow yourself")
		}
		if err == apperror.ErrAlreadyExists {
			return responses.ErrorWithMessage(c, err, "you are already following this user")
		}
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "successfully followed user")
}

func (h *HttpUserHandler) UnfollowUser(c *fiber.Ctx) error {
	ctx := c.UserContext()

	followerIDRaw := c.Locals("user_id")
	if followerIDRaw == nil {
		return responses.ErrorWithMessage(c, apperror.ErrUnauthorized, "missing authenticated user")
	}

	followerID, err := uuid.Parse(fmt.Sprint(followerIDRaw))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid follower id")
	}

	followeeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "invalid followee id")
	}

	if err := h.userUseCase.UnfollowUser(ctx, followerID, followeeID); err != nil {
		if err.Error() == "follow not found" {
			return responses.ErrorWithMessage(c, apperror.ErrInvalidData, "follow relationship not found")
		}
		return responses.Error(c, err)
	}

	return responses.Message(c, fiber.StatusOK, "successfully unfollowed user")
}

func validatePatchUser(user *entities.User) (string, error) {
	if user.UserHandle == "" {
		return "username is invalid", apperror.ErrInvalidData
	}
	return "", nil
}
