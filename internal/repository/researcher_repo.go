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

type ResearcherRepo struct {
	Client *firestore.Client
}

// 🟢 Login with password verification
func (r *ResearcherRepo) Login(email string, password string) (*models.ResearcherInfo, error) {
	ctx := context.Background()

	// Query by email field instead of using email as document ID
	docs, err := r.Client.Collection("researchers").Where("researcher_email", "==", email).Limit(1).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("researcher not found")
	}

	var researcher models.ResearcherInfo
	docs[0].DataTo(&researcher)

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(researcher.ResearcherPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &researcher, nil
}

func NewResearcherRepo(client *firestore.Client) *ResearcherRepo {
	return &ResearcherRepo{Client: client}
}

// 🟢 GetResearcherAll - fetch all researchers
func (r *ResearcherRepo) GetResearcherAll() ([]models.ResearcherInfo, error) {
	ctx := context.Background()
	docs, err := r.Client.Collection("researchers").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var researchers []models.ResearcherInfo
	for _, doc := range docs {
		var researcher models.ResearcherInfo
		doc.DataTo(&researcher)
		researchers = append(researchers, researcher)
	}
	return researchers, nil
}

// 🟢 GetResearcherByID - fetch one researcher by ID (field query - eventual consistency)
func (r *ResearcherRepo) GetResearcherByID(researcherID string) (*models.ResearcherInfo, error) {
	ctx := context.Background()
	// Try querying by researcher_id field first, then by id field (frontend format)
	docs, err := r.Client.Collection("researchers").Where("researcher_id", "==", researcherID).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(docs) == 0 {
		// Fallback to query by id field (frontend format)
		docs, err = r.Client.Collection("researchers").Where("id", "==", researcherID).Limit(1).Documents(ctx).GetAll()
		if err != nil {
			return nil, err
		}
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf("researcher not found")
	}

	// Manually map fields to support both frontend and model formats
	docData := docs[0].Data()
	researcher := models.ResearcherInfo{}

	// Map fields - support both frontend format and model format
	if val, ok := docData["id"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherID = str
		}
	} else if val, ok := docData["researcher_id"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherID = str
		}
	}

	if val, ok := docData["admin_id"]; ok {
		if str, ok := val.(string); ok {
			researcher.AdminID = str
		}
	}

	if val, ok := docData["prefix"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPrefix = str
		}
	} else if val, ok := docData["researcher_prefix"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPrefix = str
		}
	}

	if val, ok := docData["academic_position"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherAcademicPosition = str
		}
	} else if val, ok := docData["researcher_academic_position"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherAcademicPosition = str
		}
	}

	if val, ok := docData["first_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherFirstName = str
		}
	} else if val, ok := docData["researcher_first_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherFirstName = str
		}
	}

	if val, ok := docData["last_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherLastName = str
		}
	} else if val, ok := docData["researcher_last_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherLastName = str
		}
	}

	if val, ok := docData["department"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherDepartment = str
		}
	} else if val, ok := docData["researcher_department"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherDepartment = str
		}
	}

	if val, ok := docData["phone_number"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPhoneNumber = str
		}
	} else if val, ok := docData["researcher_phone_number"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPhoneNumber = str
		}
	}

	if val, ok := docData["email"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherEmail = str
		}
	} else if val, ok := docData["researcher_email"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherEmail = str
		}
	}

	if val, ok := docData["researcher_password"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPassword = str
		}
	}

	if val, ok := docData["created_at"]; ok {
		if t, ok := val.(time.Time); ok {
			researcher.CreatedAt = t
		}
	}

	if val, ok := docData["updated_at"]; ok {
		if t, ok := val.(time.Time); ok {
			researcher.UpdatedAt = t
		}
	}

	return &researcher, nil
}

