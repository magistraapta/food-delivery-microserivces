package models

type Error struct {
	Error   int    `json:"error" example:"27"`
	Message string `json:"message" example:"User not found"`
}
