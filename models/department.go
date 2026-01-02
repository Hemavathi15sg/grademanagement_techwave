package models

import (
	"errors"
	"regexp"
	"time"
)

// Department represents an academic department
type Department struct {
	ID             int       `json:"id"`
	DepartmentCode string    `json:"department_code" validate:"required,min=2,max=6,uppercase"`
	DepartmentName string    `json:"department_name" validate:"required,min=2,max=100"`
	AnnualBudget   float64   `json:"annual_budget" validate:"gte=0,lte=10000000"`
	DepartmentHead string    `json:"department_head,omitempty" validate:"omitempty,min=2,max=100"`
	Status         string    `json:"status" validate:"oneof=Active Inactive"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// BudgetEntry represents a budget change history entry
type BudgetEntry struct {
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Note      string    `json:"note,omitempty"`
}

// BudgetUtilization represents budget usage statistics
type BudgetUtilization struct {
	Allocated  float64 `json:"allocated"`
	Spent      float64 `json:"spent"`
	Remaining  float64 `json:"remaining"`
	Percentage float64 `json:"percentage"`
}

var departmentCodePattern = regexp.MustCompile(`^[A-Z]{2,6}$`)

// Validate checks if the department has all required fields
func (d *Department) Validate() error {
	if d.DepartmentName == "" {
		return errors.New("department_name is required")
	}
	if len(d.DepartmentName) < 2 || len(d.DepartmentName) > 100 {
		return errors.New("department_name must be between 2 and 100 characters")
	}
	if d.DepartmentCode == "" {
		return errors.New("department_code is required")
	}
	if !departmentCodePattern.MatchString(d.DepartmentCode) {
		return errors.New("department_code must be 2-6 uppercase letters (e.g., CS, MATH, ENG)")
	}
	if d.AnnualBudget < 0 {
		return errors.New("annual_budget cannot be negative")
	}
	if d.AnnualBudget > 10000000 {
		return errors.New("annual_budget cannot exceed $10,000,000")
	}
	if d.Status != "" && d.Status != "Active" && d.Status != "Inactive" {
		return errors.New("status must be either 'Active' or 'Inactive'")
	}
	if d.Status == "" {
		d.Status = "Active" // Default status
	}
	if d.DepartmentHead != "" && (len(d.DepartmentHead) < 2 || len(d.DepartmentHead) > 100) {
		return errors.New("department_head must be between 2 and 100 characters when provided")
	}
	return nil
}
