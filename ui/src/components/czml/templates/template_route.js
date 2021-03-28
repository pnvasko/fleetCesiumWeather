export const czmlRouteTemplate = function (id, source, destination, height, color, width) {
    return {
        id: `route-${id}`,
        name: `${id} route`,
        polyline: {
            positions: {
                cartographicDegrees: [source.lng, source.lat, height, destination.lng, destination.lat, height],
            },
            material: {
                solidColor: {
                    color: {
                        rgba: color.toBytes(),
                    },
                },
            },
            width: width,
            clampToGround: true,
        },
    }
}

// import { Color } from 'cesium'

export const _czmlRouteTemplate = function (id, source, destination, height, color, width) {
    // const solidColor = Color.toBytes(color)
    console.log('solidColor: ', color.toBytes())
    return {
        id: `route-${id}`,
        name: `${id} route`,
        polyline: {
            positions: {
                cartographicDegrees: [source.lng, source.lat, height, destination.lng, destination.lat, height],
            },
            material: {
                solidColor: {
                    color: [255,0,0,255] //color.toBytes(),
                },
            },
            width: width,
            clampToGround: true
        }
    }
}

