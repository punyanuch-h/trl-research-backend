package models

import (
	"time"
)

type AssessmentTRLPart2 struct {
	ID           string    `json:"id" firestore:"id"`
	AssessmentID string    `json:"assessment_id" firestore:"assessment_id"`
	CQ1Answer    string    `json:"cq1_answer" firestore:"cq1_answer"`
	CQ2Answer    string    `json:"cq2_answer" firestore:"cq2_answer"`
	CQ3Answer    string    `json:"cq3_answer" firestore:"cq3_answer"`
	CQ4Answer    string    `json:"cq4_answer" firestore:"cq4_answer"`
	CQ5Answer    string    `json:"cq5_answer" firestore:"cq5_answer"`
	CQ6Answer    string    `json:"cq6_answer" firestore:"cq6_answer"`
	CQ7Answer    string    `json:"cq7_answer" firestore:"cq7_answer"`
	CQ8Answer    string    `json:"cq8_answer" firestore:"cq8_answer"`
	CQ9Answer    string    `json:"cq9_answer" firestore:"cq9_answer"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
}
