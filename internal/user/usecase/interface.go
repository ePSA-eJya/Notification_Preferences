package usecase

<<<<<<< HEAD
import "Notification_Preferences/internal/entities"
=======
import (
	"context"
	"notification-pref/internal/entities"

	"github.com/google/uuid"
)
>>>>>>> 3f79743 (Add follow/unfollow flow)

type UserUseCase interface {
	Register(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, email, password string) (string, *entities.User, error)
	FindUserByID(ctx context.Context, id string) (*entities.User, error)
	FindAllUsers(ctx context.Context) ([]*entities.User, error)
	PatchUser(ctx context.Context, id string, user *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, id string) error
	FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error
	UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error
}
