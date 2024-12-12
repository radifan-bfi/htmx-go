CREATE TABLE form_submissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    form_schema_id INTEGER NOT NULL,
    submitted_values TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (form_schema_id) REFERENCES form_schemas(id)
);
