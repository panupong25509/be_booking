package models

type Error struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}

type Success struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}
