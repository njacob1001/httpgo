package routes

// UserAuthResponse struct
type UserAuthResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type UserVerification struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
