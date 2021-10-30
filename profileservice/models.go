package profileservice

type Profile struct {
	UserID      string `json:"userid,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
	Country     string `json:"country,omitempty"`
	City        string `json:"city,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Phone       string `json:"phone,omitempty"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}
