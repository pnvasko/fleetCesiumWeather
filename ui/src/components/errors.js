import { toastController } from '@ionic/vue'

const errorKey = 'appErrors'

export function createAppErrors (options) {
    return new AppErrors(options)
}

export class AppErrors {
    constructor (options = {}) {
        this._options = options
    }
    install (app, injectKey) {
        app.provide(injectKey || errorKey, this)
        app.config.globalProperties.$appErrors = this
        app.config.globalProperties.$appTxtErrors = this.errorList()
    }

    async showError (msg) {
        const toast = await toastController.create({
            color: 'danger',
            position: 'middle',
            duration: 5000,
            message: msg,
            showCloseButton: true
        })
        await toast.present()
    }

    errorList () {
        return {
            errorGeneral: 'General error',
            errorGetFleetInfo: 'Error get fleet info',
            errorSubscribeToFleet: 'Error subscribe to fleet update',
            errorGetWayPoints: 'Error get air route way points'
        }
    }
}
