package tests

import (
	"techwave/models"
	"time"
)

// DepartmentFactory provides a builder pattern for creating test department data
type DepartmentFactory struct {
	department models.Department
}

// NewDepartmentFactory creates a new factory with default valid values
func NewDepartmentFactory() *DepartmentFactory {
	return &DepartmentFactory{
		department: models.Department{
			ID:             1,
			DepartmentCode: "CS",
			DepartmentName: "Computer Science",
			AnnualBudget:   5000000.00,
			DepartmentHead: "Dr. John Smith",
			Status:         "Active",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}
}

// WithID sets the department ID
func (f *DepartmentFactory) WithID(id int) *DepartmentFactory {
	f.department.ID = id
	return f
}

// WithDepartmentName sets the department name
func (f *DepartmentFactory) WithDepartmentName(name string) *DepartmentFactory {
	f.department.DepartmentName = name
	return f
}

// WithDepartmentCode sets the department code
func (f *DepartmentFactory) WithDepartmentCode(code string) *DepartmentFactory {
	f.department.DepartmentCode = code
	return f
}

// WithAnnualBudget sets the annual budget
func (f *DepartmentFactory) WithAnnualBudget(budget float64) *DepartmentFactory {
	f.department.AnnualBudget = budget
	return f
}

// WithDepartmentHead sets the department head name
func (f *DepartmentFactory) WithDepartmentHead(headName string) *DepartmentFactory {
	f.department.DepartmentHead = headName
	return f
}

// WithStatus sets the department status
func (f *DepartmentFactory) WithStatus(status string) *DepartmentFactory {
	f.department.Status = status
	return f
}

// WithCreatedAt sets the created timestamp
func (f *DepartmentFactory) WithCreatedAt(t time.Time) *DepartmentFactory {
	f.department.CreatedAt = t
	return f
}

// WithUpdatedAt sets the updated timestamp
func (f *DepartmentFactory) WithUpdatedAt(t time.Time) *DepartmentFactory {
	f.department.UpdatedAt = t
	return f
}

// Build returns the constructed department
func (f *DepartmentFactory) Build() models.Department {
	return f.department
}

// Common scenario helpers

// NewValidDepartment creates a department with valid data
func NewValidDepartment() models.Department {
	return NewDepartmentFactory().Build()
}

// NewDepartmentWithoutHeadName creates a department without head name
func NewDepartmentWithoutHeadName() models.Department {
	return NewDepartmentFactory().
		WithDepartmentHead("").
		Build()
}

// Edge case helpers for invalid data

// NewDepartmentWithEmptyName creates a department with empty name
func NewDepartmentWithEmptyName() models.Department {
	return NewDepartmentFactory().
		WithDepartmentName("").
		Build()
}

// NewDepartmentWithEmptyCode creates a department with empty code
func NewDepartmentWithEmptyCode() models.Department {
	return NewDepartmentFactory().
		WithDepartmentCode("").
		Build()
}

// NewDepartmentWithShortName creates a department with too short name
func NewDepartmentWithShortName() models.Department {
	return NewDepartmentFactory().
		WithDepartmentName("A").
		Build()
}

// NewDepartmentWithLongName creates a department with too long name
func NewDepartmentWithLongName() models.Department {
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "A"
	}
	return NewDepartmentFactory().
		WithDepartmentName(longName).
		Build()
}

// NewDepartmentWithShortCode creates a department with too short code
func NewDepartmentWithShortCode() models.Department {
	return NewDepartmentFactory().
		WithDepartmentCode("A").
		Build()
}

// NewDepartmentWithLongCode creates a department with too long code
func NewDepartmentWithLongCode() models.Department {
	return NewDepartmentFactory().
		WithDepartmentCode("TOOLONG").
		Build()
}

// Batch creation helpers

// NewDepartmentBatch creates multiple departments with different IDs and codes
func NewDepartmentBatch(count int) []models.Department {
	departments := make([]models.Department, count)
	departmentNames := []string{"Computer Science", "Mathematics", "Physics", "Chemistry", "Biology"}
	departmentCodes := []string{"CS", "MATH", "PHYS", "CHEM", "BIO"}

	for i := 0; i < count; i++ {
		nameIdx := i % len(departmentNames)
		departments[i] = NewDepartmentFactory().
			WithID(i + 1).
			WithDepartmentName(departmentNames[nameIdx]).
			WithDepartmentCode(departmentCodes[nameIdx]).
			Build()
	}
	return departments
}

// Scenario helpers for specific test cases

// NewDepartmentWithCode creates a department with a specific code
func NewDepartmentWithCode(code string) models.Department {
	return NewDepartmentFactory().
		WithDepartmentCode(code).
		Build()
}

// NewDepartmentWithNameAndCode creates a department with specific name and code
func NewDepartmentWithNameAndCode(name, code string) models.Department {
	return NewDepartmentFactory().
		WithDepartmentName(name).
		WithDepartmentCode(code).
		Build()
}

// NewCommonDepartments creates a set of commonly used departments for testing
func NewCommonDepartments() []models.Department {
	return []models.Department{
		NewDepartmentFactory().WithID(1).WithDepartmentName("Computer Science").WithDepartmentCode("CS").Build(),
		NewDepartmentFactory().WithID(2).WithDepartmentName("Mathematics").WithDepartmentCode("MATH").Build(),
		NewDepartmentFactory().WithID(3).WithDepartmentName("Physics").WithDepartmentCode("PHYS").Build(),
		NewDepartmentFactory().WithID(4).WithDepartmentName("Chemistry").WithDepartmentCode("CHEM").Build(),
		NewDepartmentFactory().WithID(5).WithDepartmentName("Biology").WithDepartmentCode("BIO").Build(),
	}
}
