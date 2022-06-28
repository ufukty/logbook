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

        /** @type { Object.<string, Array.<function>> } */
        this.delegates = {
            placementUpdate: [],
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
        this.delegates[event].forEach((delegate) => {
            delegate();
        });
    }

    loadTestDataset() {
        this.medium.addSection("sectionID#123");
        this.medium.addSection("sectionID#124");
        this.medium.addSection("sectionID#125");

        this.medium.addRowToSection("sectionID#123", "taskID#1");
        this.medium.addRowToSection("sectionID#123", "taskID#2");
        this.medium.addRowToSection("sectionID#123", "taskID#3");
        this.medium.addRowToSection("sectionID#123", "taskID#4");
        this.medium.addRowToSection("sectionID#123", "taskID#5");
        this.medium.addRowToSection("sectionID#123", "taskID#6");
        this.medium.addRowToSection("sectionID#123", "taskID#7");
        this.medium.addRowToSection("sectionID#123", "taskID#8");
        this.medium.addRowToSection("sectionID#123", "taskID#9");
        this.medium.addRowToSection("sectionID#123", "taskID#10");
        this.medium.addRowToSection("sectionID#123", "taskID#11");
        this.medium.addRowToSection("sectionID#123", "taskID#12");
        this.medium.addRowToSection("sectionID#123", "taskID#13");
        this.medium.addRowToSection("sectionID#123", "taskID#14");
        this.medium.addRowToSection("sectionID#123", "taskID#15");
        this.medium.addRowToSection("sectionID#123", "taskID#16");
        this.medium.addRowToSection("sectionID#123", "taskID#17");
        this.medium.addRowToSection("sectionID#123", "taskID#18");
        this.medium.addRowToSection("sectionID#123", "taskID#19");
        this.medium.addRowToSection("sectionID#123", "taskID#20");
        this.medium.addRowToSection("sectionID#123", "taskID#21");
        this.medium.addRowToSection("sectionID#123", "taskID#22");
        this.medium.addRowToSection("sectionID#123", "taskID#23");
        this.medium.addRowToSection("sectionID#123", "taskID#24");
        this.medium.addRowToSection("sectionID#123", "taskID#25");
        this.medium.addRowToSection("sectionID#123", "taskID#26");

        this.medium.addRowToSection("sectionID#124", "taskID#27");
        this.medium.addRowToSection("sectionID#124", "taskID#28");
        this.medium.addRowToSection("sectionID#124", "taskID#29");
        this.medium.addRowToSection("sectionID#124", "taskID#30");
        this.medium.addRowToSection("sectionID#124", "taskID#31");
        this.medium.addRowToSection("sectionID#124", "taskID#32");
        this.medium.addRowToSection("sectionID#124", "taskID#33");
        this.medium.addRowToSection("sectionID#124", "taskID#34");
        this.medium.addRowToSection("sectionID#124", "taskID#35");
        this.medium.addRowToSection("sectionID#124", "taskID#36");
        this.medium.addRowToSection("sectionID#124", "taskID#37");
        this.medium.addRowToSection("sectionID#124", "taskID#38");
        this.medium.addRowToSection("sectionID#124", "taskID#39");
        this.medium.addRowToSection("sectionID#124", "taskID#40");
        this.medium.addRowToSection("sectionID#124", "taskID#41");
        this.medium.addRowToSection("sectionID#124", "taskID#42");
        this.medium.addRowToSection("sectionID#124", "taskID#43");
        this.medium.addRowToSection("sectionID#124", "taskID#44");
        this.medium.addRowToSection("sectionID#124", "taskID#45");
        this.medium.addRowToSection("sectionID#124", "taskID#46");
        this.medium.addRowToSection("sectionID#124", "taskID#47");
        this.medium.addRowToSection("sectionID#124", "taskID#48");
        this.medium.addRowToSection("sectionID#124", "taskID#49");
        this.medium.addRowToSection("sectionID#124", "taskID#50");
        this.medium.addRowToSection("sectionID#124", "taskID#51");
        this.medium.addRowToSection("sectionID#124", "taskID#52");
        this.medium.addRowToSection("sectionID#124", "taskID#53");
        this.medium.addRowToSection("sectionID#124", "taskID#54");
        this.medium.addRowToSection("sectionID#124", "taskID#55");
        this.medium.addRowToSection("sectionID#124", "taskID#56");
        this.medium.addRowToSection("sectionID#124", "taskID#57");
        this.medium.addRowToSection("sectionID#124", "taskID#58");
        this.medium.addRowToSection("sectionID#124", "taskID#59");

        this.medium.addRowToSection("sectionID#125", "taskID#60");
        this.medium.addRowToSection("sectionID#125", "taskID#61");
        this.medium.addRowToSection("sectionID#125", "taskID#62");
        this.medium.addRowToSection("sectionID#125", "taskID#63");
        this.medium.addRowToSection("sectionID#125", "taskID#64");
        this.medium.addRowToSection("sectionID#125", "taskID#65");
        this.medium.addRowToSection("sectionID#125", "taskID#66");
        this.medium.addRowToSection("sectionID#125", "taskID#67");
        this.medium.addRowToSection("sectionID#125", "taskID#68");
        this.medium.addRowToSection("sectionID#125", "taskID#69");
        this.medium.addRowToSection("sectionID#125", "taskID#70");
        this.medium.addRowToSection("sectionID#125", "taskID#71");
        this.medium.addRowToSection("sectionID#125", "taskID#72");
        this.medium.addRowToSection("sectionID#125", "taskID#73");
        this.medium.addRowToSection("sectionID#125", "taskID#74");
        this.medium.addRowToSection("sectionID#125", "taskID#75");
        this.medium.addRowToSection("sectionID#125", "taskID#76");
        this.medium.addRowToSection("sectionID#125", "taskID#77");
        this.medium.addRowToSection("sectionID#125", "taskID#78");
        this.medium.addRowToSection("sectionID#125", "taskID#79");
        this.medium.addRowToSection("sectionID#125", "taskID#80");
        this.medium.addRowToSection("sectionID#125", "taskID#81");
        this.medium.addRowToSection("sectionID#125", "taskID#82");
        this.medium.addRowToSection("sectionID#125", "taskID#83");
        this.medium.addRowToSection("sectionID#125", "taskID#84");
        this.medium.addRowToSection("sectionID#125", "taskID#85");
        this.medium.addRowToSection("sectionID#125", "taskID#86");
        this.medium.addRowToSection("sectionID#125", "taskID#87");
        this.medium.addRowToSection("sectionID#125", "taskID#88");
        this.medium.addRowToSection("sectionID#125", "taskID#89");
        this.medium.addRowToSection("sectionID#125", "taskID#90");
        this.medium.addRowToSection("sectionID#125", "taskID#91");
        this.medium.addRowToSection("sectionID#125", "taskID#92");
        this.medium.addRowToSection("sectionID#125", "taskID#93");
        this.medium.addRowToSection("sectionID#125", "taskID#94");
        this.medium.addRowToSection("sectionID#125", "taskID#95");
        this.medium.addRowToSection("sectionID#125", "taskID#96");
        this.medium.addRowToSection("sectionID#125", "taskID#97");
        this.medium.addRowToSection("sectionID#125", "taskID#98");
        this.medium.addRowToSection("sectionID#125", "taskID#99");
        this.medium.addRowToSection("sectionID#125", "taskID#00");

        this.notifyDelegateFor("placementUpdate");
    }

    getTextContent(objectSymbol) {
        const objectID = pSymbol.reverse(objectSymbol);
        return `text: ${objectID}`;
    }
}
