package model

import (
	"time"

	"github.com/shopspring/decimal"
)

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
	ID              int64
	Email           string
	Username        string
	Password        string
	CreatedAt       time.Time
	Bio             string
	Birthday        time.Time
	Phone           string
	ProfileImageURL string
}

func (User) TableName() string {
	return "user_t"
}

type Post struct {
	ID           int64
	Title        string
	Description  string
	Image        string
	UserID       int64
	MinimumFund  decimal.Decimal
	FundRaised   decimal.Decimal
	DeadlineDate time.Time
	CreatedAt    time.Time
	Category     string
}

func (Post) TableName() string {
	return "post"
}

type PostDetail struct {
	ID          int64
	Order       int32 `gorm:"column:order_c"`
	Title       string
	Description string
	Image       string
	PostID      int64
}

func (PostDetail) TableName() string {
	return "post_detail"
}

type Fund struct {
	ID     int64
	UserID int64
	PostID int64
	Amount decimal.Decimal
}

func (Fund) TableName() string {
	return "fund"
}

type Bookmark struct {
	ID     int64
	UserID int64
	PostID int64
}

func (Bookmark) TableName() string {
	return "bookmark"
}
