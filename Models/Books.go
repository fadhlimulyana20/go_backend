package models

import "time"

type Book struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateBookDTO struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
