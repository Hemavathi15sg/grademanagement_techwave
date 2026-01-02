package tests

import (
	"testing"

	"techwave/models"
)

func TestCalculateWeightedAverage(t *testing.T) {
	tests := []struct {
		name       string
		components []models.WeightedComponent
		expected   float64
	}{
		{
			name: "equal weights",
			components: []models.WeightedComponent{
				{Name: "assignments", Score: 80, Weight: 0.33},
				{Name: "exams", Score: 90, Weight: 0.34},
				{Name: "attendance", Score: 100, Weight: 0.33},
			},
			expected: 90.0,
		},
		{
			name: "typical weights",
			components: []models.WeightedComponent{
				{Name: "assignments", Score: 85, Weight: 0.30},
				{Name: "midterm", Score: 90, Weight: 0.30},
				{Name: "final", Score: 95, Weight: 0.40},
			},
			expected: 90.5,
		},
		{
			name: "single component",
			components: []models.WeightedComponent{
				{Name: "final", Score: 88, Weight: 1.0},
			},
			expected: 88.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.CalculateWeightedAverage(tt.components)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestApplyCurve(t *testing.T) {
	tests := []struct {
		name     string
		score    float64
		curve    *models.GradeCurve
		expected float64
	}{
		{
			name:     "no curve",
			score:    85,
			curve:    nil,
			expected: 85,
		},
		{
			name:  "curve disabled",
			score: 85,
			curve: &models.GradeCurve{Enabled: false, Adjustment: 5},
			expected: 85,
		},
		{
			name:  "positive adjustment",
			score: 75,
			curve: &models.GradeCurve{Enabled: true, Adjustment: 10},
			expected: 85,
		},
		{
			name:  "negative adjustment",
			score: 85,
			curve: &models.GradeCurve{Enabled: true, Adjustment: -5},
			expected: 80,
		},
		{
			name:  "capped at 100",
			score: 95,
			curve: &models.GradeCurve{Enabled: true, Adjustment: 10},
			expected: 100,
		},
		{
			name:  "capped at 0",
			score: 5,
			curve: &models.GradeCurve{Enabled: true, Adjustment: -10},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.ApplyCurve(tt.score, tt.curve)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvertToLetterGrade(t *testing.T) {
	tests := []struct {
		name          string
		numericScore  float64
		expectedLetter string
		expectedColor string
	}{
		{name: "A+ grade", numericScore: 98, expectedLetter: "A+", expectedColor: "#22C55E"},
		{name: "A grade", numericScore: 95, expectedLetter: "A", expectedColor: "#22C55E"},
		{name: "A- grade", numericScore: 91, expectedLetter: "A-", expectedColor: "#86EFAC"},
		{name: "B+ grade", numericScore: 88, expectedLetter: "B+", expectedColor: "#3B82F6"},
		{name: "B grade", numericScore: 84, expectedLetter: "B", expectedColor: "#3B82F6"},
		{name: "B- grade", numericScore: 81, expectedLetter: "B-", expectedColor: "#60A5FA"},
		{name: "C+ grade", numericScore: 78, expectedLetter: "C+", expectedColor: "#FBBF24"},
		{name: "C grade", numericScore: 74, expectedLetter: "C", expectedColor: "#FBBF24"},
		{name: "C- grade", numericScore: 71, expectedLetter: "C-", expectedColor: "#FCD34D"},
		{name: "D grade", numericScore: 65, expectedLetter: "D", expectedColor: "#F97316"},
		{name: "F grade", numericScore: 50, expectedLetter: "F", expectedColor: "#DC2626"},
		{name: "edge case 100", numericScore: 100, expectedLetter: "A+", expectedColor: "#22C55E"},
		{name: "edge case 0", numericScore: 0, expectedLetter: "F", expectedColor: "#DC2626"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := models.ConvertToLetterGrade(tt.numericScore)
			if result.Letter != tt.expectedLetter {
				t.Errorf("expected letter %v, got %v", tt.expectedLetter, result.Letter)
			}
			if result.Color != tt.expectedColor {
				t.Errorf("expected color %v, got %v", tt.expectedColor, result.Color)
			}
		})
	}
}

func TestCalculateGrade(t *testing.T) {
	tests := []struct {
		name           string
		request        models.CalculationRequest
		expectedGrade  string
		expectedError  bool
	}{
		{
			name: "excellent grade without curve",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "assignments", Score: 95, Weight: 0.4},
					{Name: "exams", Score: 98, Weight: 0.5},
					{Name: "attendance", Score: 100, Weight: 0.1},
				},
			},
			expectedGrade: "A+",
			expectedError: false,
		},
		{
			name: "grade with positive curve",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "final", Score: 85, Weight: 1.0},
				},
				Curve: &models.GradeCurve{Enabled: true, Adjustment: 5},
			},
			expectedGrade: "A-",
			expectedError: false,
		},
		{
			name: "failing grade",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 45, Weight: 1.0},
				},
			},
			expectedGrade: "F",
			expectedError: false,
		},
		{
			name: "invalid weights sum",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 85, Weight: 0.5},
				},
			},
			expectedError: true,
		},
		{
			name: "score out of range",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 150, Weight: 1.0},
				},
			},
			expectedError: true,
		},
		{
			name: "curve adjustment too large",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 85, Weight: 1.0},
				},
				Curve: &models.GradeCurve{Enabled: true, Adjustment: 60},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := models.CalculateGrade(tt.request)
			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result.LetterGrade.Letter != tt.expectedGrade {
					t.Errorf("expected grade %v, got %v", tt.expectedGrade, result.LetterGrade.Letter)
				}
			}
		})
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name        string
		request     models.CalculationRequest
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid request",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 85, Weight: 1.0},
				},
			},
			expectError: false,
		},
		{
			name:        "empty components",
			request:     models.CalculationRequest{},
			expectError: true,
			errorMsg:    "at least one component is required",
		},
		{
			name: "weight below zero",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 85, Weight: -0.1},
				},
			},
			expectError: true,
			errorMsg:    "component weight must be between 0 and 1",
		},
		{
			name: "weight above one",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 85, Weight: 1.5},
				},
			},
			expectError: true,
			errorMsg:    "component weight must be between 0 and 1",
		},
		{
			name: "score below zero",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: -5, Weight: 1.0},
				},
			},
			expectError: true,
			errorMsg:    "component score must be between 0 and 100",
		},
		{
			name: "score above 100",
			request: models.CalculationRequest{
				Components: []models.WeightedComponent{
					{Name: "exam", Score: 105, Weight: 1.0},
				},
			},
			expectError: true,
			errorMsg:    "component score must be between 0 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkCalculateGrade tests performance requirement (<200ms for 100 students)
func BenchmarkCalculateGrade(b *testing.B) {
	req := models.CalculationRequest{
		Components: []models.WeightedComponent{
			{Name: "assignments", Score: 85, Weight: 0.30},
			{Name: "midterm", Score: 90, Weight: 0.30},
			{Name: "final", Score: 92, Weight: 0.40},
		},
		Curve: &models.GradeCurve{Enabled: true, Adjustment: 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = models.CalculateGrade(req)
	}
}

// BenchmarkBatchCalculation simulates 100 students
func BenchmarkBatchCalculation100Students(b *testing.B) {
	requests := make([]models.CalculationRequest, 100)
	for i := 0; i < 100; i++ {
		requests[i] = models.CalculationRequest{
			Components: []models.WeightedComponent{
				{Name: "assignments", Score: 80 + float64(i%20), Weight: 0.30},
				{Name: "midterm", Score: 75 + float64(i%25), Weight: 0.30},
				{Name: "final", Score: 85 + float64(i%15), Weight: 0.40},
			},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, req := range requests {
			_, _ = models.CalculateGrade(req)
		}
	}
}
