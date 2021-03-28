import { Color } from 'cesium'

function newRouteStyle(route, airoportBg, airoportBorder, traceRoute, shadowAiroplan, width, twidth, radius) {
    return {
        Route:  Color.fromCssColorString(route),
        RouteWidth: width,
        AiroportBg: Color.fromCssColorString(airoportBg),
        AiroportBorder: Color.fromCssColorString(airoportBorder),
        AiroportRadius: radius,
        TraceRoute: Color.fromCssColorString(traceRoute),
        TraceRouteWidth: twidth,
        ShadowAiroplan: Color.fromCssColorString(shadowAiroplan),
    }
}

export const routeColor = {
    Route:  Color.fromCssColorString("#F5F5F5"),
    AiroportBg: Color.fromCssColorString("#F5F5F5"),
    AiroportBorder: Color.fromCssColorString("#F5F5F5"),
    TraceRoute: Color.fromCssColorString("#F5F5F5"),
    ShadowAiroplan: Color.fromCssColorString("#F5F5F5")
}

const airRouteColors = [
    Object.freeze(newRouteStyle('#4bc701', '#3c9a02', '#205f02', '#8d0232', 'rgba(141,2,50,0.25)', 3, 5, 10))
]

export default function (i) {
    return i >= airRouteColors.length ? airRouteColors[0] : airRouteColors[i]
}
