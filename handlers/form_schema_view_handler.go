package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"htmx-go/models"
	"htmx-go/repositories"
	"net/http"
	"strconv"
	"strings"
)

type formSchemaViewHandler struct {
	formSchemaRepository     repositories.FormSchemaRepository
	formSubmissionRepository repositories.FormSubmissionRepository
}

func NewFormSchemaViewHandler(rootRoute *echo.Echo, repository repositories.FormSchemaRepository, submissionRepo repositories.FormSubmissionRepository) {
	handler := &formSchemaViewHandler{
		formSchemaRepository:     repository,
		formSubmissionRepository: submissionRepo,
	}

	rootRoute.GET("/", handler.RenderFormSchemaList)
	rootRoute.GET("/search", handler.SearchFormSchemas)
	rootRoute.GET("/create", handler.RenderCreateForm)
	rootRoute.POST("/create", handler.CreateForm)
	rootRoute.GET("/:id", handler.RenderFormSchemaById)
	rootRoute.POST("/:id/submit", handler.SubmitForm)
}

func (h *formSchemaViewHandler) RenderFormSchemaById(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Render(http.StatusOK, "form-not-found.html", nil)
	}

	formSchema, err := h.formSchemaRepository.GetSchemaById(id)
	if err != nil {
		return c.Render(http.StatusOK, "form-not-found.html", nil)
	}

	parsedSchema := models.JSONSchema{}
	if err := json.Unmarshal(formSchema.Schema, &parsedSchema); err != nil {
		return c.Render(http.StatusOK, "form-not-found.html", nil)
	}

	return c.Render(http.StatusOK, "form.html", map[string]interface{}{
		"ID":     formSchema.ID,
		"Schema": parsedSchema,
	})
}

func (h *formSchemaViewHandler) SearchFormSchemas(c echo.Context) error {
	search := c.QueryParam("search")
	formSchemas, err := h.formSchemaRepository.GetSchema()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Error loading forms",
		})
	}

	var filteredSchemas []map[string]interface{}
	for _, schema := range formSchemas {
		var parsedSchema models.JSONSchema
		if err := json.Unmarshal(schema.Schema, &parsedSchema); err != nil {
			continue
		}

		// Case-insensitive search
		if search == "" || strings.Contains(
			strings.ToLower(parsedSchema.Title),
			strings.ToLower(search)) {
			filteredSchemas = append(filteredSchemas, map[string]interface{}{
				"ID":    schema.ID,
				"Title": parsedSchema.Title,
			})
		}
	}
	return c.Render(http.StatusOK, "partials/form-table-partial.html", map[string]interface{}{
		"Schemas": filteredSchemas,
	})
}

func (h *formSchemaViewHandler) RenderFormSchemaList(c echo.Context) error {
	formSchemas, err := h.formSchemaRepository.GetSchema()
	if err != nil {
		return c.Render(http.StatusOK, "error.html", nil)
	}

	var schemaList []map[string]interface{}
	for _, schema := range formSchemas {
		var parsedSchema models.JSONSchema
		if err := json.Unmarshal(schema.Schema, &parsedSchema); err != nil {
			continue
		}
		schemaList = append(schemaList, map[string]interface{}{
			"ID":    schema.ID,
			"Title": parsedSchema.Title,
		})
	}

	return c.Render(http.StatusOK, "form-list.html", map[string]interface{}{
		"Schemas": schemaList,
	})
}
func (h *formSchemaViewHandler) RenderCreateForm(c echo.Context) error {
	return c.Render(http.StatusOK, "create-form.html", nil)
}

func (h *formSchemaViewHandler) CreateForm(c echo.Context) error {
	title := c.FormValue("title")
	schemaStr := c.FormValue("schema")

	// Validate the schema is valid JSON
	var schemaData models.JSONSchema
	if err := json.Unmarshal([]byte(schemaStr), &schemaData); err != nil {
		return c.Render(http.StatusBadRequest, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Invalid JSON schema format",
		})
	}

	// Set the title from the form
	schemaData.Title = title

	// Convert back to JSON string
	finalSchema, err := json.Marshal(schemaData)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Error processing schema",
		})
	}

	// Save the schema
	_, err = h.formSchemaRepository.SaveSchema(string(finalSchema))
	if err != nil {
		return c.Render(http.StatusInternalServerError, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Error saving form",
		})
	}

	return c.Render(http.StatusOK, "partials/alerts/success-redirect.html", map[string]interface{}{
		"Message":     "Form created successfully! Redirecting...",
		"RedirectURL": "/",
		"Delay":       1500,
	})
}

func (h *formSchemaViewHandler) SubmitForm(c echo.Context) error {
	paramId := c.Param("id")
	formSchemaId, err := strconv.Atoi(paramId)
	if err != nil {
		return c.Render(http.StatusBadRequest, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Invalid form ID",
		})
	}

	// Parse the form data into a map
	formData := make(map[string]interface{})
	if err := c.Bind(&formData); err != nil {
		return c.Render(http.StatusBadRequest, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Invalid form data",
		})
	}

	// Convert the form data to JSON
	submittedValues, err := json.Marshal(formData)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Error processing form data",
		})
	}

	// Save the submission
	err = h.formSubmissionRepository.SaveSubmission(formSchemaId, submittedValues)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "partials/alerts/error.html", map[string]interface{}{
			"Message": "Error saving submission",
		})
	}

	return c.Render(http.StatusOK, "partials/alerts/success.html", map[string]interface{}{
		"Message": "Form submitted successfully!",
	})
}
