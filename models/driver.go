package models

type Driver struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Number    int    `json:"driver_number"`
	TeamName  string `json:"team_name"`
}
