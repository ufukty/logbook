import { pSymbol } from "./utilities.js";

export class TableViewStructuredDataMedium {
    constructor() {
        /** @type {{sections: Array.Symbol, rows: Map.<Symbol, Symbol> }} */
        this.data = {
            sections: [],
            rows: new Map(),
        };
        /** @type { Map.<Symbol, Symbol> } */
        this.mapRowSection = new Map();
    }

    /**
     * @param {string} sectionID - A string represents specific section. Could  be a UUID.
     * @param {number=} index - Optional argument. Don't pass this if you want to add a new section to the end.
     */
    addSection(sectionID, index = undefined) {
        const sectionSymbol = pSymbol.get(sectionID);
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
        const sectionSymbol = pSymbol.get(sectionID);
    }

    /** arguments are given as strings, not symbols */
    addRowToSection(sectionID, rowID, index = undefined) {
        const sectionSymbol = pSymbol.get(sectionID);
        const rowSymbol = pSymbol.get(rowID);

        /** @type {[]} */
        const placement = this.data.rows.get(sectionSymbol);
        // append to a specific index (if specified) or to the end.
        if (index === undefined) placement.push(rowSymbol);
        else placement.splice(index, 0, rowSymbol);

        this.mapRowSection.set(rowSymbol, sectionSymbol);
    }

    moveRow(rowID, newIndex) {
        const rowSymbol = pSymbol.get(rowID);
        const sectionSymbol = this.mapRowSection.get(rowSymbol);

        // remove row from old index
        /** @type {[]} */
        const placement = this.data.rows.get(sectionSymbol);
        const index = placement.indexOf(rowID);
        placement.splice(index, 1);

        // add to new index
        placement.splice(newIndex, 0, rowID);
    }

    deleteRow(rowID) {
        const rowSymbol = pSymbol.get(rowID);
        const sectionSymbol = this.mapRowSection.get(rowSymbol);

        // remove row from old index
        /** @type {[]} */
        const placement = this.data.rows.get(sectionSymbol);
        const index = placement.indexOf(rowID);
        placement.splice(index, 1);

        this.mapRowSection.delete(rowSymbol);
    }

    moveRowToAnotherSection(rowID, newSectionID, newIndex) {
        const rowSymbol = pSymbol.get(rowID);
        const currentSectionSymbol = this.mapRowSection.get(rowSymbol);
        const nextSectionSymbol = pSymbol.get(newSectionID);

        // remove row from old section & index
        /** @type {[]} */
        const currentSectionPlacement = this.data.rows.get(currentSectionSymbol);
        const index = currentSectionPlacement.indexOf(rowID);
        currentSectionPlacement.splice(index, 1);

        // add to new section & index
        /** @type {[]} */
        const nextSectionPlacement = this.data.rows.get(nextSectionSymbol);
        nextSectionPlacement.splice(newIndex, 0, rowID);

        this.mapRowSection.set(rowSymbol, nextSectionSymbol);
    }
}

export class DataSource {
    constructor() {
        this.medium = new TableViewStructuredDataMedium();

        this.delegates = {
            placementUpdate: undefined,
        };

        this.textContent = {
            "sectionID#123": "section 123 header",
            "sectionID#124": "section 124 header",
            "sectionID#125": "section 125 header",
            "taskID#234": "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
            "taskID#235": "Vivamus vitae nibh nec tortor porta congue quis eu ante",
            "taskID#236": "Aliquam rhoncus tortor nec elit molestie, a rutrum odio hendrerit",
            "taskID#237": "Donec varius feugiat purus id sagittis",
            "taskID#238": "Nunc quis fringilla tellus, sed aliquam dui",
            "taskID#239": "Proin molestie dolor eget purus molestie, a cursus mi iaculis",
        };
    }

    notifyDelegateFor(event) {
        this.delegates[event]();
    }

    loadTestDataset() {
        this.medium.addSection("sectionID#123");
        this.medium.addSection("sectionID#124");
        this.medium.addSection("sectionID#125");

        this.medium.addRowToSection("sectionID#123", "taskID#234");
        this.medium.addRowToSection("sectionID#123", "taskID#235");
        this.medium.addRowToSection("sectionID#124", "taskID#236");
        this.medium.addRowToSection("sectionID#124", "taskID#237");
        this.medium.addRowToSection("sectionID#125", "taskID#238");
        this.medium.addRowToSection("sectionID#125", "taskID#239");

        this.notifyDelegateFor("placementUpdate");
    }

    getTextContent(objectSymbol) {
        const objectID = pSymbol.reverse(objectSymbol);
        return this.textContent[objectID];
    }
}
