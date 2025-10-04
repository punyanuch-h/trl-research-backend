package models

type TrlQuestion struct {
	AssessmentID string `firestore:"assessment_id"` // ref ไปที่ assessment_trl
	QuestionCode string `firestore:"question_code"` // เช่น "RQ1", "RQ2", "CQ1"
	AnswerType   string `firestore:"answer_type"`   // "bool" หรือ "text"
	AnswerValue  string `firestore:"answer_value"`  // เก็บค่าจริงเป็น string เช่น "true"/"false" หรือข้อความ
}
