// export const czmlAirportTemplate = function (id, position, height, width, color) {

export const czmlAirportTemplate = function (id, position, height, width, color, colorBg) {
    return {
        id: `airport-${id}`,
        name: `${id} airport`,
        label: {
            fillColor: {
                rgba: [255, 255, 255, 255],
            },
            font: "12pt Lucida Console",
            horizontalOrigin: "LEFT",
            pixelOffset: {
                cartesian2: [14, 14],
            },
            style: "FILL",
            text: id,
            showBackground: true,
            backgroundColor: {
                rgba: [112, 89, 57, 200],
            },
            position: {
                cartographicDegrees: [position.lng, position.lat, 5],
            }
        },
        point: {
            color: {
                rgba: color.toBytes(),
            },
            outlineColor: {
                rgba: colorBg.toBytes(),
            },
            outlineWidth: 2,
            pixelSize: width,
        },
        position: {
            cartographicDegrees: [position.lng, position.lat, height],
        }
    }
}
