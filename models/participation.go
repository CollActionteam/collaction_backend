package models

type ParticipationRecord struct {
	UserID        string   `json:"userID"`
	Name          string   `json:"name"`
	CrowdactionID string   `json:"crowdactionID"`
	Commitments   []string `json:"commitments,omitempty"`
	Timestamp     int64    `json:"timestamp"` // TODO use date instead?
}

type ParticipationEvent struct {
	UserID        string   `json:"userID"`
	CrowdactionID string   `json:"crowdactionID"`
	Commitments   []string `json:"commitments,omitempty"`
	Count         int      `json:"count"`
}
