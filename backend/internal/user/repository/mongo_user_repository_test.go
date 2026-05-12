package repository_test

import (
	"context"
	"testing"

	"Notification_Preferences/internal/entities"
	"Notification_Preferences/internal/user/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoUserRepository_Save(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		user := &entities.User{ID: uuid.New(), Email: "test@example.com", Password: "password123", UserHandle: "test_user"}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := repo.Save(context.Background(), user)
		require.NoError(mt, err)
	})
}

func TestMongoUserRepository_FindByEmail(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		doc := bson.D{{Key: "_id", Value: uuid.New()}, {Key: "email", Value: "find@example.com"}, {Key: "password", Value: "hash"}, {Key: "user_handle", Value: "find_user"}}
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, doc),
			mtest.CreateCursorResponse(0, "test.users", mtest.NextBatch),
		)

		user, err := repo.FindByEmail(context.Background(), "find@example.com")
		require.NoError(mt, err)
		require.NotNil(mt, user)
		require.Equal(mt, "find@example.com", user.Email)
		require.Equal(mt, "find_user", user.UserHandle)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.users", mtest.FirstBatch))

		user, err := repo.FindByEmail(context.Background(), "missing@example.com")
		require.Error(mt, err)
		require.Nil(mt, user)
		require.Equal(mt, "user not found", err.Error())
	})
}

func TestMongoUserRepository_FindByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		id := uuid.New()
		doc := bson.D{{Key: "_id", Value: id}, {Key: "email", Value: "id@example.com"}, {Key: "password", Value: "hash"}, {Key: "user_handle", Value: "id_user"}}
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, doc),
			mtest.CreateCursorResponse(0, "test.users", mtest.NextBatch),
		)

		user, err := repo.FindByID(context.Background(), id.String())
		require.NoError(mt, err)
		require.NotNil(mt, user)
		require.Equal(mt, id, user.ID)
	})

	mt.Run("invalid uuid", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		user, err := repo.FindByID(context.Background(), "invalid-id")
		require.Error(mt, err)
		require.Nil(mt, user)
	})
}

func TestMongoUserRepository_FindAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		doc1 := bson.D{{Key: "_id", Value: uuid.New()}, {Key: "email", Value: "one@example.com"}, {Key: "password", Value: "hash"}, {Key: "user_handle", Value: "one"}}
		doc2 := bson.D{{Key: "_id", Value: uuid.New()}, {Key: "email", Value: "two@example.com"}, {Key: "password", Value: "hash"}, {Key: "user_handle", Value: "two"}}
		mt.AddMockResponses(
			mtest.CreateCursorResponse(2, "test.users", mtest.FirstBatch, doc1, doc2),
			mtest.CreateCursorResponse(0, "test.users", mtest.NextBatch),
		)

		users, err := repo.FindAll(context.Background())
		require.NoError(mt, err)
		require.Len(mt, users, 2)
	})
}

func TestMongoUserRepository_Patch(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: int32(1)},
			bson.E{Key: "nModified", Value: int32(1)},
		))

		err := repo.Patch(context.Background(), uuid.New().String(), &entities.User{UserHandle: "updated"})
		require.NoError(mt, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "n", Value: int32(0)},
			bson.E{Key: "nModified", Value: int32(0)},
		))

		err := repo.Patch(context.Background(), uuid.New().String(), &entities.User{UserHandle: "updated"})
		require.Error(mt, err)
		require.Equal(mt, "user not found", err.Error())
	})
}

func TestMongoUserRepository_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: int32(1)}))

		err := repo.Delete(context.Background(), uuid.New().String())
		require.NoError(mt, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: int32(0)}))

		err := repo.Delete(context.Background(), uuid.New().String())
		require.Error(mt, err)
		require.Equal(mt, "user not found", err.Error())
	})
}

func TestMongoUserRepository_IsFollowing(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("following exists", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		followerID := uuid.New()
		followeeID := uuid.New()

		doc := bson.D{
			{Key: "_id", Value: uuid.New()},
			{Key: "follower_id", Value: followerID},
			{Key: "followee_id", Value: followeeID},
		}
		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.follows", mtest.FirstBatch, doc),
			mtest.CreateCursorResponse(0, "test.follows", mtest.NextBatch),
		)

		isFollowing, err := repo.IsFollowing(context.Background(), followerID, followeeID)
		require.NoError(mt, err)
		require.True(mt, isFollowing)
	})

	mt.Run("not following", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "test.follows", mtest.FirstBatch))

		isFollowing, err := repo.IsFollowing(context.Background(), uuid.New(), uuid.New())
		require.NoError(mt, err)
		require.False(mt, isFollowing)
	})
}

func TestMongoUserRepository_CreateFollow(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		follow := &entities.Follow{
			ID:         uuid.New(),
			FollowerID: uuid.New(),
			FolloweeID: uuid.New(),
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := repo.CreateFollow(context.Background(), follow)
		require.NoError(mt, err)
	})
}

func TestMongoUserRepository_DeleteFollow(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: int32(1)}))

		err := repo.DeleteFollow(context.Background(), uuid.New(), uuid.New())
		require.NoError(mt, err)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := repository.NewMongoUserRepository(mt.DB)
		mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: int32(0)}))

		err := repo.DeleteFollow(context.Background(), uuid.New(), uuid.New())
		require.Error(mt, err)
		require.Equal(mt, "follow not found", err.Error())
	})
}
