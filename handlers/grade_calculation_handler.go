package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"techwave/models"
)

// GradeCalculationHandler handles grade calculation operations
type GradeCalculationHandler struct{}

// NewGradeCalculationHandler creates a new instance
func NewGradeCalculationHandler() *GradeCalculationHandler {
	return &GradeCalculationHandler{}
}

// CalculateGrade handles POST /api/grades/calculate
func (h *GradeCalculationHandler) CalculateGrade(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var req models.CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Invalid request body",
			"message": err.Error(),
			"alert":   "Error: Invalid Input", // Figma alert state
			"color":   "#DC2626",              // Figma error color
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Validation failed",
			"message": err.Error(),
			"alert":   "Error: Invalid Input",
			"color":   "#DC2626",
		})
		return
	}

	// Calculate grade
	result, err := models.CalculateGrade(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Calculation failed",
			"message": err.Error(),
			"alert":   "Error: Invalid Input",
			"color":   "#DC2626",
		})
		return
	}

	// Prepare response with alert state based on grade
	message := "Success: Grade Calculated"
	alertColor := "#22C55E" // Success color from Figma
	if result.NumericGrade < 70 {
		message = "Warning: Low Score"
		alertColor = "#FBBF24" // Warning color from Figma
	}

	response := models.CalculationResponse{
		Result:       result,
		CalculatedAt: time.Now(),
		Message:      message,
	}

	// Add processing time header for performance tracking
	processingTime := time.Since(startTime)
	w.Header().Set("X-Processing-Time", processingTime.String())
	w.Header().Set("X-Alert-Color", alertColor)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CalculateBatchGrades handles POST /api/grades/calculate/batch
func (h *GradeCalculationHandler) CalculateBatchGrades(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var batchReq models.BatchCalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&batchReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Invalid request body",
			"message": err.Error(),
			"alert":   "Error: Invalid Input",
			"color":   "#DC2626",
		})
		return
	}

	// Validate max 100 students for performance requirement
	if len(batchReq.Students) > 100 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Too many students",
			"message": "Maximum 100 students allowed per batch",
			"alert":   "Error: Invalid Input",
			"color":   "#DC2626",
		})
		return
	}

	if len(batchReq.Students) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "Empty batch",
			"message": "At least one student required",
			"alert":   "Error: Invalid Input",
			"color":   "#DC2626",
		})
		return
	}

	// Process all students
	response := models.BatchCalculationResponse{
		CalculatedAt: time.Now(),
	}

	for _, student := range batchReq.Students {
		req := models.CalculationRequest{
			Components: student.Components,
			Curve:      batchReq.Curve,
		}

		result, err := models.CalculateGrade(req)
		resultEntry := struct {
			StudentID int                      `json:"student_id"`
			Result    models.CalculationResult `json:"result"`
			Error     string                   `json:"error,omitempty"`
		}{
			StudentID: student.StudentID,
			Result:    result,
		}

		if err != nil {
			resultEntry.Error = err.Error()
			response.TotalErrors++
		}

		response.Results = append(response.Results, resultEntry)
		response.TotalProcessed++
	}

	processingTime := time.Since(startTime)
	response.ProcessingTime = processingTime.String()

	// Performance check: log warning if exceeds 200ms requirement
	if processingTime.Milliseconds() > 200 {
		w.Header().Set("X-Performance-Warning", "Exceeded 200ms target")
	}

	w.Header().Set("X-Processing-Time", processingTime.String())
	w.Header().Set("Content-Type", "application/json")

	// Use info alert color for successful batch processing
	message := "Info: Grade Updated"
	if response.TotalErrors > 0 {
		message = "Warning: Low Score"
	}
	w.Header().Set("X-Alert-Message", message)
	w.Header().Set("X-Alert-Color", "#3B82F6") // Info color from Figma

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetGradeScale handles GET /api/grades/scale
func (h *GradeCalculationHandler) GetGradeScale(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"scale":       models.GradeScale,
		"description": "Grade Calculation - Design System",
		"source":      "Figma Design System",
		"colors": map[string]string{
			"success": "#22C55E",
			"warning": "#FBBF24",
			"error":   "#DC2626",
			"info":    "#3B82F6",
		},
	})
}
