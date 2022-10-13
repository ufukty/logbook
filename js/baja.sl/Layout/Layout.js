import { Area } from "./Coordinates.js";

/**
 * @typedef {Symbol} ItemSymbol
 * @typedef {Symbol} CellTypeSymbol
 * @typedef {Symbol} ViewControllerSymbol
 */

export class Layout {
    constructor() {
        /** @type {Map.<Symbol, Area>} */
        this.positions = new Map();
        /** @type {Map.<Symbol, number>} */
        this.scaling = new Map();
    }
}
