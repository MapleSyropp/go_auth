package models

type UserReq struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

func NewUserReq(username, password string) *UserReq {
	return &UserReq{
		Username: username,
		Password: password,
	}
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func NewUser(id int, username, password string) *User {
	return &User{
		ID:       id,
		Username: username,
		Password: password,
	}
}

type ApiError struct {
	Error string
}
