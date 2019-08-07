package models

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

type Success struct {
	Message string `json:"message"`
}
