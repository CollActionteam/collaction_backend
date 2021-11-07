package main

type Crowdaction struct {
	CrowdactionID    string `json:"crowdactionID,omitempty"`
	Title            string `json:"title,omitempty"`
	DescriptionShort string `json:"description_short,omitempty"`
	DescriptionLong  string `json:"description_long,omitempty"`
	StartDate        string `json:"start_date,omitempty"`
	EndDate          string `json:"end_date,omitempty"`
}
