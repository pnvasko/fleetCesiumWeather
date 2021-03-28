export const isRpcResponse = (val) => {
    return (val && val.jsonrpc && val.jsonrpc === '2.0')
}

export const isDirectCall = (val) => {
    return (val && val.params)
}

export const isRpcError = (val) => {
    return (val && val.error && val.error.code !== 0)
}

export default function passToStoreHandler (eventName, event) {
    if (!eventName.startsWith('SOCKET_')) { return }
    let method = 'commit'
    let defaultNamespace = 'Sockets'
    let target = eventName.toUpperCase()
    if (this.format === 'json' && event.data) {
        event.data.split('\n').map((val) => {
            if (val !== ' ' && val !== '' && val !== null && val !== void 0) {
                let post = JSON.parse(val)
                if (isRpcResponse(post)) {
                    if (isDirectCall(post)) {
                        // this.store.dispatch(`${defaultNamespace}/${post.method}`, post.params)
                        this.store.dispatch(`${post.method}`, post.params)
                    } else {
                        this.store['commit']('Sockets/SOCKET_RPC_CALLBACK_ONMESSAGE', JSON.stringify(post))
                    }
                }
            }
        })
    } else {
        target = `${defaultNamespace}/${target}`
        this.store[method](target, event)
    }
}
