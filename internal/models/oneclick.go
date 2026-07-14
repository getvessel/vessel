package models

type OneClickApp struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Port        int    `json:"port"`
}
