package user

type UpdateProfileRequest struct {
	Email string `json:"email"`
}

type CreatePATRequest struct {
	Name string `json:"name"`
}

type CreatePATResponse struct {
	Token *PersonalAccessToken `json:"token"`
	Plain string               `json:"plain"`
}
