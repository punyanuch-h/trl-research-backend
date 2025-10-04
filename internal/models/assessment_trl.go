package models

import (
	"time"
)

type AssessmentTRL struct {
	ID             string    `json:"id" firestore:"id"`
	CaseID         string    `json:"case_id" firestore:"case_id"`
	TrlLevelResult int       `json:"trl_level_result" firestore:"trl_level_result"`
	Part1          []AssessmentTRLPart1 `json:"part1" firestore:"part1"`
	Part2          []AssessmentTRLPart2 `json:"part2" firestore:"part2"`
	CreatedAt      time.Time `json:"created_at" firestore:"created_at"`
}
