package participation

type ParticipationEvent struct {
	UserID        string `json:"userID"`
	CrowdactionID string `json:"crowdactionID"`
	Count         int    `json:"count"`
}
