import { createStore } from 'vuex'

import * as mutations from './mutations'
import * as getters from './getters'
import * as actions from './actions'
import state from './state'
import Sockets from './Sockets'
import AirRoute from './AirRoute'

const store = createStore({
    state: state,
    getters: getters,
    mutations: mutations,
    actions: actions,
    modules: {
        AirRoute,
        Sockets
    },
    strict: process.env.DEV
})

export default store
