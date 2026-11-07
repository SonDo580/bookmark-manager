package model

type HealthStatus string

const (
	HealthStatusOK HealthStatus = "OK"
)

type HealthCheckResponse struct {
	Message     HealthStatus `json:"message"`
	ServiceName string       `json:"service_name"`
	InstanceID  string       `json:"instance_id"`
}
