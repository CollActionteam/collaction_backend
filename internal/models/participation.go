package models

type JoinPayload struct {
	Password    string   `json:"password,omitempty"`
	Commitments []string `json:"commitments,omitempty"`
}

type ParticipationRecord struct {
	UserID        string   `json:"userID"`
	Name          string   `json:"name"`
	CrowdactionID string   `json:"crowdactionID"`
	Title         string   `json:"title"`
	Commitments   []string `json:"commitments,omitempty"`
	Date          string   `json:"date"`
}

type ParticipationEvent struct {
	UserID        string   `json:"userID"`
	CrowdactionID string   `json:"crowdactionID"`
	Commitments   []string `json:"commitments,omitempty"`
	Count         int      `json:"count"`
}
