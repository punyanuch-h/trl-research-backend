package models

import "time"

type Supporter struct {
	SupporterID                    string    `json:"supporter_id" firestore:"supporter_id"`
	CaseID                         string    `json:"case_id" firestore:"case_id"`
	SupportResearch                bool      `json:"support_research" firestore:"support_research"`
	SupportVDC                     bool      `json:"support_vdc" firestore:"support_vdc"`
	SupportSiEIC                   bool      `json:"support_sieic" firestore:"support_sieic"`
	NeedProtectIntellectualProperty bool     `json:"need_protect_intellectual_property" firestore:"need_protect_intellectual_property"`
	NeedCoDevelopers               bool      `json:"need_co_developers" firestore:"need_co_developers"`
	NeedActivities                 bool      `json:"need_activities" firestore:"need_activities"`
	NeedTest                       bool      `json:"need_test" firestore:"need_test"`
	NeedCapital                    bool      `json:"need_capital" firestore:"need_capital"`
	NeedPartners                   bool      `json:"need_partners" firestore:"need_partners"`
	NeedGuidelines                 bool      `json:"need_guidelines" firestore:"need_guidelines"`
	NeedCertification              bool      `json:"need_certification" firestore:"need_certification"`
	NeedAccount                    bool      `json:"need_account" firestore:"need_account"`
	Need                           string    `json:"need" firestore:"need"`
	AdditionalDocuments            string    `json:"additional_documents" firestore:"additional_documents"`
	CreatedAt                      time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt                      time.Time `json:"updated_at" firestore:"updated_at"`
}
