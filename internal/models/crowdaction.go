package models

type CrowdActionRequest struct {
	Data CrowdActionData `json:"data" validate:"required"`
}

type CrowdactionParticipant struct {
	Name   string `json:"name,omitempty"`
	UserID string `json:"userID,omitempty"`
}

type CrowdactionImages struct {
	Card   string `json:"card,omitempty"`
	Banner string `json:"banner,omitempty"`
}

type CrowdActionData struct {
	CrowdactionID      string                   `json:"crowdactionID"`
	Title              string                   `json:"title"`
	Description        string                   `json:"description"`
	Category           string                   `json:"category"`
	Subcategory        string                   `json:"subcategory"`
	Location           string                   `json:"location"`
	DateStart          string                   `json:"date_start"`
	DateEnd            string                   `json:"date_end"`
	DateLimitJoin      string                   `json:"date_limit_join"`
	PasswordJoin       string                   `json:"password_join"`
	CommitmentOptions  []CommitmentOption       `json:"commitment_options"`
	ParticipationCount int                      `json:"participant_count"`
	TopParticipants    []CrowdactionParticipant `json:"top_participants"`
	Images             CrowdactionImages        `json:"images"`
}
