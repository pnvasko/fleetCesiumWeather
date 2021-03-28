<template>
  <ion-content :fullscreen="true">
    <ion-grid>
      <ion-row>
        <ion-button @click="subscribeToFleet('SWA2558')">subscribeFleet</ion-button>
        <ion-button @click="unSubscribeToFleet('SWA2558')">unSubscribeToFleet</ion-button>
      </ion-row>
      <ion-row style="height: 730px">
        <ion-col>
          <div id="cesiumContainer" style="width:calc(100vw);height:720px;"></div>
        </ion-col>
      </ion-row>
      <ion-row>
        Debug info
      </ion-row>
      <ion-row>
        <ion-col></ion-col>
        <ion-col></ion-col>
        <ion-col></ion-col>
      </ion-row>
    </ion-grid>

  </ion-content>
</template>

<script>
import 'cesium/Build/Cesium/Widgets/widgets.css'
import { IonContent, IonGrid, IonRow, IonCol, toastController } from '@ionic/vue'
import {
  Clock,
  EllipsoidTerrainProvider,
  OpenStreetMapImageryProvider,
  Viewer,
  Cartesian3 } from "cesium"
import { createAirRoute } from '@/components/czml'

const DebugAirRoute = [
  { source: 'DFW', destination: 'SLC' },
  { source: 'MMLP', destination: 'UAO' },
  { source: 'DFW', destination: 'RDU' },
  { source: 'PHX', destination: 'CLT' },
  { source: 'MCO', destination: 'KAPA' },
  { source: 'PHL', destination: 'PWM' },
  { source: 'ATL', destination: 'LAX' },
  { source: 'KMLB', destination: 'YIP' },
  { source: 'DEN', destination: 'TPA' },
  { source: 'PBI', destination: 'IAD' },
  { source: 'PHL', destination: 'IAH' },
  { source: 'IAH', destination: 'DCA' },
  { source: 'TPA', destination: 'DTW' },
  { source: 'DEN', destination: 'SFO' },
  { source: 'PHX', destination: 'BOS' },
  { source: 'DEN', destination: 'ATL' },
  { source: 'LAS', destination: 'KMDT' },
  { source: 'KEGE', destination: 'EWR' },
  { source: 'ATL', destination: 'LAS' },
  { source: 'TUS', destination: 'MDW' },
  { source: 'CLT', destination: 'SAV' },
  { source: 'PHX', destination: 'DFW' },
  { source: 'BOS', destination: 'MIA' },
  { source: 'DEN', destination: 'EWR' },
  { source: 'PHX', destination: 'KEUG' },
  { source: 'MEM', destination: 'TJBQ' },
  { source: 'CVG', destination: 'ATL' },
  { source: 'DFW', destination: 'FAT' },
  { source: 'DFW', destination: 'MCO' },
  { source: 'FLL', destination: 'DTW' }
]

