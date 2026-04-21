package dto

import "Notification_Preferences/internal/entities"

// From entity.User to UserResponse
func ToUserResponse(user *entities.User) *UserResponse {
	return &UserResponse{
<<<<<<< HEAD
		ID:    user.ID,
		Email: user.Email,
=======
		ID:         user.ID,
		Email:      user.Email,
		UserHandle: user.UserHandle,
>>>>>>> 3f79743 (Add follow/unfollow flow)
	}
}

func ToUserResponseList(users []*entities.User) []*UserResponse {
	responses := make([]*UserResponse, len(users))
	for i, u := range users {
		responses[i] = ToUserResponse(u)
	}
	return responses
}

// From RegisterRequest to entity.User (optional, if want to use in usecase)
func ToUserEntity(req *RegisterRequest) *entities.User {
	return &entities.User{
<<<<<<< HEAD
		Email:    req.Email,
		Password: req.Password,
=======
		Email:      req.Email,
		Password:   req.Password,
		UserHandle: req.UserHandle,
>>>>>>> 3f79743 (Add follow/unfollow flow)
	}
}
