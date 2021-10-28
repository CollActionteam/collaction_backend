package profileservice

type Profile struct {
	UserID      string `json:"userid,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
	Country     string `json:"country,omitempty"`
	City        string `json:"city,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Phone       string `json:"phone,omitempty"`
}
