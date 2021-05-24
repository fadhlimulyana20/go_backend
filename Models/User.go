package models

import (
	"errors"
	"log"
	"time"

	"github.com/fadhlimulyana20/go_backend/utils"
	"gorm.io/gorm"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" validate:"required,email"`
	FirstName   string    `json:"firstName" validate:"required"`
	LastName    string    `json:"lastName" validate:"required"`
	Password    string    `json:"password" validate:"required,min=8"`
	IsValidated bool      `json:"isValidated" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *User) BeforeSave(db *gorm.DB) error {
	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	u.Password = hashedPassword
	return nil
}

func (u *User) FindUserById(db *gorm.DB, uid uint) (*User, error) {
	err := db.First(&u, uid).Error

	if err != nil {
		return &User{}, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("User not found")
	}

	return u, nil
}

func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	err := db.Where("email = ?", email).First(&u).Error

	if err != nil {
		return &User{}, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("User not found")
	}

	return u, nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	// Create new user record
	res := db.Create(&u)
	if err := res.Error; err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) UpdateUser(db *gorm.DB, uid uint) (*User, error) {
	// to hash password
	err := u.BeforeSave(db)
	if err != nil {
		log.Fatal(err)
	}

	db = db.Model(&User{}).Where("id", uid).Updates(
		map[string]interface{}{
			"email":      u.Email,
			"password":   u.Password,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
		},
	)

	if err := db.Error; err != nil {
		return &User{}, err
	}

	return u, nil
}
