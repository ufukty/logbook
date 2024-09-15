import { symbolizer } from "../baja.sl/utilities.js";

export class InfiniteSheetDataMedium {
    constructor() {
        /** @type {{sections: Array.Symbol, rows: Map.<Symbol, Symbol> }} */
        
        /** @type { Map.<Symbol, Symbol> } */
        this.mapRowSection = new Map();
    }

    /**
     * @param {string} sectionID - A string represents specific section. Could  be a UUID.
     * @param {number=} index - Optional argument. Don't pass this if you want to add a new section to the end.
     */
    addSection(sectionID, index = undefined) {
        const sectionSymbol = symbolizer.get(sectionID);
        // append to a specific index (if specified) or to the end.
        if (index === undefined) this.data.sections.push(sectionSymbol);
        else this.data.sections.splice(index, 0, sectionSymbol);
        this.data.rows.set(sectionSymbol, []);
    }

    deleteSection(sectionID) {
        // TODO:
    }

    moveSection(sectionID, newIndex) {
        // TODO:
        const sectionSymbol = symbolizer.get(sectionID);
    }

    /** arguments are given as strings, not symbols */
    addRowToSection(sectionID, rowID, index = undefined) {
        const sectionSymbol = symbolizer.get(sectionID);
        const rowSymbol = symbolizer.get(rowID);

        /** @type {[]} */
        const placement = this.data.rows.get(sectionSymbol);
        // append to a specific index (if specified) or to the end.
        if (index === undefined) placement.push(rowSymbol);
        else placement.splice(index, 0, rowSymbol);

        this.mapRowSection.set(rowSymbol, sectionSymbol);
    }

    moveRow(rowID, newIndex) {
        const rowSymbol = symbolizer.get(rowID);
        const sectionSymbol = this.mapRowSection.get(rowSymbol);

        // remove row from old index
        /** @type {[]} */
        const placement = this.data.rows.get(sectionSymbol);
        const index = placement.indexOf(rowSymbol);
        placement.splice(index, 1);

        // add to new index
        placement.splice(newIndex, 0, rowSymbol);
    }

    deleteRow(rowID) {
        const rowSymbol = symbolizer.get(rowID);
        const sectionSymbol = this.mapRowSection.get(rowSymbol);

        // remove row from old index
        /** @type {[]} */
        const placement = this.data.rows.get(sectionSymbol);
        const index = placement.indexOf(rowSymbol);
        placement.splice(index, 1);

        this.mapRowSection.delete(rowSymbol);
    }

    moveRowToAnotherSection(rowID, newSectionID, newIndex) {
        const rowSymbol = symbolizer.get(rowID);
        const currentSectionSymbol = this.mapRowSection.get(rowSymbol);
        const nextSectionSymbol = symbolizer.get(newSectionID);

        // remove row from old section & index
        /** @type {[]} */
        const currentSectionPlacement = this.data.rows.get(currentSectionSymbol);
        const index = currentSectionPlacement.indexOf(rowSymbol);
        currentSectionPlacement.splice(index, 1);

        // add to new section & index
        /** @type {[]} */
        const nextSectionPlacement = this.data.rows.get(nextSectionSymbol);
        nextSectionPlacement.splice(newIndex, 0, rowSymbol);

        this.mapRowSection.set(rowSymbol, nextSectionSymbol);
    }
}
