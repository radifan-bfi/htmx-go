package handlers

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"htmx-go/models"
	"htmx-go/repositories"
	"net/http"
	"strconv"
	"time"
)

type FormSubmissionResponse struct {
	ID              int       `json:"id"`
	FormSchemaID    int       `json:"form_schema_id"`
	SubmittedValues string    `json:"submitted_values"`
	CreatedAt       time.Time `json:"created_at"`
}

type formSubmissionHandler struct {
	formSubmissionRepository repositories.FormSubmissionRepository
	formSchemaRepository     repositories.FormSchemaRepository
}

func NewFormSubmissionViewHandler(rootRoute *echo.Echo, submissionRepo repositories.FormSubmissionRepository, schemaRepo repositories.FormSchemaRepository) {
	handler := &formSubmissionHandler{
		formSubmissionRepository: submissionRepo,
		formSchemaRepository:     schemaRepo,
	}

	rootRoute.GET("/:id/submissions", handler.GetSubmissionsByFormID)
}

func (h *formSubmissionHandler) GetSubmissionsByFormID(c echo.Context) error {
	formID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	// Get the form schema to display the title
	formSchema, err := h.formSchemaRepository.GetSchemaById(formID)
	if err != nil {
		return err
	}

	var parsedSchema models.JSONSchema
	if err := json.Unmarshal(formSchema.Schema, &parsedSchema); err != nil {
		return err
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	submissions, err := h.formSubmissionRepository.GetSubmissionsByFormID(formID, page, pageSize)
	if err != nil {
		return err
	}

	// Convert to response type with formatted JSON
	var submissionResponses []FormSubmissionResponse
	for _, submission := range submissions {
		var parsedValues interface{}
		if err := json.Unmarshal([]byte(submission.SubmittedValues), &parsedValues); err != nil {
			return err
		}

		prettyJSON, err := json.MarshalIndent(parsedValues, "", "    ")
		if err != nil {
			return err
		}

		submissionResponses = append(submissionResponses, FormSubmissionResponse{
			ID:              submission.ID,
			FormSchemaID:    submission.FormSchemaID,
			SubmittedValues: string(prettyJSON),
			CreatedAt:       submission.CreatedAt,
		})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "partials/submissions-partial.html", map[string]interface{}{
			"Submissions": submissionResponses,
			"FormID":      formID,
			"NextPage":    page + 1,
		})
	}

	return c.Render(http.StatusOK, "submissions.html", map[string]interface{}{
		"Submissions": submissionResponses,
		"FormID":      formID,
		"FormTitle":   parsedSchema.Title,
		"NextPage":    2,
	})
}
