package models

import (
	"errors"
	"time"
)

// LetterGrade represents a letter grade with associated color from Figma design system
type LetterGrade struct {
	Letter      string  `json:"letter"`
	MinPoints   float64 `json:"min_points"`
	MaxPoints   float64 `json:"max_points"`
	Color       string  `json:"color"`        // Hex color from Figma
	Description string  `json:"description"`  // e.g., "Excellent", "Good"
}

// GradeScale defines the complete grading scale from Figma design
var GradeScale = []LetterGrade{
	{Letter: "A+", MinPoints: 97, MaxPoints: 100, Color: "#22C55E", Description: "Excellent"},
	{Letter: "A", MinPoints: 93, MaxPoints: 96, Color: "#22C55E", Description: "Excellent"},
	{Letter: "A-", MinPoints: 90, MaxPoints: 92, Color: "#86EFAC", Description: "Very Good"},
	{Letter: "B+", MinPoints: 87, MaxPoints: 89, Color: "#3B82F6", Description: "Good"},
	{Letter: "B", MinPoints: 83, MaxPoints: 86, Color: "#3B82F6", Description: "Good"},
	{Letter: "B-", MinPoints: 80, MaxPoints: 82, Color: "#60A5FA", Description: "Above Average"},
	{Letter: "C+", MinPoints: 77, MaxPoints: 79, Color: "#FBBF24", Description: "Average"},
	{Letter: "C", MinPoints: 73, MaxPoints: 76, Color: "#FBBF24", Description: "Average"},
	{Letter: "C-", MinPoints: 70, MaxPoints: 72, Color: "#FCD34D", Description: "Below Average"},
	{Letter: "D", MinPoints: 60, MaxPoints: 69, Color: "#F97316", Description: "Poor"},
	{Letter: "F", MinPoints: 0, MaxPoints: 59, Color: "#DC2626", Description: "Fail"},
}

// WeightedComponent represents a single component in grade calculation
type WeightedComponent struct {
	Name   string  `json:"name"`   // e.g., "assignments", "exams", "attendance"
	Score  float64 `json:"score"`  // Actual score (0-100)
	Weight float64 `json:"weight"` // Weight percentage (0-1, sum must equal 1.0)
}

// GradeCurve defines how to adjust scores
type GradeCurve struct {
	Enabled    bool    `json:"enabled"`
	Adjustment float64 `json:"adjustment"` // Percentage to add to all grades
}

// CalculationRequest contains all parameters for grade calculation
type CalculationRequest struct {
	Components []WeightedComponent `json:"components" validate:"required,min=1,dive"`
	Curve      *GradeCurve         `json:"curve,omitempty"`
}

// CalculationResult represents the calculated grade for a single student
type CalculationResult struct {
	NumericGrade  float64      `json:"numeric_grade"`
	LetterGrade   LetterGrade  `json:"letter_grade"`
	CurveApplied  bool         `json:"curve_applied"`
	CurveAdjustment float64    `json:"curve_adjustment,omitempty"`
}

// CalculationResponse contains the result and metadata
type CalculationResponse struct {
	Result       CalculationResult `json:"result"`
	CalculatedAt time.Time         `json:"calculated_at"`
	Message      string            `json:"message"`
}

// BatchCalculationRequest for calculating multiple students at once
type BatchCalculationRequest struct {
	Students []struct {
		StudentID  int                 `json:"student_id"`
		Components []WeightedComponent `json:"components"`
	} `json:"students" validate:"required,min=1,max=100,dive"`
	Curve *GradeCurve `json:"curve,omitempty"`
}

// BatchCalculationResponse contains results for multiple students
type BatchCalculationResponse struct {
	Results []struct {
		StudentID int               `json:"student_id"`
		Result    CalculationResult `json:"result"`
		Error     string            `json:"error,omitempty"`
	} `json:"results"`
	CalculatedAt   time.Time `json:"calculated_at"`
	TotalProcessed int       `json:"total_processed"`
	TotalErrors    int       `json:"total_errors"`
	ProcessingTime string    `json:"processing_time"`
}

// Validate validates the calculation request
func (req *CalculationRequest) Validate() error {
	if len(req.Components) == 0 {
		return errors.New("at least one component is required")
	}

	totalWeight := 0.0
	for _, comp := range req.Components {
		if comp.Weight < 0 || comp.Weight > 1 {
			return errors.New("component weight must be between 0 and 1")
		}
		if comp.Score < 0 || comp.Score > 100 {
			return errors.New("component score must be between 0 and 100")
		}
		totalWeight += comp.Weight
	}

	// Allow small floating point tolerance
	if totalWeight < 0.99 || totalWeight > 1.01 {
		return errors.New("total weight of all components must equal 1.0")
	}

	if req.Curve != nil && req.Curve.Enabled {
		if req.Curve.Adjustment < -50 || req.Curve.Adjustment > 50 {
			return errors.New("curve adjustment must be between -50 and 50 percent")
		}
	}

	return nil
}

// CalculateWeightedAverage calculates the weighted average from components
func CalculateWeightedAverage(components []WeightedComponent) float64 {
	total := 0.0
	for _, comp := range components {
		total += comp.Score * comp.Weight
	}
	return total
}

// ApplyCurve applies grade curve adjustment to a score
func ApplyCurve(score float64, curve *GradeCurve) float64 {
	if curve == nil || !curve.Enabled {
		return score
	}
	adjusted := score + curve.Adjustment
	// Clamp to 0-100 range
	if adjusted > 100 {
		return 100
	}
	if adjusted < 0 {
		return 0
	}
	return adjusted
}

// ConvertToLetterGrade converts numeric score to letter grade using Figma design scale
func ConvertToLetterGrade(numericScore float64) LetterGrade {
	for _, grade := range GradeScale {
		if numericScore >= grade.MinPoints && numericScore <= grade.MaxPoints {
			return grade
		}
	}
	// Default to F if out of range (shouldn't happen with proper validation)
	return GradeScale[len(GradeScale)-1]
}

// CalculateGrade performs complete grade calculation with all steps
func CalculateGrade(req CalculationRequest) (CalculationResult, error) {
	if err := req.Validate(); err != nil {
		return CalculationResult{}, err
	}

	// Step 1: Calculate weighted average
	numericGrade := CalculateWeightedAverage(req.Components)

	// Step 2: Apply curve if specified
	curveApplied := false
	curveAdjustment := 0.0
	if req.Curve != nil && req.Curve.Enabled {
		originalGrade := numericGrade
		numericGrade = ApplyCurve(numericGrade, req.Curve)
		curveApplied = true
		curveAdjustment = numericGrade - originalGrade
	}

	// Step 3: Convert to letter grade
	letterGrade := ConvertToLetterGrade(numericGrade)

	return CalculationResult{
		NumericGrade:    numericGrade,
		LetterGrade:     letterGrade,
		CurveApplied:    curveApplied,
		CurveAdjustment: curveAdjustment,
	}, nil
}
