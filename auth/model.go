package auth

type UserCredential struct {
	UserID   string `json:"user_id" validate:"min=6,max=20"`
	Password string `json:"password" validate:"min=8,max=20"`
}

type UserUpdate struct {
	UserID   string `json:"user_id,omitempty"`
	Password string `json:"password,omitempty"`
	Nickname string `json:"nickname,omitempty" validate:"max=29"`
	Comment  string `json:"comment,omitempty" validate:"max=99"`
}

type User struct {
	UserID   string `json:"user_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

type UserData struct {
	User
	Password string
}
