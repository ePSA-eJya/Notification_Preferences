package entities

import (
	"time"

	"github.com/google/uuid"
)

// --- Custom Types for Preferences ---

type ChannelType string

const (
	InAppChannel ChannelType = "InApp"
	PushChannel  ChannelType = "Push"
	EmailChannel ChannelType = "Email"
)

type PreferenceLevel string

const (
	PrefAll       PreferenceLevel = "ALL"
	PrefFollowers PreferenceLevel = "FOLLOWERS"
	PrefFollowing PreferenceLevel = "FOLLOWING"
	PrefNone      PreferenceLevel = "NONE"
)

type ChannelConfig struct {
	InApp PreferenceLevel `bson:"in_app" json:"in_app"`
	Push  PreferenceLevel `bson:"push" json:"push"`
	Email PreferenceLevel `bson:"email" json:"email"`
}

type NotificationPreferences struct {
	Likes    ChannelConfig `bson:"likes" json:"likes"`
	Comments ChannelConfig `bson:"comments" json:"comments"`
	Follows  ChannelConfig `bson:"follows" json:"follows"`
	Posts    ChannelConfig `bson:"posts" json:"posts"`
}

// --- Main Entities ---

type User struct {
	ID          uuid.UUID               `bson:"_id" json:"id"`
	UserHandle  string                  `bson:"user_handle" json:"user_handle"`
	Email       string                  `bson:"email" json:"email"`
	Password    string                  `bson:"password" json:"-"` // Omitted from JSON responses
	Preferences NotificationPreferences `bson:"preferences" json:"preferences"`
	DeviceToken string                  `bson:"device_token"`
}

func (u *User) Initialize() {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
}

type Follow struct {
	ID         uuid.UUID `bson:"_id" json:"id"`
	FollowerID uuid.UUID `bson:"follower_id" json:"follower_id"`
	FolloweeID uuid.UUID `bson:"followee_id" json:"followee_id"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
