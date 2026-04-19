package entities

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email    string    `gorm:"uniqueIndex" json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
}

func (u *User) BeforeCreate(tx *mongo.Database) (err error) {
	u.ID = uuid.New()
	return
}
