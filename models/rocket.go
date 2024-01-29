package models

const (
	RocketStatusWaiting   string = "waiting"
	RocketStatusLaunched  string = "launched"
	RocketStatusDeployed  string = "deployed"
	RocketStatusFailed    string = "failed"
	RocketStatusCancelled string = "cancelled"
)

type Payload struct {
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}

type Telemetry struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RocketInfo struct {
	ID           string            `json:"id"`
	Model        string            `json:"model"`
	Mass         float64           `json:"mass"`
	Payload      Payload           `json:"payload"`
	Telemetry    Telemetry         `json:"telemetry"`
	Status       string            `json:"status"`
	Timestamps   map[string]string `json:"timestamps"`
	Altitude     float64           `json:"altitude"`
	Speed        float64           `json:"speed"`
	Acceleration float64           `json:"acceleration"`
	Thrust       float64           `json:"thrust"`
	Temperature  float64           `json:"temperature"`
}
