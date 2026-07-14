package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
)

// AuthResponse holds the response from the authentication endpoint.
type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

// Login authenticates a user with email and password and returns a token.
func (c *Client) Login(email, password string) (*AuthResponse, error) {
	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	resp, err := c.sendRequest("POST", "/auth/login", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed (status %d): %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	// Automatically set the token on the client for future requests
	c.Token = authResp.Token

	return &authResp, nil
}
