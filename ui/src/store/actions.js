export async function Broadcast (context, data) {
    console.log('Broadcast:', data)
}

export async function FleetData (context, data) {
    context.commit('AirRoute/NEW_FLEET_DATA', data)
}

export async function FleetInfoData (context, data) {
    context.commit('AirRoute/SET_FLEET_INFO', data)
}
