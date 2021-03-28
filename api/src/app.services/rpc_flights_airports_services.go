package app_services

import (
	models "app.models"
	"encoding/json"
	"fmt"
)

func (rfs *RpcFleetService) RPCGetAirportsGeo(req json.RawMessage) ([]byte, error) {
	rfs.log.Debug("RpcFleetService.RPCGetFleetScheduled start")
	ar := &models.AirRouteRequest{}

	if err := ar.Load(req); err != nil {
		rfs.log.Error("RpcFleetService.RPCGetAirportsGeo error load air route: ", err)
		return nil, fmt.Errorf("error load air route from query")
	}

	resp := &models.AirRouteResponse{}

	if err := resp.GetLocations(rfs.ctx, ar); err != nil {
		rfs.log.Error("RpcFleetService.RPCGetAirportsGeo error locations for air route: ", err)
		return nil, fmt.Errorf("error get locations for air route from db")
	}

	return resp.Marshal()
}
