export class NoOptionException extends Error {
    constructor(message) {
        super(message)
        this.name = "NoOptionException"
    }
}
