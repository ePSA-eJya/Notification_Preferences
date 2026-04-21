package usecase_test

import (
	"context"
	"errors"
	"os"
	"sort"
	"sync"
	"testing"

	"notification-pref/internal/entities"
	"notification-pref/internal/user/usecase"
	"notification-pref/pkg/apperror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type inMemoryUserRepository struct {
	mu      sync.RWMutex
	users   map[uuid.UUID]*entities.User
	follows map[string]*entities.Follow
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{
		users:   map[uuid.UUID]*entities.User{},
		follows: map[string]*entities.Follow{},
	}
}

func followKey(followerID, followeeID uuid.UUID) string {
	return followerID.String() + ":" + followeeID.String()
}

func cloneUser(user *entities.User) *entities.User {
	if user == nil {
		return nil
	}
	copy := *user
	return &copy
}

func (r *inMemoryUserRepository) Save(_ context.Context, user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	r.users[user.ID] = cloneUser(user)
	return nil
}

func (r *inMemoryUserRepository) FindByEmail(_ context.Context, email string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return cloneUser(user), nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *inMemoryUserRepository) FindByID(_ context.Context, id string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, ok := r.users[uuidID]
	if !ok {
		return nil, errors.New("user not found")
	}

	return cloneUser(user), nil
}

func (r *inMemoryUserRepository) FindAll(_ context.Context) ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, cloneUser(user))
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Email < users[j].Email
	})

	return users, nil
}

func (r *inMemoryUserRepository) Patch(_ context.Context, id string, user *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	existing, ok := r.users[uuidID]
	if !ok {
		return errors.New("user not found")
	}

	if user.UserHandle != "" {
		existing.UserHandle = user.UserHandle
	}

	return nil
}

func (r *inMemoryUserRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if _, ok := r.users[uuidID]; !ok {
		return errors.New("user not found")
	}

	delete(r.users, uuidID)
	return nil
}

func (r *inMemoryUserRepository) IsFollowing(_ context.Context, followerID, followeeID uuid.UUID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.follows[followKey(followerID, followeeID)]
	return ok, nil
}

func (r *inMemoryUserRepository) CreateFollow(_ context.Context, follow *entities.Follow) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.follows[followKey(follow.FollowerID, follow.FolloweeID)] = follow
	return nil
}

func (r *inMemoryUserRepository) DeleteFollow(_ context.Context, followerID, followeeID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := followKey(followerID, followeeID)
	if _, ok := r.follows[key]; !ok {
		return errors.New("follow not found")
	}

	delete(r.follows, key)
	return nil
}

type testPublisher struct {
	topic   string
	payload map[string]interface{}
	err     error
}

func (p *testPublisher) Publish(_ context.Context, topic string, payload map[string]interface{}) error {
	p.topic = topic
	p.payload = payload
	return p.err
}

type UserUseCaseTestSuite struct {
	suite.Suite
	repo      *inMemoryUserRepository
	publisher *testPublisher
	service   usecase.UserUseCase
}

func (s *UserUseCaseTestSuite) SetupTest() {
	s.repo = newInMemoryUserRepository()
	s.publisher = &testPublisher{}
	s.service = usecase.NewUserServiceWithPublisher(s.repo, s.publisher)

	// Set JWT_SECRET for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-token-generation")
}

func (s *UserUseCaseTestSuite) TearDownTest() {
	os.Unsetenv("JWT_SECRET")
}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}

