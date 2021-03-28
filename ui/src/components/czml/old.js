export const xbearing = function (start, dest) {
    const startLat = bearingToRadians(start.lat)
    const startLng = bearingToRadians(start.lng)
    const destLat = bearingToRadians(dest.lng)
    const destLng = bearingToRadians(dest.lng)

    const y = Math.sin(destLng - startLng) * Math.cos(destLat)
    const x = Math.cos(startLat) * Math.sin(destLat) -
        Math.sin(startLat) * Math.cos(destLat) * Math.cos(destLng - startLng)
    let brng = Math.atan2(y, x)
    brng = bearingToDegrees(brng)
    return (brng + 360) % 360
}

export const bearingXY = function (x1, y1, x2, y2) {
    const dx = x2-x1
    const dy = y2-y1

    console.log('dx:', dx)
    console.log('dy:', dy)

    if (dx === 0) {
        return dy > 0 ? 90 : 270
    }
    if (dy === 0) {
        return dx > 0 ? 0 : 180
    }

    return Math.atan(dy/dx) * 180 / Math.PI
}
