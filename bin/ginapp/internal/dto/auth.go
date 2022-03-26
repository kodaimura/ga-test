package dto


type LoginDto struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type SignupDto struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
