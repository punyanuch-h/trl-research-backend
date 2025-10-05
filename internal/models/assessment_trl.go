package models

import "time"

type AssessmentTrl struct {
	ID              string    `json:"id" firestore:"id"`
	CaseID          string    `json:"case_id" firestore:"case_id"`
	TrlLevelResult  int       `json:"trl_level_result" firestore:"trl_level_result"`
	Rq1Answer       bool      `json:"rq1_answer" firestore:"rq1_answer"`
	Rq2Answer       bool      `json:"rq2_answer" firestore:"rq2_answer"`
	Rq3Answer       bool      `json:"rq3_answer" firestore:"rq3_answer"`
	Rq4Answer       bool      `json:"rq4_answer" firestore:"rq4_answer"`
	Rq5Answer       bool      `json:"rq5_answer" firestore:"rq5_answer"`
	Rq6Answer       bool      `json:"rq6_answer" firestore:"rq6_answer"`
	Rq7Answer       bool      `json:"rq7_answer" firestore:"rq7_answer"`
	Cq1Answer       string    `json:"cq1_answer" firestore:"cq1_answer"`
	Cq2Answer       string    `json:"cq2_answer" firestore:"cq2_answer"`
	Cq3Answer       string    `json:"cq3_answer" firestore:"cq3_answer"`
	Cq4Answer       string    `json:"cq4_answer" firestore:"cq4_answer"`
	Cq5Answer       string    `json:"cq5_answer" firestore:"cq5_answer"`
	Cq6Answer       string    `json:"cq6_answer" firestore:"cq6_answer"`
	Cq7Answer       string    `json:"cq7_answer" firestore:"cq7_answer"`
	Cq8Answer       string    `json:"cq8_answer" firestore:"cq8_answer"`
	Cq9Answer       string    `json:"cq9_answer" firestore:"cq9_answer"`
	CreatedAt       time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" firestore:"updated_at"`
}
