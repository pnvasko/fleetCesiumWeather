import Observer from './Observer'
import Emitter from './Emitter'

export default {
    install (app, connection, opts = {}) {
        if (!connection) { throw new Error('[vue-native-socket] cannot locate connection') }
        let observer = null

        opts.$setInstance = (wsInstance) => {
            console.log('wsInstance: ', wsInstance)
            app.config.globalProperties.$socket = wsInstance
        }

        if (opts.connectManually) {
            app.config.globalProperties.$connect = (connectionUrl = connection, connectionOpts = opts) => {
                connectionOpts.$setInstance = opts.$setInstance
                observer = new Observer(connectionUrl, connectionOpts)
                app.config.globalProperties.$socket = observer.WebSocket
            }

            app.config.globalProperties.$disconnect = () => {
                if (observer && observer.reconnection) { observer.reconnection = false }
                if (app.config.globalProperties.$socket) {
                    app.config.globalProperties.$socket.close()
                    delete app.config.globalProperties.$socket
                }
            }
        } else {
            observer = new Observer(connection, opts)
            app.config.globalProperties.$socket = observer.WebSocket
        }

        const hasProxy = typeof Proxy !== 'undefined' && typeof Proxy === 'function' && /native code/.test(Proxy.toString())
        console.log('hasProxy:', hasProxy)
        app.mixin({
            created() {
                let vm = this
                let sockets = this.$options['sockets']

                if (hasProxy) {
                    this.$options.sockets = new Proxy({}, {
                        set (target, key, value) {
                            Emitter.addListener(key, value, vm)
                            target[key] = value
                            return true
                        },
                        deleteProperty (target, key) {
                            Emitter.removeListener(key, vm.$options.sockets[key], vm)
                            delete target.key
                            return true
                        }
                    })
                    if (sockets) {
                        Object.keys(sockets).forEach((key) => {
                            this.$options.sockets[key] = sockets[key]
                        })
                    }
                } else {
                    Object.seal(this.$options.sockets)

                    // if !hasProxy need addListener
                    if (sockets) {
                        Object.keys(sockets).forEach(key => {
                            Emitter.addListener(key, sockets[key], vm)
                        })
                    }
                }
            },
            beforeUnmount () {
                if (hasProxy) {
                    let sockets = this.$options['sockets']

                    if (sockets) {
                        Object.keys(sockets).forEach((key) => {
                            delete this.$options.sockets[key]
                        })
                    }
                }
            }
        })
    },
    restart () {
        console.log('vue-native-socket restart')
    }
}
