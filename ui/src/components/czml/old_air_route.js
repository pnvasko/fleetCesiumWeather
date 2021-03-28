/*
import {
    // ArcType,
    EllipsoidGeodesic,
    // Color,
    Cartographic,
    CzmlDataSource,
    Cartesian3,
    HeadingPitchRoll,
    // Math as cesiumMath,
    Transforms,
    // VelocityOrientationProperty,
    Ray } from 'cesium'
*/
// import { bearing } from './utils'
// import { czmlAiroplanTemplate, czmlAirportTemplate, czmlHeaderTemplate, czmlRouteTemplate } from './templates'
// import { NoWayPointsException } from './errors'

// import {isObject, isObjectEmpty} from "../../services/utils";

export class XAirRoute {
    _tmp () {
        let a = ref(0)
        let cart = reactive([])
        let total = computed(() => {
            return cart
        })
        let preFetch = effect(async () => {
            cart.push('1')
        })
        console.log(a)
        console.log(cart)
        console.log(total)
        console.log(preFetch)
    }
    constructor (viewer, options = {}) {
        this._viewer = viewer
        this._options = options
        this._height = 1000
        this._wayPoints = options.wayPoints
        this._name = options.name
        this._id = this._getRouteId()
        this._polyline = this._getPolyline(this._wayPoints)
        this._dataSource = null
        this.distance = this._getDistance()


        /*
        const fleetInfo = this._store.getters['AirRoute/getFleetInfo'](this._name)
        this._showAirRoute(fleetInfo)
        this._airAltitude = fleetInfo.airAltitude
        this._airSpeed = fleetInfo.airSpeed
        */

        // this.WebSocket = opts.WebSocket || (protocol === '' ? new WebSocket(connectionUrl) : new WebSocket(connectionUrl, protocol))
        // this.format = opts.format && opts.format.toLowerCase()
        // if (!options || !isObject(options) || isObjectEmpty(options)) { throw new NoOptionException() }
    }

    id () {
        return this._id
    }

    async show () {
        // https://sandcastle.cesium.com/#c=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
        // meridian(this._viewer, this._wayPoints.source.lng, Color.CYAN, 0.001)
        // parallel(this._viewer, this._wayPoints.source.lat, Color.CYAN, 0.001)

        this._dataSource = await CzmlDataSource.load(this._polyline)
        this._viewer.dataSources.add(this._dataSource)

        const ds = this._viewer.dataSources.get(0)
        const model = ds.entities.getById('airoplan-0')
        const position = Cartesian3.fromDegrees(this._wayPoints.source.lng, this._wayPoints.source.lat, this._height)
        model.position = position
        const heading = bearing(this._wayPoints.source, this._wayPoints.destination)
        const pitch = 0
        const roll = 0
        const hpr = new HeadingPitchRoll(heading, pitch, roll)
        model.orientation = Transforms.headingPitchRollQuaternion(position, hpr)
        this._viewer.zoomTo(this._dataSource)

        // const line = ds.entities.getById('line-PHX,CLT')
        console.log('distance: ', this.distance)
    }

    _getDistance () {
        let startCartographic = Cartographic.fromCartesian(Cartesian3.fromDegrees(this._wayPoints.source.lng, this._wayPoints.source.lat, this._height))
        let endCartographic = Cartographic.fromCartesian(Cartesian3.fromDegrees(this._wayPoints.destination.lng, this._wayPoints.destination.lat, this._height))
        let ellipsoidGeodesic = new EllipsoidGeodesic(startCartographic, endCartographic)
        return ellipsoidGeodesic.surfaceDistance * 0.001
    }

    async remove () {
        const ds = await this._dataSource
        this._viewer.dataSources.remove(ds, true)
        this._dataSource = null
    }

    update (route) {
        this._polyline = this._getPolyline(route)
        if (this._dataSource !== null) {
            this.remove()
        }
        this.show()
    }

    _getRouteId () {
        const sN =  this._wayPoints.source.name
        const dN =  this._wayPoints.destination.name
        return `${sN},${dN}`
    }

    _getRouteName () {
        const sN =  this._wayPoints.source.name
        const dN =  this._wayPoints.destination.name
        return `${sN},${dN}`
    }

    _getPolyline (route) {
        const color = [255, 0, 0, 255]
        return [
            lineHeaderTemplate(this._id, this._getRouteName()),
            airportTemplate(route.source.name, route.source, this._height, [255, 0, 0, 255]),
            lineTemplate(this._id, route.source, route.destination, this._height, color),
            airportTemplate(route.destination.name, route.destination, this._height, [255, 0, 0, 255]),
            airoplanCZML('0')
        ]
    }

    _pickGlobeIntersection () {
        let direction = new Cartesian3()
        let p0 = Cartesian3.fromDegrees(this._wayPoints.source.lng, this._wayPoints.source.lat, this._height)
        let p1 = Cartesian3.fromDegrees(this._wayPoints.destination.lng, this._wayPoints.destination.lat, this._height)
        Cartesian3.subtract(p1, p0, direction)
        Cartesian3.normalize(direction, direction)
        let ray = new Ray(p0, direction)
        let hitPos = this._viewer.scene.globe.pick(ray, this._viewer.scene)
        return hitPos
    }
}
