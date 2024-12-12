package models

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
)

type FormSchema struct {
	ID     int             `json:"id"`
	Schema json.RawMessage `json:"schema"`
}

type JSONSchemaProperty struct {
	Type        string                        `json:"type"`
	Description string                        `json:"description"`
	Format      string                        `json:"format,omitempty"`
	MinLength   int                           `json:"minLength,omitempty"`
	MaxLength   int                           `json:"maxLength,omitempty"`
	Pattern     string                        `json:"pattern,omitempty"`
	Enum        []string                      `json:"enum,omitempty"`
	Properties  map[string]JSONSchemaProperty `json:"properties,omitempty"`
	Items       *JSONSchemaProperty          `json:"items,omitempty"`
	Required    []string                      `json:"required,omitempty"`
	Minimum     *float64                      `json:"minimum,omitempty"`
	Maximum     *float64                      `json:"maximum,omitempty"`
}

type JSONSchema struct {
	Title      string                        `json:"title"`
	Properties map[string]JSONSchemaProperty `json:"properties"`
}

