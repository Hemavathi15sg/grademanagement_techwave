package models

import (
	"errors"
	"time"
)

// Department represents an academic department
type Department struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=2,max=100"`
	Code      string    `json:"code" validate:"required,min=2,max=10,uppercase"`
	HeadName  string    `json:"head_name,omitempty" validate:"omitempty,min=2,max=100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate checks if the department has all required fields
func (d *Department) Validate() error {
	if d.Name == "" {
		return errors.New("name is required")
	}
	if len(d.Name) < 2 || len(d.Name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}
	if d.Code == "" {
		return errors.New("code is required")
	}
	if len(d.Code) < 2 || len(d.Code) > 10 {
		return errors.New("code must be between 2 and 10 characters")
	}
	if d.HeadName != "" && (len(d.HeadName) < 2 || len(d.HeadName) > 100) {
		return errors.New("head_name must be between 2 and 100 characters when provided")
	}
	return nil
}
