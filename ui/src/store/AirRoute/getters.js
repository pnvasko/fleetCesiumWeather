export const getFleetInfo = (state) => (name) => {
    if (state.airRoutes && name in state.airRoutes) {
        return state.airRoutes[name]
    }
    return {}
}

export const getFleetRouteData = (state) => (name) => {
    if (state.fleetRouteData && name in state.fleetRouteData) {
        return state.fleetRouteData[name]
    }
    return []
}

export const getSubscribes = (state) => {
    return state.subscribes
}
