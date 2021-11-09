package participation

type ParticipationRecord struct {
	UserID        string `json:"userID"`
	Name          string `json:"name"`
	CrowdactionID string `json:"crowdactionID"`
	CommitmentID  string `json:"commitmentID"`
	Timestamp     int64  `json:"timestamp"`
}

type ParticipationEvent struct {
	UserID        string `json:"userID"`
	CrowdactionID string `json:"crowdactionID"`
	CommitmentID  string `json:"commitmentID"`
	Count         int    `json:"count"`
}

type JoinPayload struct {
	CommitmentID string `json:"commitmentID"`
}
