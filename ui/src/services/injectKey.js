import { inject } from 'vue'

export const socketKey = 'socket'

export function useSocket (key = null) {
    return inject(key !== null ? key : socketKey)
}
