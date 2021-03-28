import { EventEmitter } from 'events'
import { format, JsonRpcError, MethodNotFound } from 'json-rpc-protocol'
import { forEach, isArray, map } from 'lodash'
import { UUID} from './uuid'

const JSONRPC = '2.0'

export class InvalidJson extends JsonRpcError {
    constructor () {
        super('invalid JSON', -32700)
    }
}

export class InvalidRequest extends JsonRpcError {
    constructor () {
        super('invalid JSON-RPC request', -32600)
    }
}
const getType = (value) => (value === null ? 'null' : typeof value)

const isString = (value) => typeof value === 'string'

export const isNumber = (value) => {
    const type = typeof value
    return type === 'number' && value > negativeInf && value < positiveInf
}

const isNotificationPayload = (message) => {
    if (isString(message.method)) {
        checkParams(message.params)
        return true
    }
    return false
}

const setMessageType = (message, type) => defineProperty(message, 'type', {
    configurable: true,
    value: type,
    writable: true
})

const isRequestPayload = (message) => {
    if (isString(message.method)) {
        const { id } = message
        checkId(id)
        checkParams(message.params)
        return true
    }
    return false
}

const checkId = (id) => {
    if (!isNumber(id) && !isString(id)) {
        throw new InvalidRequest(`invalid identifier: ${getType(id)} instead of number or string`)
    }
}

const checkParams = (params) => {
    if (params !== undefined && !Array.isArray(params) && !isObject(params)) {
        throw new InvalidRequest(`invalid params: ${getType(params)} instead of undefined, array or object`)
    }
}

const checkError = (error) => {
    if (error == null || !isInteger(error.code) || !isString(error.message)) {
        throw new InvalidRequest(`invalid error: ${getType(error)} instead of {code, message}`)
    }
}

const isResponsePayload = (message) => {
    if (isObject(message.result)) {
        checkId(message.id)
        checkParams(message.data)
        return true
    } else if (!isString(message.method)) {
        if (!isErrorResponse(message)) {
            const { id } = message
            checkId(id)
            return true
        }
    }
    return false
}

const isErrorPayload = (message) => {
    // Todo fix isErrorPayload
    if (!isString(message.method)) {
        const { id } = message
        if (id !== null) {
            checkId(id)
        }
        checkError(message.error)
        return true
    }
    return false
}

const { defineProperty } = Object
const negativeInf = Number.NEGATIVE_INFINITY
const positiveInf = Number.POSITIVE_INFINITY
const isErrorResponse = ({ error }) => error !== undefined
export const isInteger = (value) => isNumber(value) && value % 1 === 0

export const isObject = (value) => {
    const type = typeof value
    return value !== null && (type === 'object' || type === 'function')
}

function makeAsync (fn) {
    return function () {
        return new Promise(resolve => resolve(fn.apply(this, arguments)))
    }
}

const parseMessage = message => {
    try {
        return parse(message)
    } catch (error) {
        throw format.error(null, error)
    }
}

function defaultOnMessage (message) {
    if (message.type === 'request') {
        throw new MethodNotFound(message.method)
    }
}

function parse (message) {
    if (isString(message)) {
        try {
            message = JSON.parse(message)
        } catch (error) {
            if (error instanceof SyntaxError) {
                throw new InvalidJson()
            }

            throw error
        }
    }

    if (Array.isArray(message)) {
        return message.map(parse)
    }

    if (isNotificationPayload(message)) {
        setMessageType(message, 'notification')
    } else if (isRequestPayload(message)) {
        setMessageType(message, 'request')
    } else if (isResponsePayload(message)) {
        setMessageType(message, 'response')
    } else if (isErrorPayload(message)) {
        setMessageType(message, 'error')
    } else {
        throw new InvalidJson()
    }

    return message
}

function noop () {}

export class Peer extends EventEmitter {
    constructor (onMessage = defaultOnMessage) {
        super()

        this._asyncEmitError = () => {
            process.nextTick.bind(process, this.emit.bind(this), 'error')
        }
        this._handle = makeAsync(onMessage)
        this._deferreds = Object.create(null)
    }

    _getDeferred (id) {
        const deferred = this._deferreds[id]
        delete this._deferreds[id]
        return deferred
    }

    async exec (message, data) {
        message = parseMessage(message)
        if (isArray(message)) {
            const results = []
            // Only returns non empty results.
            await Promise.all(
                map(message, message => {
                    return this.exec(message, data).then(result => {
                        if (result !== undefined) {
                            results.push(result)
                        }
                    })
                })
            )

            return results
        }
        const { type } = message
        if (type === 'error') {
            const { id } = message
            // Some errors do not have an identifier, simply discard them.
            if (id === null) {
                return
            }

            const { error } = message
            let err = new JsonRpcError(error.message, error.code, error.data)
            // TODO: it would be great if we could return an error with of
            // a more specific type (and custom types with registration).
            this._getDeferred(id).reject(err)
        } else if (type === 'response') {
            this._getDeferred(message.id).resolve(message.result)
        } else if (type === 'notification') {
            this._handle(message, data).catch(noop)
        } else {
            return this._handle(message, data)
                .then(result =>
                    format.response(message.id, result === undefined ? null : result)
                )
                .catch(error =>
                    format.error(
                        message.id,

                        // If the method name is not defined, default to the method passed
                        // in the request.
                        error instanceof MethodNotFound && !error.data
                            ? new MethodNotFound(message.method)
                            : error
                    )
                )
        }
    }

    failPendingRequests (reason) {
        const { _deferreds: deferreds } = this

        forEach(deferreds, ({ reject }, id) => {
            reject(reason)
            delete deferreds[id]
        })
    }

    request (method, params) {
        return new Promise((resolve, reject) => {
            const requestId = this._createID()

            // this.push(format.request(requestId, method, params))
            let req = {
                id: requestId,
                jsonrpc: JSONRPC,
                method: method,
                payload: params
            }
            this.push(JSON.stringify(req))
            this._deferreds[requestId] = { resolve, reject }
        })
    }

    async notify (method, params) {
        this.push(format.notification(method, params))
    }

    end (data, encoding, cb) {
        if (typeof data === 'function') {
            process.nextTick(data)
        } else {
            if (typeof encoding === 'function') {
                process.nextTick(encoding)
            } else if (typeof cb === 'function') {
                process.nextTick(cb)
            }

            if (data !== undefined) {
                this.write(data)
            }
        }
    }

    pipe (writable) {
        const listeners = {
            data: data => writable.write(data),
            end: () => {
                writable.end()
                clean()
            }
        }

        const clean = () =>
            forEach(listeners, (listener, event) => {
                this.removeListener(event, listener)
            })

        forEach(listeners, (listener, event) => {
            this.on(event, listener)
        })

        return writable
    }

    push (data) {
        return data === null ? this.emit('end') : this.emit('data', data)
    }

    write (message) {
        let cb
        const n = arguments.length
        if (n > 1 && typeof (cb = arguments[n - 1]) === 'function') {
            process.nextTick(cb)
        }

        this.exec(String(message)).then(response => {
            if (response !== undefined) {
                this.push(response)
            }
        }, this._asyncEmitError)

        // indicates that other calls to `write` are allowed
        return true
    }

    _createID () {
        return UUID()
    }
}
