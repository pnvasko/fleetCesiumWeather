export function getConnectedStatus (state) {
    return state.isConnected
}

export function getReconnectionStatus (state) {
    return state.reconnection
}

export function getStateUi (state) {
    return state.stateUi
}

export function getReconnectErrorStatus (state) {
    return state.reconnectError
}

export function getSocket (state) {
    return state.socket
}

export function getRpcCallbackMsg (state) {
    return state.rpcmessage
}
