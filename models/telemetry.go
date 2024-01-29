package models

type TelemetryData struct {
	RocketID     string
	Altitude     float32
	Speed        float32
	Acceleration float32
	Thrust       float32
	Temperature  float32
}
