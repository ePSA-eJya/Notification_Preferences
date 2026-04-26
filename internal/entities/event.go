package entities

import (
	"time"

	"github.com/google/uuid"
)

type ActionType string

const (
	Posted    ActionType = "POSTED"
	Liked     ActionType = "LIKED"
	Commented ActionType = "COMMENTED"
	Followed  ActionType = "FOLLOWED"
)

// <Actor> did <ActionType> on <EntityType>
type Event struct {
	ID         uuid.UUID  `bson:"_id" json:"id"`
	ActorID    uuid.UUID  `bson:"actor_id" json:"actor_id"`
	ActionType ActionType `bson:"action_type" json:"action_type"` // e.g., "COMMENT", "LIKE"
	EntityID   uuid.UUID  `bson:"entity_id" json:"entity_id"`
	EntityType string     `bson:"entity_type" json:"entity_type"` // e.g., "POST"
	Payload    string     `bson:"payload" json:"payload"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
}
