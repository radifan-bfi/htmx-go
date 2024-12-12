package handlers

import (
	"github.com/labstack/echo/v4"
	"htmx-go/repositories"
	"htmx-go/responses"
	"net/http"
	"strconv"
)

type formSchemaHandler struct {
	formSchemaRepository repositories.FormSchemaRepository
}

func NewFormSchemaHandler(rootRoute *echo.Group, formSchemaRepository repositories.FormSchemaRepository) {
	handler := &formSchemaHandler{
		formSchemaRepository: formSchemaRepository,
	}

	rootRoute.GET("/form-schemas", handler.GetSchema)
	rootRoute.GET("/form-schemas/:id", handler.GetSchemaById)
	rootRoute.POST("/form-schemas", handler.CreateSchema)
}

func (h *formSchemaHandler) GetSchema(c echo.Context) error {
	formSchemas, err := h.formSchemaRepository.GetSchema()
	if err != nil {
		return responses.ErrInvalidRequest(err)
	}

	return c.JSON(http.StatusOK, formSchemas)
}

func (h *formSchemaHandler) GetSchemaById(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		return responses.ErrInvalidRequest(err)
	}

	formSchema, err := h.formSchemaRepository.GetSchemaById(id)
	if err != nil {
		return responses.ErrInvalidRequest(err)
	}

	return c.JSON(http.StatusOK, formSchema)
}

func (h *formSchemaHandler) CreateSchema(c echo.Context) error {
	var request string
	if err := c.Bind(request); err != nil {
		return responses.ErrInvalidRequest(err)
	}

	id, err := h.formSchemaRepository.SaveSchema(request)
	if err != nil {
		return responses.ErrInvalidRequest(err)
	}

	return c.JSON(http.StatusCreated, map[string]int64{
		"id": id,
	})
}
