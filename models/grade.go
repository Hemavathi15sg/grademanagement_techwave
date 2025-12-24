package models

type Grade struct {
    ID        int    `json:"id"`
    StudentID int    `json:"student_id"`
    Value     string `json:"value"`
    Subject   string `json:"subject"`
}