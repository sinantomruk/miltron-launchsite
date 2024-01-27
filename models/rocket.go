package models

type Payload struct {
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}

type Telemetry struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Timestamps struct {
	Launched  string `json:"launched"`
	Deployed  string `json:"deployed"`
	Failed    string `json:"failed"`
	Cancelled string `json:"cancelled"`
}

type RocketInfo struct {
	ID           string     `json:"id"`
	Model        string     `json:"model"`
	Mass         int        `json:"mass"`
	Payload      Payload    `json:"payload"`
	Telemetry    Telemetry  `json:"telemetry"`
	Status       string     `json:"status"`
	Timestamps   Timestamps `json:"timestamps"`
	Altitude     float64    `json:"altitude"`
	Speed        float64    `json:"speed"`
	Acceleration float64    `json:"acceleration"`
	Thrust       int        `json:"thrust"`
	Temperature  float64    `json:"temperature"`
}
