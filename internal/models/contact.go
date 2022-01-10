package models

type EmailContactRequest struct {
	Email   string `json:"email" validate:"required,email" binding:"required"`
	Subject string `json:"subject" validate:"required,lte=50" binding:"required"`
	Message string `json:"message" validate:"required,lte=500" binding:"required"`
	//TODO 11.01.22 mrsoftware: fix regx
	AppVersion string `json:"app_version" validate:"required" binding:"required"` // ,regexp=^(?:ios|android) [0-9]+\\.[0-9]+\\.[0-9]+\\+[0-9]+$
}

type EmailData struct {
	Recipient  string
	Message    string
	Subject    string
	Sender     string
	ReplyEmail string
}
