package repository

<<<<<<< HEAD
import "Notification_Preferences/internal/entities"
=======
import (
	"context"
	"notification-pref/internal/entities"
>>>>>>> 3f79743 (Add follow/unfollow flow)

	"github.com/google/uuid"
)

// UserRepository is the strict contract that any database implementation must follow
type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, id string) (*entities.User, error)
	FindAll(ctx context.Context) ([]*entities.User, error) // Added ctx here too
	Patch(ctx context.Context, id string, user *entities.User) error
	Delete(ctx context.Context, id string) error
	IsFollowing(ctx context.Context, followerID, followeeID uuid.UUID) (bool, error)
	CreateFollow(ctx context.Context, follow *entities.Follow) error
	DeleteFollow(ctx context.Context, followerID, followeeID uuid.UUID) error
}
