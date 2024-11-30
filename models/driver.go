package models

type Driver struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NameAcronym string `json:"name_acronym"`
	Number      int    `json:"driver_number"`
	TeamName    string `json:"team_name"`
}
