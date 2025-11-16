package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"trl-research-backend/internal/models"

	"cloud.google.com/go/firestore"
	"golang.org/x/crypto/bcrypt"
)

type AdminRepo struct {
	Client *firestore.Client
}

func NewAdminRepo(client *firestore.Client) *AdminRepo {
	return &AdminRepo{Client: client}
}

// 🟢 Get all admins
func (r *AdminRepo) GetAdminAll() ([]models.AdminInfo, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("admin_info").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var admins []models.AdminInfo
	for _, doc := range docs {
		var admin models.AdminInfo
		doc.DataTo(&admin)
		admins = append(admins, admin)
	}
	return admins, nil
}

// 🟢 Get admin by ID
func (r *AdminRepo) GetAdminByID(adminID string) (*models.AdminInfo, error) {
	ctx := context.Background()

	// Query by admin_id field instead of document ID
	docs, err := r.Client.Collection("admin_info").Where("admin_id", "==", adminID).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("admin not found")
	}

	var admin models.AdminInfo
	docs[0].DataTo(&admin)
	return &admin, nil
}

// 🟢 Get admin by email
func (r *AdminRepo) GetAdminByEmail(email string) (*models.AdminInfo, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("admin_info").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}

	var admin models.AdminInfo
	doc.DataTo(&admin)
	return &admin, nil
}

// 🟢 Create admin (auto-generate AdminID)
func (r *AdminRepo) CreateAdmin(admin *models.AdminInfo) error {
	ctx := context.Background()

	// find last ID to generate next
	docs, err := r.Client.Collection("admin_info").OrderBy("admin_id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "AD-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["admin_id"].(string)
		numStr := strings.TrimPrefix(lastID, "SI-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("AD-%05d", n+1)
		}
	}

	// assign values
	admin.AdminID = nextID
	now := time.Now()
	admin.CreatedAt = now

	// save to Firestore using email as document ID
	_, err = r.Client.Collection("admin_info").Doc(admin.AdminEmail).Set(ctx, admin)
	return err
}

// 🟢 Login with password verification
func (r *AdminRepo) Login(email string, password string) (*models.AdminInfo, error) {
	ctx := context.Background()

	// Query by email field instead of using email as document ID
	docs, err := r.Client.Collection("admin_info").Where("admin_email", "==", email).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("admin not found")
	}

	var admin models.AdminInfo
	docs[0].DataTo(&admin)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(admin.AdminPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &admin, nil
}

// 🟢 Update password
func (r *AdminRepo) UpdatePasswordByEmail(email string, password string) error {
	ctx := context.Background()
	_, err := r.Client.Collection("admin_info").Doc(email).Set(ctx, map[string]interface{}{"admin_password": password}, firestore.MergeAll)
	return err
}

// 🟢 Update admin by ID
func (r *AdminRepo) UpdateAdminByID(id string, data *models.AdminInfo) error {
	ctx := context.Background()
	// update data from request body to this admin
	// Note: Documents are stored with email as document ID, not admin_id
	// So we use data.AdminEmail (which should be set from existingAdmin)
	data.UpdatedAt = time.Now()

	// Convert struct to map for Firestore update (using MergeAll like other repos)
	updateMap := map[string]interface{}{
		"admin_id":                data.AdminID,
		"admin_prefix":            data.AdminPrefix,
		"admin_academic_position": data.AdminAcademicPosition,
		"admin_first_name":        data.AdminFirstName,
		"admin_last_name":         data.AdminLastName,
		"admin_department":        data.AdminDepartment,
		"admin_phone_number":      data.AdminPhoneNumber,
		"admin_email":             data.AdminEmail,
		"admin_password":          data.AdminPassword,
		"case_id":                 data.CaseID,
		"created_at":              data.CreatedAt,
		"updated_at":              data.UpdatedAt,
	}

	docRef := r.Client.Collection("admin_info").Doc(data.AdminEmail)
	_, err := docRef.Set(ctx, updateMap, firestore.MergeAll)
	return err
}

// 🟢 Delete admin
func (r *AdminRepo) DeleteAdmin(email string) error {
	ctx := context.Background()
	_, err := r.Client.Collection("admin_info").Doc(email).Delete(ctx)
	return err
}