export default  {
  name: 'Home',
  components: { IonContent, IonGrid, IonRow, IonCol },
  data () {
    return {
      cesiumViewer: null,
      airRouteSelected: '',
      wayPoints: null,
      airRoute: null
    }
  },
  computed: {
    debugAirRoute () {
      return DebugAirRoute.map(row => {
        return `${row.source},${row.destination}`
      })
    }
  },
  mounted() {
    window.CESIUM_BASE_URL = 'http://192.168.3.254:8080/'
    this.initMap()
  },
  methods: {
    initMap () {
      this.cesiumViewer = new Viewer("cesiumContainer", {
        imageryProvider: new OpenStreetMapImageryProvider({
          url: '/maps'
        }),
        animation: false,
        baseLayerPicker: false,
        fullscreenButton: false,
        geocoder: false,
        homeButton: false,
        infoBox: false,
        sceneModePicker: false,
        selectionIndicator: false,
        timeline: false,
        navigationHelpButton: false,
        scene3DOnly: true,
        clock: new Clock(),
        selectedImageryProviderViewModel: undefined,
        selectedTerrainProviderViewModel: undefined,
        terrainProvider: new EllipsoidTerrainProvider(),
        contextOptions: {
          id: "cesiumCanvas",
          webgl: {
            preserveDrawingBuffer: true
          }
        }
      })
      this.cesiumViewer.scene.globe.enableLighting = false
      this.cesiumViewer.scene.globe.showGroundAtmosphere = false
      this.cesiumViewer.scene.debugShowFramesPerSecond = false
      this.cesiumViewer.scene.canvas.id = "cesiumCanvas"
      this.cesiumViewer.shadows = true
      this.cesiumViewer.camera.flyTo({
        destination: Cartesian3.fromDegrees(-122.384, 37.62, 6000000)
      })
    },
    changeAirRouteSelected (ev) {
      this.airRouteSelected =  ev.detail.value
    },
    async subscribeToFleet (fleet) {
      if (this.airRoute !== null) {
        return
      }
      let data, err;

      ({ data, err } = await this.$rpc.safeCallWithTimeout('Fleet.GetFleetInfo', { fleet: fleet }, 5000))

      if (err !== null) {
        const message =  err ?  err.message || 'Test error' : ''
        await this.$appErrors.showError(`${this.$appTxtErrors.errorGetFleetInfo}: ${message}`)
        return
      }
      this.$store.commit('AirRoute/SET_FLEET_INFO', data);

      ({ data, err } = await this.$rpc.safeCallWithTimeout('Subscribe', { name: fleet }, 5000))
      if (err !== null) {
        const message =  err ?  err.message || 'Test error' : ''
        await this.$appErrors.showError(`${this.$appTxtErrors.errorSubscribeToFleet}: ${message}`)
        return
      }
      this.$store.commit('AirRoute/SET_SUBSCRIBE', data)

      this.airRoute = createAirRoute(this.cesiumViewer, this.$store, fleet, { height: 5000 } )
      this.airRoute.init()
      this.airRoute.show()
    },
    async unSubscribeToFleet (fleet) {
      let data, err;
      ({ data, err } = await this.$rpc.safeCallWithTimeout('UnSubscribe', { name: fleet }, 5000))
      if (err !== null) {
        const message =  err ?  err.message || 'Test error' : ''
        await this.$appErrors.showError(`${this.$appTxtErrors.errorSubscribeToFleet}: ${message}`)
        return
      }
      this.$store.commit('AirRoute/REMOVE_SUBSCRIBE', data)
      if (this.airRoute !== null) {
        await this.airRoute.destroy()
        this.airRoute = null
      }
    },
    async getGetFleetInfo () {
      const {
        data,
        err
      } = await this.$rpc.safeCallWithTimeout('Fleet.GetFleetInfo', { fleet: 'SWA2558' }, 5000)
      if (err !== null) {
        console.log('getAirRouteWayPoints err: ', err)
        this.wayPoints = null
        const message = err.message || ''
        const toast = await toastController.create({
          color: 'danger',
          position: 'middle',
          duration: 5000,
          message: `${this.$appTxtErrors.errorGetFleetInfo}: ${message}`,
          showCloseButton: true
        })
        await toast.present()
        return
      }
      this.$store.commit('AirRoute/SET_FLEET_INFO', data)
    },
    async getAirRouteWayPoints () {
      const [ss, dd] = this.airRouteSelected.split(",")
      const {
        data,
        err
      } = await this.$rpc.safeCallWithTimeout('Fleet.GetAirportsGeo', { source: ss, destination: dd }, 5000)
      if (err !== null) {
        console.log('getAirRouteWayPoints err: ', err)
        this.wayPoints = null
        const toast = await toastController.create({
          color: 'danger',
          position: 'middle',
          duration: 5000,
          message: this.$appErrors.errorGetWayPoints,
          showCloseButton: true
        })
        await toast.present()
        return
      }
      this.wayPoints = [
        { name: ss, lat: data.source.lat, lng: data.source.lng },
        { name: dd, lat: data.destination.lat, lng: data.destination.lng }
      ]
    },
    async showAirRoute () {
      if (this.airRoute !== null) {
        this.airRoute.remove()
      }
      //     constructor (viewer, store, name, wayPoints, options = {}) {

      this.airRoute = createAirRoute(this.cesiumViewer, this.$store, 'A4321', this.wayPoints, { height: 5000 } )
      this.airRoute.show()
    },
  },
  watch: {
    'airRouteSelected': async function () {
      await this.getAirRouteWayPoints()
    },
    'wayPoints': async function () {
      await this.showAirRoute()
    }
  }
}
</script>
<style lang="sass" scoped></style>
