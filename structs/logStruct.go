package structs

type Logs struct {
	TimeCreated      string `json:"timeCreated"`
	ID               int    `json:"id"`
	LevelDisplayName string `json:"levelDisplayName"`
	Message          string `json:"message"`
}
