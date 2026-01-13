package models

import "time"

type AssessmentTrl struct {
	ID             string `json:"id" firestore:"id" form:"id"`
	CaseID         string `json:"case_id" firestore:"case_id" form:"case_id"`
	TrlLevelResult int    `json:"trl_level_result" firestore:"trl_level_result" form:"trl_level_result"`

	Rq1Answer      bool     `json:"rq1_answer" firestore:"rq1_answer" form:"rq1_answer"`
	Rq1Attachments []string `json:"rq1_attachment" firestore:"rq1_attachment"`
	Rq2Answer      bool     `json:"rq2_answer" firestore:"rq2_answer" form:"rq2_answer"`
	Rq2Attachments []string `json:"rq2_attachment" firestore:"rq2_attachment"`
	Rq3Answer      bool     `json:"rq3_answer" firestore:"rq3_answer" form:"rq3_answer"`
	Rq3Attachments []string `json:"rq3_attachment" firestore:"rq3_attachment"`
	Rq4Answer      bool     `json:"rq4_answer" firestore:"rq4_answer" form:"rq4_answer"`
	Rq4Attachments []string `json:"rq4_attachment" firestore:"rq4_attachment"`
	Rq5Answer      bool     `json:"rq5_answer" firestore:"rq5_answer" form:"rq5_answer"`
	Rq5Attachments []string `json:"rq5_attachment" firestore:"rq5_attachment"`
	Rq6Answer      bool     `json:"rq6_answer" firestore:"rq6_answer" form:"rq6_answer"`
	Rq6Attachments []string `json:"rq6_attachment" firestore:"rq6_attachment"`
	Rq7Answer      bool     `json:"rq7_answer" firestore:"rq7_answer" form:"rq7_answer"`
	Rq7Attachments []string `json:"rq7_attachment" firestore:"rq7_attachment"`

	Cq1Answer      []string `json:"cq1_answer" firestore:"cq1_answer" form:"cq1_answer"`
	Cq1Attachments []string `json:"cq1_attachment" firestore:"cq1_attachment"`
	Cq2Answer      []string `json:"cq2_answer" firestore:"cq2_answer" form:"cq2_answer"`
	Cq2Attachments []string `json:"cq2_attachment" firestore:"cq2_attachment"`
	Cq3Answer      []string `json:"cq3_answer" firestore:"cq3_answer" form:"cq3_answer"`
	Cq3Attachments []string `json:"cq3_attachment" firestore:"cq3_attachment"`
	Cq4Answer      []string `json:"cq4_answer" firestore:"cq4_answer" form:"cq4_answer"`
	Cq4Attachments []string `json:"cq4_attachment" firestore:"cq4_attachment"`
	Cq5Answer      []string `json:"cq5_answer" firestore:"cq5_answer" form:"cq5_answer"`
	Cq5Attachments []string `json:"cq5_attachment" firestore:"cq5_attachment"`
	Cq6Answer      []string `json:"cq6_answer" firestore:"cq6_answer" form:"cq6_answer"`
	Cq6Attachments []string `json:"cq6_attachment" firestore:"cq6_attachment"`
	Cq7Answer      []string `json:"cq7_answer" firestore:"cq7_answer" form:"cq7_answer"`
	Cq7Attachments []string `json:"cq7_attachment" firestore:"cq7_attachment"`
	Cq8Answer      []string `json:"cq8_answer" firestore:"cq8_answer" form:"cq8_answer"`
	Cq8Attachments []string `json:"cq8_attachment" firestore:"cq8_attachment"`
	Cq9Answer      []string `json:"cq9_answer" firestore:"cq9_answer" form:"cq9_answer"`
	Cq9Attachments []string `json:"cq9_attachment" firestore:"cq9_attachment"`

	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}
