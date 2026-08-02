package dto

type UpdateProfileRequest struct {
	Name *string `json:"name" binding:"omitempty,max=100"`
	Bio  *string `json:"bio" binding:"omitempty,max=500"`
}
