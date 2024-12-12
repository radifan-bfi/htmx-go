package models

import (
	"encoding/json"
	"time"
)

type FormSubmission struct {
	ID              int             `json:"id"`
	FormSchemaID    int             `json:"form_schema_id"`
	SubmittedValues json.RawMessage `json:"submitted_values"`
	CreatedAt       time.Time       `json:"created_at"`
}
