import {
    BoundingSphere,
    Cartesian2,
    Cartesian3
} from 'cesium'
import ParticleSystem from './particle_system'
import { DataProcess } from './load_data'
import { viewRectangleToLonLatRange } from './utils'

export class Wind {
    constructor (viewer, options) {
        console.log('Wind.constructor', options)
        this.viewer = viewer
        this.options = options ? options : {}

        this.viewerParameters = {
            lonRange: new Cartesian2(),
            latRange: new Cartesian2(),
            pixelSize: 0.0
        }
        this.globeBoundingSphere = new BoundingSphere(Cartesian3.ZERO, 0.99 * 6378137.0)
        console.log('Wind.constructor', this.globeBoundingSphere)
        this.updateViewerParameters()
        DataProcess.loadData().then(data => {
            this.particleSystem = new ParticleSystem(this.viewer.scene.context, data, this.options, this.viewerParameters)
            this.addPrimitives()
            this.setupEventListeners()
        })
        // this.imageryLayers = this.viewer.imageryLayers
        // this.setGlobeLayer(this.options)
    }

    addPrimitives () {
        this.viewer.scene.primitives.add(this.particleSystem.particlesComputing.primitives.calculateSpeed);
        this.viewer.scene.primitives.add(this.particleSystem.particlesComputing.primitives.updatePosition);
        this.viewer.scene.primitives.add(this.particleSystem.particlesComputing.primitives.postProcessingPosition);

        this.viewer.scene.primitives.add(this.particleSystem.particlesRendering.primitives.segments);
        this.viewer.scene.primitives.add(this.particleSystem.particlesRendering.primitives.trails);
        this.viewer.scene.primitives.add(this.particleSystem.particlesRendering.primitives.screen);
    }

    updateViewerParameters () {
        const viewRectangle = this.viewer.camera.computeViewRectangle(this.viewer.scene.globe.ellipsoid);
        const lonLatRange = viewRectangleToLonLatRange(viewRectangle);
        this.viewerParameters.lonRange.x = lonLatRange.lon.min;
        this.viewerParameters.lonRange.y = lonLatRange.lon.max;
        this.viewerParameters.latRange.x = lonLatRange.lat.min;
        this.viewerParameters.latRange.y = lonLatRange.lat.max;

        var pixelSize = this.viewer.camera.getPixelSize(
            this.globeBoundingSphere,
            this.viewer.scene.drawingBufferWidth,
            this.viewer.scene.drawingBufferHeight
        );

        if (pixelSize > 0) {
            this.viewerParameters.pixelSize = pixelSize;
        }
    }

    setGlobeLayer (userInput) {
        console.log('setGlobeLayer: ', userInput)
    }

    setupEventListeners () {
        const that = this

        this.viewer.camera.moveStart.addEventListener(function () {
            that.viewer.scene.primitives.show = false
        })

        this.viewer.camera.moveEnd.addEventListener(function () {
            that.updateViewerParameters()
            that.particleSystem.applyViewerParameters(that.viewerParameters)
            that.viewer.scene.primitives.show = true
        })

        let resized = false
        window.addEventListener("resize", function () {
            resized = true
            that.viewer.scene.primitives.show = false
            that.viewer.scene.primitives.removeAll()
        })

        this.viewer.scene.preRender.addEventListener(function () {
            if (resized) {
                that.particleSystem.canvasResize(that.viewer.scene.context)
                resized = false
                that.addPrimitives()
                that.viewer.scene.primitives.show = true
            }
        })

        window.addEventListener('particleSystemOptionsChanged', function () {
            that.particleSystem.applyUserInput(that.options)
        })

        window.addEventListener('layerOptionsChanged', function () {
            that.setGlobeLayer(that.options)
        })
    }
}
