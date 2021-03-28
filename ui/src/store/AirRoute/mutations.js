export function SET_FLEET_INFO (state, data) {
    state.airRoutes[data.fleet] = {
        airRoute: data.air_route,
        airAltitude: data.air_altitude,
        airSpeed: data.air_speed
    }
}

export function NEW_FLEET_DATA (state, data) {
    console.info('NEW_FLEET_DATA:', data)
    if (!(data.fleet in state.fleetRouteData)) {
        state.fleetRouteData[data.fleet] = []
    }
    state.fleetRouteData[data.fleet].push({
        lat: data.lat,
        lng: data.lng,
        altitude: data.speed,
        speed: data.alt,
        timestamp: data.timestamp
    })
}

export function SET_SUBSCRIBE (state, data) {
    state.subscribes[data.fleet] = data.timestamp
}

export function REMOVE_SUBSCRIBE (state, data) {
    if (data.chanel in state.subscribes) {
        delete state.subscribes[data.chanel]
    }
    if (data.chanel in state.fleetRouteData) {
        delete state.fleetRouteData[data.chanel]
    }
}


