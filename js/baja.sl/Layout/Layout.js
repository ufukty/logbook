import { Position } from "./Coordinates.js";

export class Layout {
    constructor() {
        /** @type {Map.<Symbol, Position>} */
        this.positions = new Map();
        /** @type {Map.<Symbol, number>} */
        this.scaling = new Map();
    }
}
