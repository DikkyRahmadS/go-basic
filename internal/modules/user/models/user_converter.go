package models

func ToUserResponse(user *User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserResponses(users []User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for i := range users {
		responses = append(responses, ToUserResponse(&users[i]))
	}
	return responses
}
