package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"trl-research-backend/internal/models"
)

type AppointmentRepo struct {
	Client *firestore.Client
}

func NewAppointmentRepo(client *firestore.Client) *AppointmentRepo {
	return &AppointmentRepo{Client: client}
}

// 🟢 GetAppointmentAll
func (r *AppointmentRepo) GetAppointmentAll() ([]models.Appointment, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("appointments").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var appointments []models.Appointment
	for _, doc := range docs {
		var ap models.Appointment
		doc.DataTo(&ap)
		appointments = append(appointments, ap)
	}
	return appointments, nil
}

// 🟢 GetAppointmentByID
func (r *AppointmentRepo) GetAppointmentByID(appointmentID string) (*models.Appointment, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("appointments").Doc(appointmentID).Get(ctx)
	if err != nil {
		return nil, err
	}

	var ap models.Appointment
	doc.DataTo(&ap)
	return &ap, nil
}

// 🟢 GetAppointmentByCaseID
func (r *AppointmentRepo) GetAppointmentByCaseID(caseID string) (*models.Appointment, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("appointments").Where("case_id", "==", caseID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var ap models.Appointment
	doc[0].DataTo(&ap)
	return &ap, nil
}

// 🟢 CreateAppointment - auto generate ID AP-00001
func (r *AppointmentRepo) CreateAppointment(ap *models.Appointment) error {
	ctx := context.Background()

	docs, err := r.Client.Collection("appointments").OrderBy("appointment_id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "AP-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["appointment_id"].(string)
		numStr := strings.TrimPrefix(lastID, "AP-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("AP-%05d", n+1)
		}
	}

	ap.AppointmentID = nextID
	now := time.Now()
	ap.CreatedAt = now
	ap.UpdatedAt = now

	_, err = r.Client.Collection("appointments").Doc(ap.AppointmentID).Set(ctx, ap)
	return err
}

// 🟢 UpdateAppointmentByID
func (r *AppointmentRepo) UpdateAppointmentByID(appointmentID string, data map[string]interface{}) error {
	ctx := context.Background()
	data["updated_at"] = time.Now()
	_, err := r.Client.Collection("appointments").Doc(appointmentID).Set(ctx, data, firestore.MergeAll)
	return err
}
