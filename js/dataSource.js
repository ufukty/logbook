import { TableViewPlacementData } from "./tableViewPlacementData.js";

export class DataSource {
    constructor() {
        this.placement = new TableViewPlacementData();

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
        this.placement.addSection("sectionID#123");
        this.placement.addSection("sectionID#124");
        this.placement.addSection("sectionID#125");

        this.placement.addRowToSection("sectionID#123", "taskID#234");
        this.placement.addRowToSection("sectionID#123", "taskID#235");
        this.placement.addRowToSection("sectionID#124", "taskID#236");
        this.placement.addRowToSection("sectionID#124", "taskID#237");
        this.placement.addRowToSection("sectionID#125", "taskID#238");
        this.placement.addRowToSection("sectionID#125", "taskID#239");

        this.notifyDelegateFor("placementUpdate");
    }

    getTextContent(objectSymbol) {
        const objectID = pSymbol.reverse(objectSymbol);
        return this.textContent[objectID];
    }
}
