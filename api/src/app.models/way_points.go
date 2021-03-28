package app_models

type AirRoute []*AirRouteWayPoint

type AirRouteWayPoint struct {
	AirportId string  `json:"airport_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
