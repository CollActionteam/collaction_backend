package main

type Crowdaction struct {
	CrowdactionID string `json:"crowdactionID,omitempty"`
	Title         string `json:"title,omitempty"`
	Description   string `json:"description,omitempty"`
	Category      string `json:"category,omitempty"`
	Subcategory   string `json:"subcategory,omitempty"`
	Location      string `json:"location,omitempty"`
	DateStart     string `json:"date_start,omitempty"`
	DateEnd       string `json:"date_end,omitempty"`
	DateLimitJoin string `json:"date_limit_join,omitempty"`
	Code          string `json:"code,omitempty"`
}
