package repositories

import (
	"database/sql"
	"encoding/json"
	"time"
)

type FormSubmission struct {
	ID              int       `json:"id"`
	FormSchemaID    int       `json:"form_schema_id"`
	SubmittedValues string    `json:"submitted_values"`
	CreatedAt       time.Time `json:"created_at"`
}

type FormSubmissionRepository interface {
	SaveSubmission(formSchemaId int, submittedValues json.RawMessage) error
	GetSubmissions(page, pageSize int) ([]FormSubmission, error)
	GetSubmissionsByFormID(formID, page, pageSize int) ([]FormSubmission, error)
}

type formSubmissionRepository struct {
	db *sql.DB
}

func NewFormSubmissionRepository(db *sql.DB) FormSubmissionRepository {
	return &formSubmissionRepository{
		db: db,
	}
}

func (r *formSubmissionRepository) SaveSubmission(formSchemaId int, submittedValues json.RawMessage) error {
	_, err := r.db.Exec(
		"INSERT INTO form_submissions (form_schema_id, submitted_values) VALUES (?, ?)",
		formSchemaId,
		submittedValues,
	)
	return err
}

func (r *formSubmissionRepository) GetSubmissions(page, pageSize int) ([]FormSubmission, error) {
	offset := (page - 1) * pageSize
	rows, err := r.db.Query(`
		SELECT id, form_schema_id, submitted_values, created_at 
		FROM form_submissions 
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?`,
		pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []FormSubmission
	for rows.Next() {
		var submission FormSubmission
		err := rows.Scan(
			&submission.ID,
			&submission.FormSchemaID,
			&submission.SubmittedValues,
			&submission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, nil
}

func (r *formSubmissionRepository) GetSubmissionsByFormID(formID, page, pageSize int) ([]FormSubmission, error) {
	offset := (page - 1) * pageSize
	rows, err := r.db.Query(`
		SELECT id, form_schema_id, submitted_values, created_at 
		FROM form_submissions 
		WHERE form_schema_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?`,
		formID, pageSize, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []FormSubmission
	for rows.Next() {
		var submission FormSubmission
		err := rows.Scan(
			&submission.ID,
			&submission.FormSchemaID,
			&submission.SubmittedValues,
			&submission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, submission)
	}
	return submissions, nil
}
