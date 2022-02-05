package models

type CrowdActionRequest struct {
	Title       string
	Description string
	AppVersion  string
}

type CrowdAction struct {
	CrowdactionID string
	Title         string
	Description   string
}
