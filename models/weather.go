package models

type Precipitation struct {
	Probability float64 `json:"probability"`
	Rain        bool    `json:"rain"`
	Snow        bool    `json:"snow"`
	Sleet       bool    `json:"sleet"`
	Hail        bool    `json:"hail"`
}

type Wind struct {
	Direction string  `json:"direction"`
	Angle     float64 `json:"angle"`
	Speed     float64 `json:"speed"`
}

type Weather struct {
	Temperature   float64       `json:"temperature"`
	Humidity      float64       `json:"humidity"`
	Pressure      float64       `json:"pressure"`
	Precipitation Precipitation `json:"precipitation"`
	Time          string        `json:"time"`
	Wind          Wind          `json:"wind"`
}
