package routes

// UserAuthResponse struct
type UserAuthResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// UserVerification struct
type UserVerification struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ArticlesIdentidicators struct
type ArticlesIdentidicators struct {
	Articles []string `json:"articles"`
}
