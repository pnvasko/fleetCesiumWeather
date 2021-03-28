// import Vue from 'vue'

export function SOCKET_ONOPEN (state, event) {
    console.info('SOCKET_ONOPEN state:', state)
    console.info('SOCKET_ONOPEN:', event)
    state.socket = event.currentTarget
    // state.socket = Vue.prototype.$socket = event.currentTarget
    state.isConnected = true
    state.reconnection = false
    state.stateUi = 'idle'
}

export function SOCKET_ONCLOSE (state, event) {
    console.info('SOCKET_ONCLOSE:', event)
    state.isConnected = false
    state.reconnection = false
    state.stateUi = 'offline'
}

export function SOCKET_ONERROR (state, event) {
    console.info('SOCKET_ONERROR', event)
    state.isConnected = false
    state.reconnection = false
    state.stateUi = 'error'
}

export function SOCKET_ONMESSAGE (state, message) {
    console.info('SOCKET_ONMESSAGE:', message)
}

export function SOCKET_RECONNECT (state, count) {
    console.info('SOCKET_RECONNECT:', count)
    state.reconnection = true
    state.stateUi = 'reconnection'
}

export function SOCKET_RECONNECT_ERROR (state) {
    console.info('SOCKET_RECONNECT_ERROR:')
    state.reconnectError = true
    state.reconnection = false
    state.stateUi = 'error'
}

export function SOCKET_RPC_CALLBACK_ONMESSAGE (state, message) {
    state.rpcmessage = message
}
