import {
    // ArcType,
    // EllipsoidGeodesic,
    // Color,
    // Cartographic,
    CzmlDataSource,
    Cartesian3,
    HeadingPitchRoll,
    // Math as cesiumMath,
    Transforms, EllipsoidGeodesic, Cartographic
    // VelocityOrientationProperty,
    // Ray
} from 'cesium'
import routeStyles from './style'
import { bearing } from './utils'
import { czmlHeaderTemplate, czmlAiroplanTemplate, czmlAirportTemplate, czmlRouteTemplate } from './templates'

export class AirRoute {
    constructor (viewer, store, name, wayPoints, options = {}) {
        this._viewer = viewer
        this._store = store
        this._name = name

        this._options = options
        this._style = options.routeStyles || routeStyles(0)
        this._height = options.height || 5
        this._groundHeight = options.groundHeight || 5
        this._widthRouteLine = options.widthRouteLine || 2
        this._animationInterval = options.animationInterval || 50

        this._animationTimer = null
        this._dataSource = null
        this._geodesic = null

        this._czmlModel = null
        this._airRoute = null
        this._routeData = null
        this._airAltitude = 0
        this._airSpeed = 0
        this._orientation = null
        this._nextDestination = 0

        this._airRouteIds = {}
        this._airportIds = {}

        this._unwatchFleetInfo = this._store.watch(
            () => this._store.getters['AirRoute/getFleetInfo'](this._name),
            async (data) => {
                this._updateFleetInfo(data)
                const ds = this._viewer.dataSources.get(0)
                // ds.load(this._getCzmlModel())
                if (ds !== null) {
                    const arp = this._getAirRoutePolylines()
                    const ars = this._getAirports()
                    const arpIds = arp.map(row => row.id)
                    const arsIds = ars.map(row => row.id)
                    const ids = [...arpIds, ...arsIds]

                    const airIds = [...Object.keys(this._airRouteIds), ...Object.keys(this._airportIds)]
                    let totalIds = airIds.length

                    for (let i = 0; i < totalIds; i++) {
                        let eId = airIds[i]
                        if (!ids.includes(eId)) {
                            ds.entities.removeById(eId)
                            if (eId.startsWith('route')) {
                                delete this._airRouteIds[eId]
                            } else if (eId.startsWith('airport')) {
                                delete this._airportIds[eId]
                            }
                        }
                    }
                    ds.process([...arp, ...ars])
                }
            }, { deep:true }
        )

        this._unwatchFleetRouteData = this._store.watch(
            () => this._store.getters['AirRoute/getFleetRouteData'](this._name),
            async (data) => {
                console.log('AirRoute/getFleetRouteData:', data)
                this._routeData = data
            }, { deep:true }
        )
    }

    init () {
        const fleetInfo = this._store.getters['AirRoute/getFleetInfo'](this._name)
        this._updateFleetInfo(fleetInfo)
        const source = this._airRoute[0]
        const destination = this._airRoute[1]
        if (this._airRoute.length > 2) {
            this._nextDestination = 2
        }

        const startPosition = {
            lat: source.latitude,
            lng: source.longitude,
            altitude: 5,
            speed: 0,
            timestamp: 0,
        }
        this._routeData = [startPosition]
        this._orientation = this._getOrientation(
            { lat: source.latitude, lng: source.longitude },
            { lat: destination.latitude, lng: destination.longitude },
        )

        this._czmlModel = this._getCzmlModel()
        this._geodesic = new EllipsoidGeodesic(
            Cartographic.fromDegrees(source.longitude, source.latitude, this._height),
            Cartographic.fromDegrees(destination.longitude, destination.latitude, this._height)
        )
    }

    _getOrientation (source, destination) {
        const position = Cartesian3.fromDegrees(source.lng, source.lat, this._height)
        const heading = bearing(source, destination)
        const pitch = 0
        const roll = 0
        const hpr = new HeadingPitchRoll(heading, pitch, roll)
        return Transforms.headingPitchRollQuaternion(position, hpr)
    }

    async show () {
        if (this._czmlModel === null || this._czmlModel < 2) {
            return
        }
        this._dataSource = await CzmlDataSource.load(this._czmlModel)
        this._viewer.dataSources.add(this._dataSource)

        const ds = this._viewer.dataSources.get(0)
        const model = ds.entities.getById(`airoplan-${this._name}`)
        model.orientation = this._orientation

        // this._airAltitude = data.airAltitude
        if (this._airSpeed > 0) {
            this.animationStart()
        }
    }

