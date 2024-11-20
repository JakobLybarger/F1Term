package models

type Session struct {
	Name string `json:"session_name"`
	Key  int    `json:"session_key"`
}
