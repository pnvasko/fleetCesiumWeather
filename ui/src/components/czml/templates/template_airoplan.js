export const czmlAiroplanTemplate = function (id, position) {
    return {
        id: `airoplan-${id}`,
        position: {
            cartographicDegrees: [position.lng, position.lat, 5],
        },
        orientation: {
            interpolationAlgorithm: "LAGRANGE",
            unitQuaternion: [ 0.0, 0.0, 0.0, 1.0 ]
        },
        model: {
            interpolationAlgorithm: "LAGRANGE",
            gltf: "models/flightradar24/b788.glb",
            minimumPixelSize: 50,
            scale: 0.001
        },
        shadows: "ENABLED"
    }
}
