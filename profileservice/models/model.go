package models

type Profile struct {
	UserId      string `json:"userid"`
	DisplayName string `json:"displayname"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Bio         string `json:"bio"`
}
