package model

import (
	"errors"

	"gorm.io/gorm"
)

// UserModel define a user base info struct.
type UserModel struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"id"`
	Username  string `gorm:"column:username" json:"username"`
	Nickname  string `gorm:"column:nickname" json:"nickname"`
	Phone     string `gorm:"column:phone" json:"phone"`
	Email     string `gorm:"column:email" json:"email"`
	Password  string `gorm:"column:password" json:"password"`
	Avatar    string `gorm:"column:avatar" json:"avatar"`
	Gender    string `gorm:"column:gender" json:"gender"`
	Birthday  string `gorm:"column:birthday" json:"birthday"`
	Bio       string `gorm:"column:bio" json:"bio"`
	LoginAt   int64  `gorm:"column:login_at" json:"login_at"` // login time for last times
	Status    int32  `gorm:"column:status" json:"status"`
	CreatedAt int64  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt int64  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名
func (u *UserModel) TableName() string {
	return "user_info"
}

// BeforeCreate a prefix func for create
func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Username == "admin" || u.Username == "root" || u.Username == "administrator" {
		return errors.New("invalid name")
	}
	return
}
