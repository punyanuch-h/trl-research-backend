package repository

import (
	"trl-research-backend/internal/database"
	"trl-research-backend/internal/models"
)

type AppointmentRepo struct{}

func (r *AppointmentRepo) CreateOrUpdate(data *models.Appointment) error {
	_, err := database.FirestoreClient.Collection("appointment").
		Doc(data.ID).
		Set(database.Ctx, data)
	return err
}

func (r *AppointmentRepo) GetByID(id string) (*models.Appointment, error) {
	doc, err := database.FirestoreClient.Collection("appointment").
		Doc(id).Get(database.Ctx)
	if err != nil {
		return nil, err
	}
	var obj models.Appointment
	_ = doc.DataTo(&obj)
	return &obj, nil
}

func (r *AppointmentRepo) Delete(id string) error {
	_, err := database.FirestoreClient.Collection("appointment").
		Doc(id).Delete(database.Ctx)
	return err
}
