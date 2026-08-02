package dto

type UserResponse struct {
	Id                int     `json:"id"`
	Email             string  `json:"email"`
	Name              string  `json:"name"`
	Bio               string  `json:"bio"`
	ProfilePictureURL *string `json:"profile_picture_url"`
}