func (s *UserUseCaseTestSuite) TestRegister() {
	ctx := context.Background()

	user := &entities.User{
		Email:      "register@example.com",
		Password:   "password123",
		UserHandle: "register_user",
	}

	err := s.service.Register(ctx, user)
	s.NoError(err)
	s.NotEmpty(user.ID)

	// Verify password is hashed
	s.NotEqual("password123", user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	s.NoError(err)
}

func (s *UserUseCaseTestSuite) TestRegister_DuplicateEmail() {
	ctx := context.Background()

	user1 := &entities.User{
		Email:      "duplicate@example.com",
		Password:   "password123",
		UserHandle: "user_1",
	}
	err := s.service.Register(ctx, user1)
	s.NoError(err)

	// Try to register with same email
	user2 := &entities.User{
		Email:      "duplicate@example.com",
		Password:   "password456",
		UserHandle: "user_2",
	}
	err = s.service.Register(ctx, user2)
	s.Error(err)
	s.Equal(apperror.ErrAlreadyExists, err)
}

func (s *UserUseCaseTestSuite) TestLogin() {
	ctx := context.Background()

	// Register a user first
	user := &entities.User{
		Email:      "login@example.com",
		Password:   "password123",
		UserHandle: "login_user",
	}
	err := s.service.Register(ctx, user)
	s.NoError(err)

	// Login with correct credentials
	token, loggedInUser, err := s.service.Login(ctx, "login@example.com", "password123")
	s.NoError(err)
	s.NotEmpty(token)
	s.NotNil(loggedInUser)
	s.Equal(user.Email, loggedInUser.Email)
}

func (s *UserUseCaseTestSuite) TestLogin_WrongPassword() {
	ctx := context.Background()

	// Register a user first
	user := &entities.User{
		Email:      "wrongpass@example.com",
		Password:   "password123",
		UserHandle: "wrongpass_user",
	}
	err := s.service.Register(ctx, user)
	s.NoError(err)

	// Login with wrong password
	token, loggedInUser, err := s.service.Login(ctx, "wrongpass@example.com", "wrongpassword")
	s.Error(err)
	s.Empty(token)
	s.Nil(loggedInUser)
}

func (s *UserUseCaseTestSuite) TestLogin_UserNotFound() {
	ctx := context.Background()

	token, loggedInUser, err := s.service.Login(ctx, "notfound@example.com", "password123")
	s.Error(err)
	s.Empty(token)
	s.Nil(loggedInUser)
}

func (s *UserUseCaseTestSuite) TestFindUserByID() {
	ctx := context.Background()

	// Register a user first
	user := &entities.User{
		Email:      "findbyid@example.com",
		Password:   "password123",
		UserHandle: "find_user",
	}
	err := s.service.Register(ctx, user)
	s.NoError(err)

	// Find by ID
	found, err := s.service.FindUserByID(ctx, user.ID.String())
	s.NoError(err)
	s.NotNil(found)
	s.Equal(user.ID, found.ID)
	s.Equal(user.Email, found.Email)
}

func (s *UserUseCaseTestSuite) TestFindAllUsers() {
	ctx := context.Background()

	// Register multiple users
	users := []*entities.User{
		{Email: "all1@example.com", Password: "pass1", UserHandle: "user_1"},
		{Email: "all2@example.com", Password: "pass2", UserHandle: "user_2"},
		{Email: "all3@example.com", Password: "pass3", UserHandle: "user_3"},
	}

	for _, user := range users {
		err := s.service.Register(ctx, user)
		s.NoError(err)
	}

	// Find all
	allUsers, err := s.service.FindAllUsers(ctx)
	s.NoError(err)
	s.Len(allUsers, 3)
}

func (s *UserUseCaseTestSuite) TestPatchUser() {
	ctx := context.Background()

	// Register a user first
	user := &entities.User{
		Email:      "patch@example.com",
		Password:   "password123",
		UserHandle: "original_user",
	}
	err := s.service.Register(ctx, user)
	s.NoError(err)

	// Update user
	updateData := &entities.User{
		UserHandle: "updated_user",
	}
	updated, err := s.service.PatchUser(ctx, user.ID.String(), updateData)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("updated_user", updated.UserHandle)
	s.Equal(user.Email, updated.Email)
}

func (s *UserUseCaseTestSuite) TestDeleteUser() {
	ctx := context.Background()

	// Register a user first
	user := &entities.User{
		Email:      "delete@example.com",
		Password:   "password123",
		UserHandle: "delete_user",
	}
	err := s.service.Register(ctx, user)
	s.NoError(err)

	// Delete user
	err = s.service.DeleteUser(ctx, user.ID.String())
	s.NoError(err)

	// Verify deletion
	found, err := s.service.FindUserByID(ctx, user.ID.String())
	s.Error(err)
	s.Nil(found)
}

func (s *UserUseCaseTestSuite) TestFollowUser() {
	ctx := context.Background()
	followerID := uuid.New()
	followeeID := uuid.New()

	err := s.service.FollowUser(ctx, followerID, followeeID)
	s.NoError(err)

	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followeeID)
	s.NoError(err)
	s.True(isFollowing)
	s.Equal("social_events", s.publisher.topic)
	s.Equal("FOLLOW", s.publisher.payload["action_type"])
	s.Equal(followerID.String(), s.publisher.payload["actor_id"])
	s.Equal(followeeID.String(), s.publisher.payload["entity_id"])
	s.Equal("USER", s.publisher.payload["entity_type"])
}

func (s *UserUseCaseTestSuite) TestFollowUser_CannotFollowSelf() {
	ctx := context.Background()
	userID := uuid.New()

	err := s.service.FollowUser(ctx, userID, userID)
	s.Error(err)
	s.Equal(apperror.ErrInvalidData, err)
}

func (s *UserUseCaseTestSuite) TestFollowUser_AlreadyFollowing() {
	ctx := context.Background()
	followerID := uuid.New()
	followeeID := uuid.New()

	err := s.service.FollowUser(ctx, followerID, followeeID)
	s.NoError(err)

	err = s.service.FollowUser(ctx, followerID, followeeID)
	s.Error(err)
	s.Equal(apperror.ErrAlreadyExists, err)
}

func (s *UserUseCaseTestSuite) TestUnfollowUser() {
	ctx := context.Background()
	followerID := uuid.New()
	followeeID := uuid.New()

	err := s.service.FollowUser(ctx, followerID, followeeID)
	s.NoError(err)

	err = s.service.UnfollowUser(ctx, followerID, followeeID)
	s.NoError(err)

	isFollowing, err := s.repo.IsFollowing(ctx, followerID, followeeID)
	s.NoError(err)
	s.False(isFollowing)
}
