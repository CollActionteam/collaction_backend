package models

type CommitmentOption struct {
	Id          string             `json:"id"`
	Label       string             `json:"label"`
	Description string             `json:"description"`
	Requires    []CommitmentOption `json:"requires,omitempty"`
	Points      int                `json:"points"`
}
