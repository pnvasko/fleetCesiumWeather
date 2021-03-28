export const RedLine = [
    {
        id: "document",
        name: "CZML Geometries: Polyline",
        version: "1.0",
    },
    {
        id: "redLine",
        name: "Red line clamped to terain",
        polyline: {
            positions: {
                cartographicDegrees: [-75, 35, 9000, -125, 35, 9000],
            },
            material: {
                solidColor: {
                    color: {
                        rgba: [255, 0, 0, 255],
                    },
                },
            },
            width: 5,
            clampToGround: true,
        },
    }
]
