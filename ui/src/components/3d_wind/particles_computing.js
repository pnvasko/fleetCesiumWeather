import { Cartesian2, Cartesian3, Sampler, PixelFormat, PixelDatatype, TextureMinificationFilter, TextureMagnificationFilter, ShaderSource } from 'cesium'
import { CustomPrimitive } from './custom_primitive'
import { randomizeParticles } from './data_process'
import { createTexture } from './utils'
import { calculateSpeed, updatePosition, postProcessingPosition } from './glsl'

export default class ParticlesComputing {
    constructor (context, data, userInput, viewerParameters) {
        this.data = data
        this.createWindTextures(context, data)
        this.createParticlesTextures(context, userInput, viewerParameters)
        this.createComputingPrimitives(data, userInput, viewerParameters)
    }

    createWindTextures (context, data) {
        const windTextureOptions = {
            context: context,
            width: data.dimensions.lon,
            height: data.dimensions.lat * data.dimensions.lev,
            pixelFormat: PixelFormat.LUMINANCE,
            pixelDatatype: PixelDatatype.FLOAT,
            flipY: false,
            sampler: new Sampler({
                // the values of texture will not be interpolated
                minificationFilter: TextureMinificationFilter.NEAREST,
                magnificationFilter: TextureMagnificationFilter.NEAREST
            })
        }

        this.windTextures = {
            U: createTexture(windTextureOptions, data.U.array),
            V: createTexture(windTextureOptions, data.V.array)
        }
    }

    createParticlesTextures (context, userInput, viewerParameters) {
        const particlesTextureOptions = {
            context: context,
            width: userInput.particlesTextureSize,
            height: userInput.particlesTextureSize,
            pixelFormat: PixelFormat.RGBA,
            pixelDatatype: PixelDatatype.FLOAT,
            flipY: false,
            sampler: new Sampler({
                // the values of texture will not be interpolated
                minificationFilter: TextureMinificationFilter.NEAREST,
                magnificationFilter: TextureMagnificationFilter.NEAREST
            })
        };

        const particlesArray = randomizeParticles(this.data, userInput.maxParticles, viewerParameters)
        const zeroArray = new Float32Array(4 * userInput.maxParticles).fill(0)

        this.particlesTextures = {
            previousParticlesPosition: createTexture(particlesTextureOptions, particlesArray),
            currentParticlesPosition: createTexture(particlesTextureOptions, particlesArray),
            nextParticlesPosition: createTexture(particlesTextureOptions, particlesArray),
            postProcessingPosition: createTexture(particlesTextureOptions, particlesArray),

            particlesSpeed: createTexture(particlesTextureOptions, zeroArray)
        }
    }

    destroyParticlesTextures() {
        Object.keys(this.particlesTextures).forEach((key) => {
            this.particlesTextures[key].destroy()
        })
    }

    createComputingPrimitives (data, userInput, viewerParameters) {
        const dimension = new Cartesian3(data.dimensions.lon, data.dimensions.lat, data.dimensions.lev)
        const minimum = new Cartesian3(data.lon.min, data.lat.min, data.lev.min)
        const maximum = new Cartesian3(data.lon.max, data.lat.max, data.lev.max)
        const interval = new Cartesian3(
            (maximum.x - minimum.x) / (dimension.x - 1),
            (maximum.y - minimum.y) / (dimension.y - 1),
            dimension.z > 1 ? (maximum.z - minimum.z) / (dimension.z - 1) : 1.0
        )
        const uSpeedRange = new Cartesian2(data.U.min, data.U.max)
        const vSpeedRange = new Cartesian2(data.V.min, data.V.max)

        const that = this

        this.primitives = {
            calculateSpeed: new CustomPrimitive({
                commandType: 'Compute',
                uniformMap: {
                    U: function () {
                        return that.windTextures.U
                    },
                    V: function () {
                        return that.windTextures.V
                    },
                    currentParticlesPosition: function () {
                        return that.particlesTextures.currentParticlesPosition
                    },
                    dimension: function () {
                        return dimension
                    },
                    minimum: function () {
                        return minimum
                    },
                    maximum: function () {
                        return maximum
                    },
                    interval: function () {
                        return interval
                    },
                    uSpeedRange: function () {
                        return uSpeedRange
                    },
                    vSpeedRange: function () {
                        return vSpeedRange
                    },
                    pixelSize: function () {
                        return viewerParameters.pixelSize
                    },
                    speedFactor: function () {
                        return userInput.speedFactor
                    }
                },
                fragmentShaderSource: new ShaderSource({
                    sources: [calculateSpeed]
                }),
                outputTexture: this.particlesTextures.particlesSpeed,
                preExecute: function () {
                    // swap textures before binding
                    const temp = that.particlesTextures.previousParticlesPosition
                    that.particlesTextures.previousParticlesPosition = that.particlesTextures.currentParticlesPosition
                    that.particlesTextures.currentParticlesPosition = that.particlesTextures.postProcessingPosition
                    that.particlesTextures.postProcessingPosition = temp

                    // keep the outputTexture up to date
                    that.primitives.calculateSpeed.commandToExecute.outputTexture = that.particlesTextures.particlesSpeed
                }
            }),

            updatePosition: new CustomPrimitive({
                commandType: 'Compute',
                uniformMap: {
                    currentParticlesPosition: function () {
                        return that.particlesTextures.currentParticlesPosition
                    },
                    particlesSpeed: function () {
                        return that.particlesTextures.particlesSpeed
                    }
                },
                fragmentShaderSource: new ShaderSource({
                    sources: [updatePosition]
                }),
                outputTexture: this.particlesTextures.nextParticlesPosition,
                preExecute: function () {
                    // keep the outputTexture up to date
                    that.primitives.updatePosition.commandToExecute.outputTexture = that.particlesTextures.nextParticlesPosition
                }
            }),

            postProcessingPosition: new CustomPrimitive({
                commandType: 'Compute',
                uniformMap: {
                    nextParticlesPosition: function () {
                        return that.particlesTextures.nextParticlesPosition
                    },
                    particlesSpeed: function () {
                        return that.particlesTextures.particlesSpeed
                    },
                    lonRange: function () {
                        return viewerParameters.lonRange
                    },
                    latRange: function () {
                        return viewerParameters.latRange
                    },
                    randomCoefficient: function () {
                        return Math.random()
                    },
                    dropRate: function () {
                        return userInput.dropRate
                    },
                    dropRateBump: function () {
                        return userInput.dropRateBump
                    }
                },
                fragmentShaderSource: new ShaderSource({
                    sources: [postProcessingPosition]
                }),
                outputTexture: this.particlesTextures.postProcessingPosition,
                preExecute: function () {
                    // keep the outputTexture up to date
                    that.primitives.postProcessingPosition.commandToExecute.outputTexture = that.particlesTextures.postProcessingPosition
                }
            })
        }
    }
}
