package models

type Meeting struct {
	MeetingKey   int    `json:"meeting_key"`
	OfficialName string `json:"meeting_official_name"`
}
