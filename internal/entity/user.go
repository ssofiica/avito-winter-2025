package entity

import (
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (s AuthRequest) Valid() bool {
	return s.Name != "" && s.Password != ""
}

type AuthResponse struct {
	Token string `json:"token"`
}

type AuthData struct {
	ID       uint
	Name     string
	Password string
}

type User struct {
	ID    uint32
	Name  string
	Coins uint32
}

type Password string

func (p *Password) IsEqual(comparing string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*p), []byte(comparing))
	return err == nil
}

func (p *Password) Hash(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		*p = Password("")
		return err
	}
	*p = Password(string(bytes))
	return nil
}
