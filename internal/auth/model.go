package auth

import "vessel.dev/vessel/internal/user"

// AuthResult represents the response returned upon successful authentication.
type AuthResult struct {
	Token string     `json:"token"`
	User  *user.User `json:"user"`
}
