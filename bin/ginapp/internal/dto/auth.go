package dto


type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
