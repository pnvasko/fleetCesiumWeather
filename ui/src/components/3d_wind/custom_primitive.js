import { defaultValue, defined, destroyObject, BufferUsage, ComputeCommand, ClearCommand, Color, DrawCommand, Matrix4, Pass, VertexArray, ShaderProgram, RenderState } from 'cesium'

export class CustomPrimitive {
    constructor (options) {
        this.commandType = options.commandType

        this.geometry = options.geometry
        this.attributeLocations = options.attributeLocations
        this.primitiveType = options.primitiveType

        this.uniformMap = options.uniformMap

        this.vertexShaderSource = options.vertexShaderSource
        this.fragmentShaderSource = options.fragmentShaderSource

        this.rawRenderState = options.rawRenderState
        this.framebuffer = options.framebuffer

        this.outputTexture = options.outputTexture

        this.autoClear = defaultValue(options.autoClear, false)
        this.preExecute = options.preExecute

        this.show = true
        this.commandToExecute = undefined
        this.clearCommand = undefined
        if (this.autoClear) {
            this.clearCommand = new ClearCommand({
                color: new Color(0.0, 0.0, 0.0, 0.0),
                depth: 1.0,
                framebuffer: this.framebuffer,
                pass: Pass.OPAQUE
            })
        }
    }

    createCommand (context) {
        switch (this.commandType) {
            case 'Draw': {
                return new DrawCommand({
                    owner: this,
                    vertexArray: VertexArray.fromGeometry({
                        context: context,
                        geometry: this.geometry,
                        attributeLocations: this.attributeLocations,
                        bufferUsage: BufferUsage.STATIC_DRAW,
                    }),
                    primitiveType: this.primitiveType,
                    uniformMap: this.uniformMap,
                    modelMatrix: Matrix4.IDENTITY,
                    shaderProgram: ShaderProgram.fromCache({
                        context: context,
                        attributeLocations: this.attributeLocations,
                        vertexShaderSource: this.vertexShaderSource,
                        fragmentShaderSource: this.fragmentShaderSource
                    }),
                    framebuffer: this.framebuffer,
                    renderState: RenderState.fromCache(this.rawRenderState),
                    pass: Pass.OPAQUE
                })
            }
            case 'Compute': {
                return new ComputeCommand({
                    owner: this,
                    fragmentShaderSource: this.fragmentShaderSource,
                    uniformMap: this.uniformMap,
                    outputTexture: this.outputTexture,
                    persists: true
                })
            }
        }
    }

    setGeometry (context, geometry) {
        this.geometry = geometry
        this.commandToExecute.vertexArray = VertexArray.fromGeometry({
            context: context,
            geometry: this.geometry,
            attributeLocations: this.attributeLocations,
            bufferUsage: BufferUsage.STATIC_DRAW,
        })
    }

    update (frameState) {
        if (!this.show) {
            return
        }

        if (!defined(this.commandToExecute)) {
            this.commandToExecute = this.createCommand(frameState.context)
        }

        if (defined(this.preExecute)) {
            this.preExecute()
        }

        if (defined(this.clearCommand)) {
            frameState.commandList.push(this.clearCommand)
        }

        frameState.commandList.push(this.commandToExecute)
    }

    isDestroyed () {
        return false
    }

    destroy () {
        if (defined(this.commandToExecute)) {
            this.commandToExecute.shaderProgram = this.commandToExecute.shaderProgram && this.commandToExecute.shaderProgram.destroy()
        }
        return destroyObject(this)
    }
}
