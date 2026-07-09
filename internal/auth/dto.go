package auth

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
