package app_services

import (
	models "app.models"
	"encoding/json"
	"fmt"
)

func (rfs *RpcFleetService) RPCGetFleetInfo(req json.RawMessage) ([]byte, error) {
	rfs.log.Debug("RpcFleetService.RPCGetFleetInfo start")
	fr := &models.FleetInfoRequest{}
	if err := fr.Load(req); err != nil {
		rfs.log.Error("RpcFleetService.RPCGetFleetInfo error load fleet request: ", err)
		return nil, fmt.Errorf("can't load request from query")
	}
	// return nil, fmt.Errorf("error get fleet info")
	resp, err := fr.Get()
	if err != nil {
		rfs.log.Error("RpcFleetService.RPCGetFleetInfo error get fleet request: ", err)
		return nil, fmt.Errorf("can't get data from db")
	}
	resp.Origin = "PHX"
	resp.Destination = "CLT"
	resp.AirSpeed = 900
	resp.AirAltitude = 9000

	resp.AirRoute = models.AirRoute{
		&models.AirRouteWayPoint{AirportId: "PHX", Longitude: -112.01200103759766, Latitude: 33.43429946899414},
		&models.AirRouteWayPoint{AirportId: "DFW", Longitude: -97.038002, Latitude: 32.896801},
		&models.AirRouteWayPoint{AirportId: "RDU", Longitude: -78.7874984741211, Latitude: 35.877601623535156},
		&models.AirRouteWayPoint{AirportId: "CLT", Longitude: -80.94309997558, Latitude: 35.2140007019043},
	}

	return resp.Marshal()
}
