import { EventEmitter } from 'events'
import fromEvents from 'promise-toolbox/fromEvents'
import { Peer } from './peer'
import { MethodNotFound } from 'json-rpc-peer'
import { BaseError } from 'make-error'

export class ConnectionError extends BaseError {}

export const CLOSED = 'closed'
export const CONNECTING = 'connecting'
export const MESSAGE = 'message'
export const OPEN = 'open'

function safePromiseWithTimeout (promise, timeout) {
    return new Promise((resolve) => {
        let timerId = setTimeout(() => { resolve({ data: null, err: { 'message': 'timeout error.' } }) }, timeout)
        promise.then(data => {
            clearTimeout(timerId)
            resolve({ data: data, err: null })
        }).catch(error => {
            clearTimeout(timerId)
            resolve({ data: null, err: error })
        })
    })
}

function safePromise (promise) {
    return promise.then(data => [ data, null ]).catch(error => [ null, error ])
}

export default class JsonRpcWebSocketClient extends EventEmitter {
    constructor (opts) {
        super()

        this._opts = opts
        this._app = opts.app
        this._socket = opts.socket
        this._store = opts.store

        const peer = (this._peer = new Peer(message => {
            if (message.type !== 'notification') {
                throw new MethodNotFound()
            }
            this.emit('notification', message)
        }).on('data', message => {
            this.send(message)
        }).on('error', message => {
            this.send(message)
        }))

        this.on(CLOSED, () => {
            peer.failPendingRequests(
                new ConnectionError('connection has been closed')
            )
        })

        this.on(MESSAGE, message => {
            peer.write(message)
        })
    }

    safeCallWithTimeout (method, payload, timeout, callback) {
        let params = { ...payload }
        params.callback = callback
        this._store.commit('Sockets/SOCKET_RPC_CALLBACK_ONMESSAGE', null)
        return safePromiseWithTimeout(this._peer.request(method, params), timeout)
    }

    safecall (method, payload, callback) {
        let params = { ...payload }
        params.callback = callback
        this._store.commit('Sockets/SOCKET_RPC_CALLBACK_ONMESSAGE', null)
        return safePromise(this._peer.request(method, params))
    }

    call (method, params) {
        return this._peer.request(method, params)
    }

    notify (method, params) {
        return this._peer.notify(method, params)
    }

    init () {
        this._store.watch(() => this._store.getters['Sockets/getConnectedStatus'], async () => {
            this._socket = this._store.getters['Sockets/getSocket']
        })

        this._store.watch(() => this._store.getters['Sockets/getRpcCallbackMsg'], async (data) => {
            if (data !== null && data !== '') {
                this.emit(MESSAGE, data)
            }
        })
    }

    directSocketInit () {
        return fromEvents(this._socket, ['open'], ['close', 'error']).then(
            () => {
                this._socket.addEventListener('close', this._onClose)
                this._socket.addEventListener('error', error => {
                    this.emit('error', error)
                })
                this._socket.addEventListener('message', ({ data }) => {
                    this.emit(MESSAGE, data)
                })
            }
        )
    }

    send (data) {
        // console.log('JsonRpcWebSocketClient send start: ', this._app.$vnsocket)
        // this._app.app.config.globalProperties.$disconnect()
        // this._app.app.config.globalProperties.$connect()
        if (this._socket.readyState === this._socket.OPEN) {
            this._socket.send(data)
        } else {
            console.log('Todo JsonRpcWebSocketClient terminate connection start', this._socket.readyState)
        }
    }

    _onClose () {
        this.emit(CLOSED)
    }
}
