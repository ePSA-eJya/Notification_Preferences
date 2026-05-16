package usecase

import (
	"context"
	"time"

	"Notification_Preferences/internal/broker"
	"Notification_Preferences/internal/entities"
	feedRepoPkg "Notification_Preferences/internal/feed/repository"
	userRepoPkg "Notification_Preferences/internal/user/repository"

	"github.com/google/uuid"
)

type FeedUsecase struct {
	feedRepo   feedRepoPkg.FeedRepository
	userRepo   userRepoPkg.UserRepository
	publisher  *broker.KafkaProducer
	eventTopic string
}

func NewFeedUsecase(feedRepo feedRepoPkg.FeedRepository, userRepo userRepoPkg.UserRepository, publisher *broker.KafkaProducer, eventTopic string) *FeedUsecase {
	return &FeedUsecase{feedRepo: feedRepo, userRepo: userRepo, publisher: publisher, eventTopic: eventTopic}
}

func (s *FeedUsecase) CreatePost(ctx context.Context, post *entities.Post) error {
	if post.ID == uuid.Nil {
		post.ID = uuid.New()
	}
	post.CreatedAt = time.Now().UTC()

	// Fetch user handle for the author
	handle, err := s.userRepo.GetNameByUserID(ctx, post.UserID)
	if err == nil {
		post.UserHandle = handle
	}

	if err := s.feedRepo.CreatePost(ctx, post); err != nil {
		return err
	}

	// Fetch followers of the author and include the author themself so the post
	// appears immediately in their own feed timeline.
	followers, err := s.userRepo.GetFollowers(ctx, post.UserID)
	if err != nil {
		return err
	}

	// Build feed items and bulk insert
	items := make([]*entities.FeedItem, 0, len(followers)+1)
	items = append(items, &entities.FeedItem{
		ID:        uuid.New(),
		UserID:    post.UserID,
		PostID:    post.ID,
		AuthorID:  post.UserID,
		CreatedAt: post.CreatedAt,
	})
	for _, fid := range followers {
		if fid == post.UserID {
			continue
		}
		item := &entities.FeedItem{
			ID:        uuid.New(),
			UserID:    fid,
			PostID:    post.ID,
			AuthorID:  post.UserID,
			CreatedAt: post.CreatedAt,
		}
		items = append(items, item)
	}
	if len(items) > 0 {
		if err := s.feedRepo.AddFeedItems(ctx, items); err != nil {
			return err
		}
	}

	// Publish event to Kafka
	if s.publisher != nil {
		event := entities.Event{
			ID:         uuid.New(),
			ActorID:    post.UserID,
			EntityID:   post.ID,
			EntityType: "POST",
			ActionType: entities.Posted,
			CreatedAt:  time.Now().UTC(),
		}
		_ = s.publisher.Publish(ctx, s.eventTopic, event)
	}

	return nil
}

func (s *FeedUsecase) LikePost(ctx context.Context, like *entities.Like) error {
	if like.ID == uuid.Nil {
		like.ID = uuid.New()
	}
	like.CreatedAt = time.Now().UTC()

	if err := s.feedRepo.SaveLike(ctx, like); err != nil {
		return err
	}

	if s.publisher != nil {
		event := entities.Event{
			ID:         uuid.New(),
			ActorID:    like.UserID,
			EntityID:   like.PostID,
			EntityType: "POST",
			ActionType: entities.Liked,
			CreatedAt:  time.Now().UTC(),
		}
		_ = s.publisher.Publish(ctx, s.eventTopic, event)
	}
	return nil
}

func (s *FeedUsecase) UnlikePost(ctx context.Context, postID, userID uuid.UUID) error {
	return s.feedRepo.RemoveLike(ctx, postID, userID)
}

func (s *FeedUsecase) IsPostLikedByUser(ctx context.Context, postID, userID uuid.UUID) (bool, error) {
	return s.feedRepo.IsPostLikedByUser(ctx, postID, userID)
}

func (s *FeedUsecase) CommentOnPost(ctx context.Context, comment *entities.Comment) error {
	if comment.ID == uuid.Nil {
		comment.ID = uuid.New()
	}
	comment.CreatedAt = time.Now().UTC()
	if handle, err := s.userRepo.GetNameByUserID(ctx, comment.UserID); err == nil {
		comment.UserHandle = handle
	}

	if err := s.feedRepo.SaveComment(ctx, comment); err != nil {
		return err
	}

	if s.publisher != nil {
		event := entities.Event{
			ID:         uuid.New(),
			ActorID:    comment.UserID,
			EntityID:   comment.PostID,
			EntityType: "POST",
			ActionType: entities.Commented,
			CreatedAt:  time.Now().UTC(),
		}
		_ = s.publisher.Publish(ctx, s.eventTopic, event)
	}
	return nil
}

func (s *FeedUsecase) GetPostComments(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*entities.Comment, error) {
	comments, err := s.feedRepo.GetCommentsByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, err
	}

	for _, c := range comments {
		handle, err := s.userRepo.GetNameByUserID(ctx, c.UserID)
		if err == nil {
			c.UserHandle = handle
		}
	}

	return comments, nil
}

func (s *FeedUsecase) GetFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Post, error) {
	items, err := s.feedRepo.GetUserTimeline(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	posts := make([]*entities.Post, 0, len(items))
	for _, it := range items {
		p, err := s.feedRepo.GetPostByID(ctx, it.PostID.String())
		if err != nil {
			// skip missing posts
			continue
		}
		// Populate UserHandle
		handle, err := s.userRepo.GetNameByUserID(ctx, p.UserID)
		if err == nil {
			p.UserHandle = handle
		}
		posts = append(posts, p)
	}
	return posts, nil
}
