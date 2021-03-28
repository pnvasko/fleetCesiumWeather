<template>
  <ion-content :fullscreen="true">
    <div id="cesiumContainer" style="width:calc(100vw);height:calc(100vh);"></div>
    <div class="cesium_custom_toolbars">
      <ion-button size="small">Test</ion-button>
    </div>
  </ion-content>
</template>

<script>
import 'cesium/Build/Cesium/Widgets/widgets.css'
import { IonContent, IonButton } from '@ionic/vue'
import {Viewer, EllipsoidTerrainProvider, Clock, UrlTemplateImageryProvider, OpenStreetMapImageryProvider} from 'cesium'
import { defaultSystemOptions, defaultLayerOptions, Wind } from '@/components/3d_wind'

export default  {
  name: '3DWindHome',
  components: { IonContent, IonButton },
  data () {
    return {
      earth: null,
      cesiumViewer: null,
      wind: null,
      maxParticles: defaultSystemOptions.maxParticles,
      particleHeight: defaultSystemOptions.particleHeight,
      fadeOpacity: defaultSystemOptions.fadeOpacity,
      dropRate: defaultSystemOptions.dropRate,
      dropRateBump: defaultSystemOptions.dropRateBump,
      speedFactor: defaultSystemOptions.speedFactor,
      lineWidth: defaultSystemOptions.lineWidth,
      globeLayer: defaultLayerOptions.globeLayer,
      WMS_URL: defaultLayerOptions.WMS_URL
    }
  },
  mounted() {
    // window.CESIUM_BASE_URL = JSON.stringify('')
    window.CESIUM_BASE_URL = 'http://192.168.3.254:8080/'
    this.initMap()
  },
  methods: {
    getUserInput () {
      const particlesTextureSize = Math.ceil(Math.sqrt(this.maxParticles))
      this.maxParticles = particlesTextureSize * particlesTextureSize

      return {
        particlesTextureSize: particlesTextureSize,
        maxParticles: this.maxParticles,
        particleHeight: this.particleHeight,
        fadeOpacity: this.fadeOpacity,
        dropRate: this.dropRate,
        dropRateBump: this.dropRateBump,
        speedFactor: this.speedFactor,
        lineWidth: this.lineWidth,
        globeLayer: this.globeLayer,
        WMS_URL: this.WMS_URL
      }
    },
    initMap () {
      // https://cesium.com/docs/tutorials/optimizing-quotas/
      // https://sandcastle.cesium.com/index.html?src=Offline.html
      // https://cesium.com/docs/cesiumjs-ref-doc/TileMapServiceImageryProvider.html?classFilter=TileMap
      // ./wms-get-linux-amd64 --url https://ows.terrestris.de/osm/service --layer OSM-WMS --zooms 0,1,2,3,4,5,6 --bbox -180,-88,180,88
      const imageryProvider = new OpenStreetMapImageryProvider({
        xurl: 'https://a.tile.openstreetmap.org/',
        url: '/maps'
      })
      this.cesiumViewer = new Viewer("cesiumContainer", {
        imageryProvider: imageryProvider,
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
      this.wind = new Wind(this.cesiumViewer, this.getUserInput())
      // this.changeBaseMap("gg")
    },
    changeBaseMap () {
      this.cesiumViewer.imageryLayers.removeAll()
      const url = 'http://mt1.google.cn/vt/lyrs=s&hl=zh-CN&x={x}&y={y}&z={z}&s=Gali'
      this.cesiumViewer.imageryLayers.addImageryProvider(
              new UrlTemplateImageryProvider({ url: url })
      )
    }
  }
}
</script>
<style lang="sass" scoped>
  @import "../css/variables.sass"
  #cesiumContainer
    width: 100%
    height: 100%

  .cesium_custom_toolbars
    position: absolute
    top: 10px
    left: 10px
    width: 200px
    height: 200px
    background-color: rgba($primary-dark, .1)
    border: 1px solid $dark
    border-radius: 5px
    .button-native
      background-color: $dark
</style>
