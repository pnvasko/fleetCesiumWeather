import { ArcType, Cartesian3 } from 'cesium'

export const bearingToRadians = function (degrees) {
    return degrees * Math.PI / 180
}

export const bearingToDegrees = function (radians) {
    return radians * 180 / Math.PI
}

export const bearing = function (start, dest) {
    const x1 = start.lat
    const y1 = start.lng
    const x2 = dest.lat
    const y2 = dest.lng
    const dy = bearingToRadians(y2-y1)
    const y = Math.sin(dy) * Math.cos(bearingToRadians(x2))
    const x = Math.cos(bearingToRadians(x1))*Math.sin(bearingToRadians(x2)) - Math.sin(bearingToRadians(x1))*Math.cos(bearingToRadians(x2))*Math.cos(dy)
    const az = ((Math.atan2(y, x) * 180 / Math.PI) + 360) % 360
    return bearingToRadians(az)
}

export function parallel (viewer, latitude, color, granularity) {
    var name = "Parallel " + latitude;
    return viewer.entities.add({
        name: name,
        polyline: {
            positions: Cartesian3.fromDegreesArray([
                -180,
                latitude,
                -90,
                latitude,
                0,
                latitude,
                90,
                latitude,
                180,
                latitude
            ]),
            width: 2,
            arcType: ArcType.RHUMB,
            material: color,
            granularity: granularity,
        },
    });
}

export function meridian (viewer, longitude, color, granularity) {
    var name = "Meridian " + longitude;
    return viewer.entities.add({
        name: name,
        polyline: {
            positions: Cartesian3.fromDegreesArray([
                longitude,
                90,
                longitude,
                0,
                longitude,
                -90,
            ]),
            width: 2,
            arcType: ArcType.RHUMB,
            material: color,
            granularity: granularity,
        },
    });
}