// 🟢 GetResearcherByIDDirect - fetch one researcher by ID using document ID lookup (immediate consistency)
func (r *ResearcherRepo) GetResearcherByIDDirect(researcherID string) (*models.ResearcherInfo, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("researchers").Doc(researcherID).Get(ctx)
	if err != nil {
		return nil, err
	}

	// Manually map fields because Firestore uses frontend format (id, first_name)
	// but struct expects model format (researcher_id, researcher_first_name)
	docData := doc.Data()
	researcher := models.ResearcherInfo{}

	// Map fields - support both frontend format and model format
	if val, ok := docData["id"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherID = str
		}
	} else if val, ok := docData["researcher_id"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherID = str
		}
	}

	if val, ok := docData["admin_id"]; ok {
		if str, ok := val.(string); ok {
			researcher.AdminID = str
		}
	}

	if val, ok := docData["prefix"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPrefix = str
		}
	} else if val, ok := docData["researcher_prefix"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPrefix = str
		}
	}

	if val, ok := docData["academic_position"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherAcademicPosition = str
		}
	} else if val, ok := docData["researcher_academic_position"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherAcademicPosition = str
		}
	}

	if val, ok := docData["first_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherFirstName = str
		}
	} else if val, ok := docData["researcher_first_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherFirstName = str
		}
	}

	if val, ok := docData["last_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherLastName = str
		}
	} else if val, ok := docData["researcher_last_name"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherLastName = str
		}
	}

	if val, ok := docData["department"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherDepartment = str
		}
	} else if val, ok := docData["researcher_department"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherDepartment = str
		}
	}

	if val, ok := docData["phone_number"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPhoneNumber = str
		}
	} else if val, ok := docData["researcher_phone_number"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPhoneNumber = str
		}
	}

	if val, ok := docData["email"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherEmail = str
		}
	} else if val, ok := docData["researcher_email"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherEmail = str
		}
	}

	if val, ok := docData["researcher_password"]; ok {
		if str, ok := val.(string); ok {
			researcher.ResearcherPassword = str
		}
	}

	if val, ok := docData["created_at"]; ok {
		if t, ok := val.(time.Time); ok {
			researcher.CreatedAt = t
		}
	}

	if val, ok := docData["updated_at"]; ok {
		if t, ok := val.(time.Time); ok {
			researcher.UpdatedAt = t
		}
	}

	return &researcher, nil
}

// 🟢 GetResearcherByCaseID
func (r *ResearcherRepo) GetResearcherByCaseID(caseID string) (*models.ResearcherInfo, error) {
	ctx := context.Background()
	doc, err := r.Client.Collection("researchers").Where("case_id", "==", caseID).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var researcher models.ResearcherInfo
	doc[0].DataTo(&researcher)
	return &researcher, nil
}

// 🟢 CreateResearcher - auto-generate ResearcherID and create new record
func (r *ResearcherRepo) CreateResearcher(researcher *models.ResearcherInfo) error {
	ctx := context.Background()

	// find last ID to generate next
	docs, err := r.Client.Collection("researchers").OrderBy("researcher_id", firestore.Desc).Limit(1).Documents(ctx).GetAll()
	nextID := "RS-00001"
	if err == nil && len(docs) > 0 {
		lastID := docs[0].Data()["researcher_id"].(string)
		numStr := strings.TrimPrefix(lastID, "RS-")
		if n, err := strconv.Atoi(numStr); err == nil {
			nextID = fmt.Sprintf("RS-%05d", n+1)
		}
	}

	researcher.ResearcherID = nextID
	now := time.Now()
	researcher.CreatedAt = now
	researcher.UpdatedAt = now

	_, err = r.Client.Collection("researchers").Doc(researcher.ResearcherID).Set(ctx, researcher)
	return err
}

// 🟢 UpdateResearcherByID - partial update with UpdatedAt
func (r *ResearcherRepo) UpdateResearcherByID(researcherID string, data *models.ResearcherInfo) error {
	ctx := context.Background()
	// Note: For researchers, document ID = researcher_id
	data.UpdatedAt = time.Now()

	// Convert struct to map for Firestore update (using MergeAll like other repos)
	// Update BOTH frontend format and model format fields to keep them in sync
	updateMap := map[string]interface{}{
		// Model format fields
		"researcher_id":                data.ResearcherID,
		"admin_id":                     data.AdminID,
		"researcher_prefix":            data.ResearcherPrefix,
		"researcher_academic_position": data.ResearcherAcademicPosition,
		"researcher_first_name":        data.ResearcherFirstName,
		"researcher_last_name":         data.ResearcherLastName,
		"researcher_department":        data.ResearcherDepartment,
		"researcher_phone_number":      data.ResearcherPhoneNumber,
		"researcher_email":             data.ResearcherEmail,
		"researcher_password":          data.ResearcherPassword,
		"created_at":                   data.CreatedAt,
		"updated_at":                   data.UpdatedAt,
		// Frontend format fields (keep in sync with model format)
		"id":                data.ResearcherID,
		"prefix":            data.ResearcherPrefix,
		"academic_position": data.ResearcherAcademicPosition,
		"first_name":        data.ResearcherFirstName,
		"last_name":         data.ResearcherLastName,
		"department":        data.ResearcherDepartment,
		"phone_number":      data.ResearcherPhoneNumber,
		"email":             data.ResearcherEmail,
	}

	docRef := r.Client.Collection("researchers").Doc(researcherID)
	_, err := docRef.Set(ctx, updateMap, firestore.MergeAll)
	return err
}
