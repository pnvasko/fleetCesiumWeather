import {
    Math,
} from 'cesium'

export function randomizeParticles (ldata, maxParticles, viewerParameters) {
    let array = new Float32Array(4 * maxParticles)

    for (var i = 0; i < maxParticles; i++) {
        array[4 * i] = Math.randomBetween(viewerParameters.lonRange.x, viewerParameters.lonRange.y)
        array[4 * i + 1] = Math.randomBetween(viewerParameters.latRange.x, viewerParameters.latRange.y)
        array[4 * i + 2] = Math.randomBetween(ldata.lev.min, ldata.lev.max)
        array[4 * i + 3] = 0.0
    }

    return array
}
