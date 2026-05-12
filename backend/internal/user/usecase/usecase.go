package usecase

import (
	"context"
	"log"
	"os"
	"time"

	"Notification_Preferences/internal/entities"
	"Notification_Preferences/internal/user/repository"
	"Notification_Preferences/pkg/apperror"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, payload interface{}) error
}

type noopEventPublisher struct{}

func (n *noopEventPublisher) Publish(_ context.Context, _ string, _ interface{}) error {
	return nil
}

// UserService struct
type UserService struct {
	repo           repository.UserRepository
	eventPublisher EventPublisher
}

// Init UserService
func NewUserService(repo repository.UserRepository) UserUseCase {
	return NewUserServiceWithPublisher(repo, nil)
}

func NewUserServiceWithPublisher(repo repository.UserRepository, publisher EventPublisher) UserUseCase {
	if publisher == nil {
		publisher = &noopEventPublisher{}
	}

	return &UserService{repo: repo, eventPublisher: publisher}
}

// Register user (hash password)
func (s *UserService) Register(ctx context.Context, user *entities.User) error {
	existingUser, _ := s.repo.FindByEmail(ctx, user.Email)
	if existingUser != nil {
		return apperror.ErrAlreadyExists
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPwd)
	user.ID = uuid.New() // Nayi unique ID banayega
	return s.repo.Save(ctx, user)
}

// Login user (check email + password)
func (s *UserService) Login(ctx context.Context, email, password string) (string, *entities.User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, err
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *UserService) FindUserByID(ctx context.Context, id string) (*entities.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) FindAllUsers(ctx context.Context) ([]*entities.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) PatchUser(ctx context.Context, id string, user *entities.User) (*entities.User, error) {
	if err := s.repo.Patch(ctx, id, user); err != nil {
		return nil, err
	}

	updatedUser, _ := s.repo.FindByID(ctx, id)
	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) FollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return apperror.ErrInvalidData
	}

	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return err
	}
	if isFollowing {
		return apperror.ErrAlreadyExists
	}

	newFollow := &entities.Follow{
		ID:         uuid.New(),
		FollowerID: followerID,
		FolloweeID: followeeID,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.CreateFollow(ctx, newFollow); err != nil {
		return err
	}

	event := entities.Event{
		ID:         uuid.New(),
		ActorID:    followerID,
		EntityID:   followeeID,
		EntityType: "USER",
		ActionType: entities.Followed,
		CreatedAt:  time.Now().UTC(),
	}
	log.Printf("Publishing FOLLOW event: %+v", event)
	return s.eventPublisher.Publish(ctx, "social_events", event)
}

func (s *UserService) UnfollowUser(ctx context.Context, followerID, followeeID uuid.UUID) error {
	if followerID == followeeID {
		return apperror.ErrInvalidData
	}

	return s.repo.DeleteFollow(ctx, followerID, followeeID)
}

func (s *UserService) GetFollowers(ctx context.Context, userID uuid.UUID) ([]*entities.User, error) {
	ids, err := s.repo.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	followers := make([]*entities.User, 0, len(ids))
	for _, id := range ids {
		user, err := s.repo.FindByID(ctx, id.String())
		if err == nil {
			followers = append(followers, user)
		}
	}
	return followers, nil
}

func (s *UserService) GetFollowing(ctx context.Context, userID uuid.UUID) ([]*entities.User, error) {
	ids, err := s.repo.GetFollowing(ctx, userID)
	if err != nil {
		return nil, err
	}

	following := make([]*entities.User, 0, len(ids))
	for _, id := range ids {
		user, err := s.repo.FindByID(ctx, id.String())
		if err == nil {
			following = append(following, user)
		}
	}
	return following, nil
}
