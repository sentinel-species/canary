package types

type StatusResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type DependencyCheckResponse struct {
	FileContent string `json:"fileContent"`
}
