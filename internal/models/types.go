package models

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func NewUser(name, password string) *User {
	return &User{
		Name:     name,
		Password: password,
	}
}

// test
type response struct {
	Data string `json:"data"`
}
