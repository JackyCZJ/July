package store

type UserModel struct {
	Username string `json:"username" validate:"min=1,max=32"`
	Password string `json:"password,omitempty" validate:"min=1,max=32"`
}

type UserInformation struct {
	*UserModel `json:"user"`
	UserId     string `json:"user_id"`
	Email      string `json:"email"`
	Rank       int    `json:"rank"`
	Gander     int    `json:"gander"`
	Phone      string `json:"phone"`
}
