package auth

import "vessel.dev/vessel/internal/user"

type AuthResult struct {
	Token string     `json:"token"`
	User  *user.User `json:"user"`
}
