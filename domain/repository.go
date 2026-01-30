package domain

type ResponseFormat struct {
	Type   string
	Schema *string
}

type Tool struct {
	Name           string
	Description    string
	ParametersJSON string
}

type ToolCall struct {
	Id        string
	Name      string
	Arguments string
}

type GenerationParams struct {
	Temperature    *float32
	MaxTokens      *int32
	TopK           *int32
	TopP           *float32
	ResponseFormat *ResponseFormat
	Tools          []Tool
}
