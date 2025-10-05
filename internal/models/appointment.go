package models

import "time"

type Appointment struct {
	AppointmentID string    `json:"appointment_id" firestore:"appointment_id"`
	CaseID        string    `json:"case_id" firestore:"case_id"`
	Date          time.Time `json:"date" firestore:"date"`
	Status        string    `json:"status" firestore:"status"`
	Location      string    `json:"location" firestore:"location"`
	Note          string    `json:"note" firestore:"note"`
	Summary       string    `json:"summary" firestore:"summary"`
	CreatedAt     time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" firestore:"updated_at"`
}
