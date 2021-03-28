import {
    Appearance,
    ComponentDatatype,
    Geometry,
    GeometryAttribute,
    GeometryAttributes,
    defined,
    Framebuffer,
    Math,
    Texture
} from 'cesium'

export function viewRectangleToLonLatRange(viewRectangle) {
    const range = {}
    const postiveWest = Math.mod(viewRectangle.west, Math.TWO_PI)
    const postiveEast = Math.mod(viewRectangle.east, Math.TWO_PI)
    const width = viewRectangle.width

    let longitudeMin
    let longitudeMax;
    if (width > Math.THREE_PI_OVER_TWO) {
        longitudeMin = 0.0;
        longitudeMax = Math.TWO_PI;
    } else {
        if (postiveEast - postiveWest < width) {
            longitudeMin = postiveWest;
            longitudeMax = postiveWest + width;
        } else {
            longitudeMin = postiveWest;
            longitudeMax = postiveEast;
        }
    }
    range.lon = {
        min: Math.toDegrees(longitudeMin),
        max: Math.toDegrees(longitudeMax)
    }
    const south = viewRectangle.south
    const north = viewRectangle.north
    const height = viewRectangle.height

    let extendHeight = height > Math.PI / 12 ? height / 2 : 0
    let extendedSouth = Math.clampToLatitudeRange(south - extendHeight)
    let extendedNorth = Math.clampToLatitudeRange(north + extendHeight)

    if (extendedSouth < -Math.PI_OVER_THREE) {
        extendedSouth = -Math.PI_OVER_TWO;
    }
    if (extendedNorth > Math.PI_OVER_THREE) {
        extendedNorth = Math.PI_OVER_TWO;
    }

    range.lat = {
        min: Math.toDegrees(extendedSouth),
        max: Math.toDegrees(extendedNorth)
    }

    return range
}

export function createTexture (options, typedArray) {
    if (defined(typedArray)) {
        // typed array needs to be passed as source option, this is required by Cesium.Texture
        const source = {}
        source.arrayBufferView = typedArray;
        options.source = source
    }

    return new Texture(options)
}

export function createFramebuffer (context, colorTexture, depthTexture) {
    return new Framebuffer({
        context: context,
        colorTextures: [colorTexture],
        depthTexture: depthTexture
    })
}

export function createRawRenderState (options) {
    const translucent = true
    const closed = false
    const existing = {
        viewport: options.viewport,
        depthTest: options.depthTest,
        depthMask: options.depthMask,
        blending: options.blending
    }

    return Appearance.getDefaultRenderState(translucent, closed, existing)
}

export function  getFullscreenQuad () {
    return new Geometry({
        attributes: new GeometryAttributes({
            position: new GeometryAttribute({
                componentDatatype: ComponentDatatype.FLOAT,
                componentsPerAttribute: 3,
                //  v3----v2
                //  |     |
                //  |     |
                //  v0----v1
                values: new Float32Array([
                    -1, -1, 0, // v0
                    1, -1, 0, // v1
                    1, 1, 0, // v2
                    -1, 1, 0, // v3
                ])
            }),
            st: new GeometryAttribute({
                componentDatatype: ComponentDatatype.FLOAT,
                componentsPerAttribute: 2,
                values: new Float32Array([
                    0, 0,
                    1, 0,
                    1, 1,
                    0, 1,
                ])
            })
        }),
        indices: new Uint32Array([3, 2, 0, 0, 2, 1])
    })
}
