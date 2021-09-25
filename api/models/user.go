package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User contains the user information of a user
type User struct {
	Model
	Username     string `json:"username" validate:"required" gorm:"unique"`
	Password     string `json:"password,omitempty" validate:"required,gte=8" gorm:"-"`
	PasswordHash string `json:"-"`
}

// BeforeSave hashes the password and removes the password before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password == "" {
		// Don't do anything to password hash if nothing is changed
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = ""
	if err != nil {
		return err
	}

	u.PasswordHash = string(passwordHash)

	return nil
}
