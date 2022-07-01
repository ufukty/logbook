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
            objectUpdate: [],
        };
    }

    notifyDelegateFor(event, ...args) {
        this.delegates[event].forEach((delegate) => {
            delegate(...args);
        });
    }

    loadTestDataset() {
        this.medium.addSection("sectionID#123");
        this.medium.addSection("sectionID#124");
        this.medium.addSection("sectionID#125");

        this.rowSections = new Map([
            ["taskID#1", "sectionID#123"],
            ["taskID#2", "sectionID#123"],
            ["taskID#3", "sectionID#123"],
            ["taskID#4", "sectionID#123"],
            ["taskID#5", "sectionID#123"],
            ["taskID#6", "sectionID#123"],
            ["taskID#7", "sectionID#123"],
            ["taskID#8", "sectionID#123"],
            ["taskID#9", "sectionID#123"],
            ["taskID#10", "sectionID#123"],
            ["taskID#11", "sectionID#123"],
            ["taskID#12", "sectionID#123"],
            ["taskID#13", "sectionID#123"],
            ["taskID#14", "sectionID#123"],
            ["taskID#15", "sectionID#123"],
            ["taskID#16", "sectionID#123"],
            ["taskID#17", "sectionID#123"],
            ["taskID#18", "sectionID#123"],
            ["taskID#19", "sectionID#123"],
            ["taskID#20", "sectionID#123"],
            ["taskID#21", "sectionID#123"],
            ["taskID#22", "sectionID#123"],
            ["taskID#23", "sectionID#123"],
            ["taskID#24", "sectionID#123"],
            ["taskID#25", "sectionID#123"],
            ["taskID#26", "sectionID#123"],
            ["taskID#27", "sectionID#124"],
            ["taskID#28", "sectionID#124"],
            ["taskID#29", "sectionID#124"],
            ["taskID#30", "sectionID#124"],
            ["taskID#31", "sectionID#124"],
            ["taskID#32", "sectionID#124"],
            ["taskID#33", "sectionID#124"],
            ["taskID#34", "sectionID#124"],
            ["taskID#35", "sectionID#124"],
            ["taskID#36", "sectionID#124"],
            ["taskID#37", "sectionID#124"],
            ["taskID#38", "sectionID#124"],
            ["taskID#39", "sectionID#124"],
            ["taskID#40", "sectionID#124"],
            ["taskID#41", "sectionID#124"],
            ["taskID#42", "sectionID#124"],
            ["taskID#43", "sectionID#124"],
            ["taskID#44", "sectionID#124"],
            ["taskID#45", "sectionID#124"],
            ["taskID#46", "sectionID#124"],
            ["taskID#47", "sectionID#124"],
            ["taskID#48", "sectionID#124"],
            ["taskID#49", "sectionID#124"],
            ["taskID#50", "sectionID#124"],
            ["taskID#51", "sectionID#124"],
            ["taskID#52", "sectionID#124"],
            ["taskID#53", "sectionID#124"],
            ["taskID#54", "sectionID#124"],
            ["taskID#55", "sectionID#124"],
            ["taskID#56", "sectionID#124"],
            ["taskID#57", "sectionID#124"],
            ["taskID#58", "sectionID#124"],
            ["taskID#59", "sectionID#124"],
            ["taskID#60", "sectionID#125"],
            ["taskID#61", "sectionID#125"],
            ["taskID#62", "sectionID#125"],
            ["taskID#63", "sectionID#125"],
            ["taskID#64", "sectionID#125"],
            ["taskID#65", "sectionID#125"],
            ["taskID#66", "sectionID#125"],
            ["taskID#67", "sectionID#125"],
            ["taskID#68", "sectionID#125"],
            ["taskID#69", "sectionID#125"],
            ["taskID#70", "sectionID#125"],
            ["taskID#71", "sectionID#125"],
            ["taskID#72", "sectionID#125"],
            ["taskID#73", "sectionID#125"],
            ["taskID#74", "sectionID#125"],
            ["taskID#75", "sectionID#125"],
            ["taskID#76", "sectionID#125"],
            ["taskID#77", "sectionID#125"],
            ["taskID#78", "sectionID#125"],
            ["taskID#79", "sectionID#125"],
            ["taskID#80", "sectionID#125"],
            ["taskID#81", "sectionID#125"],
            ["taskID#82", "sectionID#125"],
            ["taskID#83", "sectionID#125"],
            ["taskID#84", "sectionID#125"],
            ["taskID#85", "sectionID#125"],
            ["taskID#86", "sectionID#125"],
            ["taskID#87", "sectionID#125"],
            ["taskID#88", "sectionID#125"],
            ["taskID#89", "sectionID#125"],
            ["taskID#90", "sectionID#125"],
            ["taskID#91", "sectionID#125"],
            ["taskID#92", "sectionID#125"],
            ["taskID#93", "sectionID#125"],
            ["taskID#94", "sectionID#125"],
            ["taskID#95", "sectionID#125"],
            ["taskID#96", "sectionID#125"],
            ["taskID#97", "sectionID#125"],
            ["taskID#98", "sectionID#125"],
            ["taskID#99", "sectionID#125"],
            ["taskID#00", "sectionID#125"],
        ]);

        this.objectContents = new Map([
            ["sectionID#123", "text content for sectionID#123"],
            ["sectionID#124", "text content for sectionID#124"],
            ["sectionID#125", "text content for sectionID#125"],
            ["taskID#1", "text content for taskID#1"],
            ["taskID#2", "text content for taskID#2"],
            ["taskID#3", "text content for taskID#3"],
            ["taskID#4", "text content for taskID#4"],
            ["taskID#5", "text content for taskID#5"],
            ["taskID#6", "text content for taskID#6"],
            ["taskID#7", "text content for taskID#7"],
            ["taskID#8", "text content for taskID#8"],
            ["taskID#9", "text content for taskID#9"],
            ["taskID#10", "text content for taskID#10"],
            ["taskID#11", "text content for taskID#11"],
            ["taskID#12", "text content for taskID#12"],
            ["taskID#13", "text content for taskID#13"],
            ["taskID#14", "text content for taskID#14"],
            ["taskID#15", "text content for taskID#15"],
            ["taskID#16", "text content for taskID#16"],
            ["taskID#17", "text content for taskID#17"],
            ["taskID#18", "text content for taskID#18"],
            ["taskID#19", "text content for taskID#19"],
            ["taskID#20", "text content for taskID#20"],
            ["taskID#21", "text content for taskID#21"],
            ["taskID#22", "text content for taskID#22"],
            ["taskID#23", "text content for taskID#23"],
            ["taskID#24", "text content for taskID#24"],
            ["taskID#25", "text content for taskID#25"],
            ["taskID#26", "text content for taskID#26"],

            ["taskID#27", "text content for taskID#27"],
            ["taskID#28", "text content for taskID#28"],
            ["taskID#29", "text content for taskID#29"],
            ["taskID#30", "text content for taskID#30"],
            ["taskID#31", "text content for taskID#31"],
            ["taskID#32", "text content for taskID#32"],
            ["taskID#33", "text content for taskID#33"],
            ["taskID#34", "text content for taskID#34"],
            ["taskID#35", "text content for taskID#35"],
            ["taskID#36", "text content for taskID#36"],
            ["taskID#37", "text content for taskID#37"],
            ["taskID#38", "text content for taskID#38"],
            ["taskID#39", "text content for taskID#39"],
            ["taskID#40", "text content for taskID#40"],
            ["taskID#41", "text content for taskID#41"],
            ["taskID#42", "text content for taskID#42"],
            ["taskID#43", "text content for taskID#43"],
            ["taskID#44", "text content for taskID#44"],
            ["taskID#45", "text content for taskID#45"],
            ["taskID#46", "text content for taskID#46"],
            ["taskID#47", "text content for taskID#47"],
            ["taskID#48", "text content for taskID#48"],
            ["taskID#49", "text content for taskID#49"],
            ["taskID#50", "text content for taskID#50"],
            ["taskID#51", "text content for taskID#51"],
            ["taskID#52", "text content for taskID#52"],
            ["taskID#53", "text content for taskID#53"],
            ["taskID#54", "text content for taskID#54"],
            ["taskID#55", "text content for taskID#55"],
            ["taskID#56", "text content for taskID#56"],
            ["taskID#57", "text content for taskID#57"],
            ["taskID#58", "text content for taskID#58"],
            ["taskID#59", "text content for taskID#59"],

            ["taskID#60", "text content for taskID#60"],
            ["taskID#61", "text content for taskID#61"],
            ["taskID#62", "text content for taskID#62"],
            ["taskID#63", "text content for taskID#63"],
            ["taskID#64", "text content for taskID#64"],
            ["taskID#65", "text content for taskID#65"],
            ["taskID#66", "text content for taskID#66"],
            ["taskID#67", "text content for taskID#67"],
            ["taskID#68", "text content for taskID#68"],
            ["taskID#69", "text content for taskID#69"],
            ["taskID#70", "text content for taskID#70"],
            ["taskID#71", "text content for taskID#71"],
            ["taskID#72", "text content for taskID#72"],
            ["taskID#73", "text content for taskID#73"],
            ["taskID#74", "text content for taskID#74"],
            ["taskID#75", "text content for taskID#75"],
            ["taskID#76", "text content for taskID#76"],
            ["taskID#77", "text content for taskID#77"],
            ["taskID#78", "text content for taskID#78"],
            ["taskID#79", "text content for taskID#79"],
            ["taskID#80", "text content for taskID#80"],
            ["taskID#81", "text content for taskID#81"],
            ["taskID#82", "text content for taskID#82"],
            ["taskID#83", "text content for taskID#83"],
            ["taskID#84", "text content for taskID#84"],
            ["taskID#85", "text content for taskID#85"],
            ["taskID#86", "text content for taskID#86"],
            ["taskID#87", "text content for taskID#87"],
            ["taskID#88", "text content for taskID#88"],
            ["taskID#89", "text content for taskID#89"],
            ["taskID#90", "text content for taskID#90"],
            ["taskID#91", "text content for taskID#91"],
            ["taskID#92", "text content for taskID#92"],
            ["taskID#93", "text content for taskID#93"],
            ["taskID#94", "text content for taskID#94"],
            ["taskID#95", "text content for taskID#95"],
            ["taskID#96", "text content for taskID#96"],
            ["taskID#97", "text content for taskID#97"],
            ["taskID#98", "text content for taskID#98"],
            ["taskID#99", "text content for taskID#99"],
            ["taskID#00", "text content for taskID#00"],
        ]);

        for (const [rowID, sectionID] of this.rowSections.entries()) {
            this.medium.addRowToSection(sectionID, rowID);
        }

        this.notifyDelegateFor("placementUpdate");

        const prom = new Promise((resolve, reject) => {
            setTimeout(resolve, 1000);
        });

        prom.then(() => {
            this.objectContents.set(
                "taskID#1",
                "Lorem ipsum dolor sit amet consectetur adipisicing elit. Omnis voluptatum labore in hic possimus dolor. Aliquam tempore unde quia natus hic optio modi excepturi. Reprehenderit natus recusandae dolores rerum omnis?"
            );
            // this.notifyDelegateFor("placementUpdate");
            this.notifyDelegateFor("objectUpdate", new Set([pSymbol.get("taskID#1")]));
        }).then(() => {
            console.log("updated");
        });
    }

    getTextContent(objectSymbol) {
        const objectID = pSymbol.reverse(objectSymbol);
        return this.objectContents.get(objectID);
    }
}
