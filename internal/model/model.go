package model

import "time"

type SignUp struct {
	ID               int32
	Email            string
	Username         string
	Password         string
	VerificationCode string
	Expire           time.Time
}

func (SignUp) TableName() string {
	return "sign_up"
}

type User struct {
	ID        int64
	Email     string
	Username  string
	Password  string
	CreatedAt time.Time
}

func (User) TableName() string {
	return "user_t"
}
