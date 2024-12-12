package requests

type CreateFormSchemaRequest struct {
	Title      string                 `json:"title"`
	Properties map[string]interface{} `json:"properties"`
}
