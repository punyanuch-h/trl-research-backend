package models

import (
	"time"
)

type AssessmentTRLPart1 struct {
	ID           string    `json:"id" firestore:"id"`
	AssessmentID string    `json:"assessment_id" firestore:"assessment_id"`
	RQ1Answer    bool      `json:"rq1_answer" firestore:"rq1_answer"`
	RQ2Answer    bool      `json:"rq2_answer" firestore:"rq2_answer"`
	RQ3Answer    bool      `json:"rq3_answer" firestore:"rq3_answer"`
	RQ4Answer    bool      `json:"rq4_answer" firestore:"rq4_answer"`
	RQ5Answer    bool      `json:"rq5_answer" firestore:"rq5_answer"`
	RQ6Answer    bool      `json:"rq6_answer" firestore:"rq6_answer"`
	RQ7Answer    bool      `json:"rq7_answer" firestore:"rq7_answer"`
	CreatedAt    time.Time `json:"created_at" firestore:"created_at"`
}