    async animationReStart () {
        console.log('animationReStart', this._nextDestination)
        await this.animationStop()
        if (this._nextDestination > 0 && this._nextDestination < this._airRoute.length) {
            const source = this._airRoute[this._nextDestination - 1]
            const destination = this._airRoute[this._nextDestination]
            if (this._airRoute.length > this._nextDestination) {
                this._nextDestination = this._nextDestination + 1
            }
            this._orientation = this._getOrientation(
                { lat: source.latitude, lng: source.longitude },
                { lat: destination.latitude, lng: destination.longitude },
            )
            this._geodesic = new EllipsoidGeodesic(
                Cartographic.fromDegrees(source.longitude, source.latitude, this._height),
                Cartographic.fromDegrees(destination.longitude, destination.latitude, this._height)
            )
            this.animationStart()
        }
    }

    async animationStart () {
        console.log('animationStart: ')
        const ds = this._viewer.dataSources.get(0)
        const model = ds.entities.getById(`airoplan-${this._name}`)
        // Todo debug Speed
        // this._animationInterval = 50
        model.orientation = this._orientation
        this._airSpeed = 400000000

        let distance = this._geodesic.surfaceDistance
        let totalTime = 3600 * distance / this._airSpeed
        let unitFraction = this._animationInterval / ( totalTime * 1000)
        console.log('animationInterval: ', this._animationInterval)
        console.log('distance: ', distance)
        console.log('time: ', totalTime)
        console.log('unitFraction: ', unitFraction)
        let fraction = 0
        if (this._animationTimer === null) {
            this._animationTimer = setInterval(() => {
                if (fraction >= 1.0) {
                    console.log('animationStart clearInterval: ')
                    clearInterval(this._animationTimer)
                    this._animationTimer = null
                    this.animationReStart()
                }
                let fractionalPoint = this._geodesic.interpolateUsingFraction(fraction)
                fractionalPoint.height = this._airAltitude
                model.position = Cartographic.toCartesian(fractionalPoint)
                fraction += unitFraction
            }, this._animationInterval)

            /*
            this._animationTimer = setInterval(() => {
                startLatitude += 0.05
                startLongitude += 0.05
                model.position = Cartesian3.fromDegrees(startLongitude, startLatitude, this._height)
            }, this._animationInterval)
            */
        }
    }

    async animationStop () {
        if (this._animationTimer !== null) {
            clearInterval(this._animationTimer)
            this._animationTimer = null
        }
    }

    async update () { }

    async remove () {
        if (this._animationTimer !== null) {
            clearInterval(this._animationTimer)
            this._animationTimer = null
        }
        if (this._dataSource === null) {
            return
        }
        this._viewer.dataSources.remove(this._dataSource, true)
        this._dataSource = null
        this._airRouteIds = {}
        this._airportIds = {}
    }

    async destroy () {
        console.log('AirRoute.destroy:')
        this._unwatchFleetInfo()
        this._unwatchFleetRouteData()
        await this.remove()
    }

    _getCzmlModel () {
        return [
            czmlHeaderTemplate(this._name),
            ...this._getAirRoutePolylines(),
            ...this._getAirports(),
            ... this._getAirplane()
        ]
    }

    _updateFleetInfo (data) {
        this._airRoute = data.airRoute
        this._airAltitude = data.airAltitude
        this._airSpeed = data.airSpeed * 1000
    }

    _getAirplane () {
        let airplane =[]
        const position = this._routeData[this._routeData.length - 1]
        airplane.push(czmlAiroplanTemplate(this._name, position))
        return airplane
    }

    _getAirRoutePolylines () {
        if (this._airRoute === null || this._airRoute.length < 2) {
            return []
        }

        const rLen = this._airRoute.length
        let polylines = []
        for (let i = 0; i < rLen - 1; i++ ) {
            this._airRouteIds[`route-${this._airRoute[i].airport_id}-${this._airRoute[i + 1].airport_id}`] = true
            polylines.push(this._getAirRoutePolyline(this._airRoute[i], this._airRoute[i + 1], this._height, this._style.Route, this._style.RouteWidth))
        }
        return polylines
    }

    _getAirports () {
        if (this._airRoute === null || this._airRoute.length < 2) {
            return []
        }
        const rLen = this._airRoute.length
        let airports = []
        for (let i = 0; i < rLen; i++ ) {
            this._airportIds[`airport-${this._airRoute[i].airport_id}`] = true
            airports.push(this._getAirport(this._airRoute[i], this._groundHeight, this._style.AiroportRadius, this._style.AiroportBg, this._style.AiroportBorder))
        }
        return airports
    }

    _getAirRoutePolyline (start, end, height, color, width) {
        return czmlRouteTemplate(
            `${start.airport_id}-${end.airport_id}`,
            { lat: start.latitude, lng: start.longitude },
            { lat: end.latitude, lng: end.longitude },
            height,
            color,
            width
        )
    }

    _getAirport (position, height, width, color, colorBr) {
        return czmlAirportTemplate(
            position.airport_id,
            { lat: position.latitude, lng: position.longitude },
            height,
            width,
            color,
            colorBr
        )
    }
}

export function createAirRoute (viewer, route, options) {
    return new AirRoute(viewer, route, options)
}
