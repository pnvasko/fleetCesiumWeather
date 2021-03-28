import VueNativeSock from './vue-native-websocket'
import { isObject, isObjectEmpty } from '../utils'
import { NoOptionException } from '../errors'
import { socketKey } from '../injectKey'
import { mutations } from './mutations'
import passToStoreHandler from './handlers'
import JsonRpcWebSocketClient from './wssocket_rpc'

export class WsSocket {
    constructor (options = {}) {
        if (!options || !isObject(options) || isObjectEmpty(options)) {
            throw new NoOptionException(
                `
                    you are passing ${options} as your root module. please provide a valid object format
                    your object should contain [router, store, strict, plugins]
                `
            )
        }
        const {
            router,
            store,
            plugins = [],
            strict = false
        } = options

        this.strict = strict
        this._store = store
        this._router = router
        this._wsuri = process.env.VUE_APP_WS_API_URL
        this.options = {
            store: this._store,
            mutations: mutations,
            connectManually: false,
            reconnection: true,
            reconnectionAttempts: 0,
            reconnectionDelay: 3000,
            format: 'json',
            passToStoreHandler: passToStoreHandler
        }

        plugins.forEach(plugin => plugin(this))
    }

    install (app, injectKey) {
        app.provide(injectKey || socketKey, this)
        app.config.globalProperties.$wssocket = this
        this._store.commit(mutations.SOCKET_ONCLOSE, 'init')
        app.use(VueNativeSock, this._wsuri, this.options)
        app.config.globalProperties.$rpc = new JsonRpcWebSocketClient({ app: app, socket: app.config.globalProperties.$socket, store: this._store })
        app.config.globalProperties.$rpc.init()
    }
}
