package repositories

import (
	"database/sql"
	"encoding/json"
	"htmx-go/models"
)

type FormSchemaRepository interface {
	GetSchema() ([]models.FormSchema, error)
	GetSchemaById(id int) (*models.FormSchema, error)
	SaveSchema(schemaJSON string) (int64, error)
}

type formSchemaRepository struct {
	db *sql.DB
}

func NewFormSchemaRepository(db *sql.DB) FormSchemaRepository {
	return &formSchemaRepository{
		db,
	}
}

func (r *formSchemaRepository) GetSchema() ([]models.FormSchema, error) {
	rows, err := r.db.Query("SELECT id, schema FROM form_schemas ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var formSchemas []models.FormSchema
	for rows.Next() {
		var formSchema models.FormSchema
		var schemaStr string
		if err := rows.Scan(&formSchema.ID, &schemaStr); err != nil {
			return formSchemas, err
		}
		formSchema.Schema = json.RawMessage(schemaStr)
		formSchemas = append(formSchemas, formSchema)
	}

	if err = rows.Err(); err != nil {
		return formSchemas, err
	}

	return formSchemas, nil
}

func (r *formSchemaRepository) GetSchemaById(id int) (*models.FormSchema, error) {
	var schema models.FormSchema
	var schemaStr string
	err := r.db.QueryRow("SELECT id, schema FROM form_schemas WHERE id = ?", id).Scan(&schema.ID, &schemaStr)
	if err != nil {
		return nil, err
	}
	schema.Schema = json.RawMessage(schemaStr)
	return &schema, nil
}

func (r *formSchemaRepository) SaveSchema(schemaJSON string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO form_schemas (schema) VALUES (?)", schemaJSON)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
