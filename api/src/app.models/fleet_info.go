package app_models

import (
	base "app.base"
	"bytes"
)

type FleetInfoRequest struct {
	Fleet string `json:"fleet"`
}

func (freq *FleetInfoRequest) Get() (*FleetInfoResponse, error) {
	fresp := FleetInfoResponse{
		Fleet: freq.Fleet,
	}
	return &fresp, nil
}

func (freq *FleetInfoRequest) Load(data []byte) error {
	return base.JsonProcessor.NewDecoder(bytes.NewReader(data)).Decode(freq)
}

type FleetInfoResponse struct {
	Fleet       string   `json:"fleet"`
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
	AirRoute    AirRoute `json:"air_route"`
	AirSpeed    float64  `json:"air_speed"`
	AirAltitude float64  `json:"air_altitude"`
}

func (fresp *FleetInfoResponse) Marshal() ([]byte, error) {
	data, err := base.JsonProcessor.Marshal(fresp)
	return data, err
}
