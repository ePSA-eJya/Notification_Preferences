package usecase

import (
	"Notification_Preferences/internal/entities"
	"context"

	"github.com/google/uuid"
)

type UserUseCase interface {
	Register(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, email, password string) (string, *entities.User, error)
	FindUserByID(ctx context.Context, id string) (*entities.User, error)
	FindAllUsers(ctx context.Context) ([]*entities.User, error)
	PatchUser(ctx context.Context, id string, user *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, id string) error
	FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]*entities.User, error)
	GetFollowing(ctx context.Context, userID uuid.UUID) ([]*entities.User, error)
}
