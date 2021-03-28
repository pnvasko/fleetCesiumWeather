import netcdfjs from 'netcdfjs'

export const DataProcess = (function () {
    let data

    const loadNetCDF = function (filePath) {
        return new Promise(function (resolve, reject) {
            const request = new XMLHttpRequest()
            request.open('GET', filePath)
            request.responseType = 'arraybuffer'

            request.onerror = function (e) {
                console.log('loadNetCDF request.onerror:', e)
                reject(new TypeError('Network request failed'))
            }

            request.onload = function () {
                const arrayToMap = function (array) {
                    return array.reduce(function (map, object) {
                        map[object.name] = object
                        return map
                    }, {})
                }

                const NetCDF = new netcdfjs(request.response)

                data = {}
                const dimensions = arrayToMap(NetCDF.dimensions);
                data.dimensions = {}
                data.dimensions.lon = dimensions['lon'].size
                data.dimensions.lat = dimensions['lat'].size
                data.dimensions.lev = dimensions['lev'].size

                const variables = arrayToMap(NetCDF.variables)
                const uAttributes = arrayToMap(variables['U'].attributes)
                const vAttributes = arrayToMap(variables['V'].attributes)

                data.lon = {}
                data.lon.array = new Float32Array(NetCDF.getDataVariable('lon').flat())
                data.lon.min = Math.min(...data.lon.array)
                data.lon.max = Math.max(...data.lon.array)

                data.lat = {}
                data.lat.array = new Float32Array(NetCDF.getDataVariable('lat').flat())
                data.lat.min = Math.min(...data.lat.array)
                data.lat.max = Math.max(...data.lat.array)

                data.lev = {}
                data.lev.array = new Float32Array(NetCDF.getDataVariable('lev').flat())
                data.lev.min = Math.min(...data.lev.array)
                data.lev.max = Math.max(...data.lev.array)

                data.U = {}
                data.U.array = new Float32Array(NetCDF.getDataVariable('U').flat())
                data.U.min = uAttributes['min'].value
                data.U.max = uAttributes['max'].value

                data.V = {}
                data.V.array = new Float32Array(NetCDF.getDataVariable('V').flat())
                data.V.min = vAttributes['min'].value
                data.V.max = vAttributes['max'].value

                resolve(data)
            }

            request.send()
        })
    }

    const loadData = async function () {
        // var ncFilePath = fileOptions.dataDirectory + fileOptions.dataFile;
        const ncFilePath = '../data/demo.nc'
        await loadNetCDF(ncFilePath)
        return data
    }

    return {
        loadData: loadData
    }
})()
