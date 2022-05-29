import { adoption, domElementReuseCollector, createElement } from "../utilities.js";
import AbstractViewController from "./AbstractViewController.js";

function inBetween(a, b, c) {
    if (a <= b && c <= c) return true;
    else return false;
}

function checkCollusion(item_y1, item_y2, viewport_y1, viewport_y2) {
    /*
            * * * * * * *  (y1)                     * * * * * * *  (y1)                
            *           *                           *           *              
        + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + +  (y1)
        +   *           *                           *           *           +
        +   * * * * * * *  (y2)                     *           *           +
        +                                           *           *           +
        +                                           *           *           +       <=  viewport
        +                                           *           *           +
        +                 * * * * * * *  (y1)       *           *           +
        +                 *           *             *           *           +
        + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + + +  (y2)
                          *           *             *           *                              
                          * * * * * * *  (y2)       * * * * * * *  (y2)                               
    */
    if (item_y2 < viewport_y1 || item_y1 > viewport_y2)
        // if item starts after viewport ends, or item ends before viewport starts,
        // then the item is not in viewport.
        return false;
    else return true;
}

class InfiniteSheet extends AbstractViewController {
    constructor() {
        super();
        this.container = document.getElementById("infinite-sheet");
        this.contentArea = createElement("div", ["content-area"]);
        this.anchorPosition = createElement("div", ["anchor-position"]);
        adoption(this.container, [adoption(this.contentArea, [this.anchorPosition])]);

        this.state = {
            sectionHeaderTexts: [
                "section0",
                "section1",
                "section2",
                "section3",
                "section4",
                "section5",
                "section6",
                "section7",
                "section10",
                "section11",
                "section12",
                "section13",
                "section14",
                "section15",
                "section16",
                "section17",
                "section20",
                "section21",
                "section22",
                "section23",
                "section24",
                "section25",
                "section26",
                "section27",
                // "section8",
                // "section9",
            ],
            effectiveOrdering: [
                ["sec0 row1", "sec0 row2", "sec0 row3", "sec0 row4", "sec0 row5"],
                ["sec1 row1", "sec1 row2", "sec1 row3", "sec1 row4", "sec1 row5"],
                ["sec2 row1", "sec2 row2", "sec2 row3", "sec2 row4", "sec2 row5"],
                ["sec3 row1", "sec3 row2", "sec3 row3", "sec3 row4", "sec3 row5"],
                ["sec4 row1", "sec4 row2", "sec4 row3", "sec4 row4", "sec4 row5"],
                ["sec5 row1", "sec5 row2", "sec5 row3", "sec5 row4", "sec5 row5"],
                ["sec6 row1", "sec6 row2", "sec6 row3", "sec6 row4", "sec6 row5"],
                ["sec7 row1", "sec7 row2", "sec7 row3", "sec7 row4", "sec7 row5"],
                ["sec10 row1", "sec10 row2", "sec10 row3", "sec10 row4", "sec10 row5"],
                ["sec11 row1", "sec11 row2", "sec11 row3", "sec11 row4", "sec11 row5"],
                ["sec12 row1", "sec12 row2", "sec12 row3", "sec12 row4", "sec12 row5"],
                ["sec13 row1", "sec13 row2", "sec13 row3", "sec13 row4", "sec13 row5"],
                ["sec14 row1", "sec14 row2", "sec14 row3", "sec14 row4", "sec14 row5"],
                ["sec15 row1", "sec15 row2", "sec15 row3", "sec15 row4", "sec15 row5"],
                ["sec16 row1", "sec16 row2", "sec16 row3", "sec16 row4", "sec16 row5"],
                ["sec17 row1", "sec17 row2", "sec17 row3", "sec17 row4", "sec17 row5"],
                ["sec0 row1", "sec0 row2", "sec0 row3", "sec0 row4", "sec0 row5"],
                ["sec1 row1", "sec1 row2", "sec1 row3", "sec1 row4", "sec1 row5"],
                ["sec2 row1", "sec2 row2", "sec2 row3", "sec2 row4", "sec2 row5"],
                ["sec3 row1", "sec3 row2", "sec3 row3", "sec3 row4", "sec3 row5"],
                ["sec4 row1", "sec4 row2", "sec4 row3", "sec4 row4", "sec4 row5"],
                ["sec5 row1", "sec5 row2", "sec5 row3", "sec5 row4", "sec5 row5"],
                ["sec6 row1", "sec6 row2", "sec6 row3", "sec6 row4", "sec6 row5"],
                ["sec7 row1", "sec7 row2", "sec7 row3", "sec7 row4", "sec7 row5"],
                ["sec0 row1", "sec0 row2", "sec0 row3", "sec0 row4", "sec0 row5"],
                ["sec1 row1", "sec1 row2", "sec1 row3", "sec1 row4", "sec1 row5"],
                ["sec2 row1", "sec2 row2", "sec2 row3", "sec2 row4", "sec2 row5"],
                ["sec3 row1", "sec3 row2", "sec3 row3", "sec3 row4", "sec3 row5"],
                ["sec4 row1", "sec4 row2", "sec4 row3", "sec4 row4", "sec4 row5"],
                ["sec5 row1", "sec5 row2", "sec5 row3", "sec5 row4", "sec5 row5"],
                ["sec6 row1", "sec6 row2", "sec6 row3", "sec6 row4", "sec6 row5"],
                ["sec7 row1", "sec7 row2", "sec7 row3", "sec7 row4", "sec7 row5"],
            ],
        };

        this.allocatedSectionElements = {};
        this.allocatedItemElements = {};
        this.visibleHeaderElements = {};
        this.visibleRowElements = {};


        document.addEventListener("scroll", this.renderVisible.bind(this));
    }

    numberOfSections() {
        return this.state.effectiveOrdering.length;
    }

    numberOfRowsPerSection(section) {
        return this.state.effectiveOrdering[section].length;
    }

    /*
     *
     *
     *
     *
     *
     */

    getPaddingTopForTable() {
        return 100;
    }

    getHeightOfSectionHeader(section) {
        return 50;
    }

    getHeightOfElement(section, row) {
        return 50;
    }

    calculateElementBounds() {
        let itemElementBounds = [];
        let headerElementBounds = [];
        let lastPosition = 0;

        lastPosition += this.getPaddingTopForTable();

        for (let i = 0; i < this.numberOfSections(); i++) {
            const headerHeight = this.getHeightOfSectionHeader(i);
            headerElementBounds.push({ y1: lastPosition, y2: lastPosition + headerHeight });
            lastPosition += headerHeight;

            itemElementBounds.push([]);
            for (let j = 0; j < this.numberOfRowsPerSection(i); j++) {
                const itemHeight = this.getHeightOfElement(i, j);
                itemElementBounds[i].push({ y1: lastPosition, y2: lastPosition + itemHeight });
                lastPosition += itemHeight;
            }
        }

        this.itemElementBounds = itemElementBounds;
        this.headerElementBounds = headerElementBounds;
        this.lastPosition = lastPosition;
    }

    /*
     *
     *
     *
     *
     *
     */

    saveReferenceOfAllocatedSectionElement(section, element) {
        this.allocatedSectionElements[section] = element;
    }

    getReferenceOfAllocatedSectionElement(section) {
        return this.allocatedSectionElements[section];
    }

    saveReferenceOfAllocatedRowElement(section, row, element) {
        if (!this.allocatedItemElements.hasOwnProperty(section)) {
            this.allocatedItemElements[section] = {};
        }
        this.allocatedItemElements[section][row] = element;
    }

    getReferenceOfAllocatedRowElement(section, row) {
        return this.allocatedItemElements[section][row];
    }

    /*
     *
     *
     *
     *
     *
     */

    getHeaderElement(section) {
        const element = domElementReuseCollector.get("infiniteSheetHeader");
        this.saveReferenceOfAllocatedSectionElement(section, element);

        const text = this.state.sectionHeaderTexts[section];
        element.setContent(text);
        element.setPosition(this.headerElementBounds[section].y1);
        element.setData({ section: section });
    }

    releaseHeaderElement(section) {
        const element = this.getReferenceOfAllocatedSectionElement(section);
        domElementReuseCollector.free("infiniteSheetHeader", element);
    }

    getRowElement(section, row) {
        const element = domElementReuseCollector.get("infiniteSheetRow");
        this.saveReferenceOfAllocatedRowElement(section, row, element);

        const content = this.state.effectiveOrdering[section][row];
        element.setContent(content);
        element.setPosition(this.itemElementBounds[section][row].y1);
        element.setData({ section: section, row: row });
    }

    releaseRowElement(section, row) {
        const element = this.getReferenceOfAllocatedRowElement(section, row);
        domElementReuseCollector.free("infiniteSheetRow", element);
    }

    /*
     *
     *
     *
     *
     *
     */

    showHeaderOnce(section) {
        if (!this.visibleHeaderElements.hasOwnProperty(section)) {
            this.visibleHeaderElements[section] = false;
        }

        if (this.visibleHeaderElements[section]) {
            return;
        } else {
            this.getHeaderElement(section);
            this.visibleHeaderElements[section] = true;
        }
    }

    hideHeaderOnce(section) {
        if (!this.visibleHeaderElements.hasOwnProperty(section)) {
            this.visibleHeaderElements[section] = false;
        }

        if (!this.visibleHeaderElements[section]) {
            return;
        } else {
            this.releaseHeaderElement(section);
            this.visibleHeaderElements[section] = false;
        }
    }

    showRowOnce(section, row) {
        if (!this.visibleRowElements.hasOwnProperty(section)) {
            this.visibleRowElements[section] = {};
        }
        if (!this.visibleRowElements[section].hasOwnProperty(row)) {
            this.visibleRowElements[section][row] = false;
        }

        if (this.visibleRowElements[section][row]) {
            return;
        } else {
            this.getRowElement(section, row);
            this.visibleRowElements[section][row] = true;
        }
    }

    hideRowOnce(section, row) {
        if (!this.visibleRowElements.hasOwnProperty(section)) {
            this.visibleRowElements[section] = {};
        }
        if (!this.visibleRowElements[section].hasOwnProperty(row)) {
            this.visibleRowElements[section][row] = false;
        }

        if (!this.visibleRowElements[section][row]) {
            return;
        } else {
            this.releaseRowElement(section, row);
            this.visibleRowElements[section][row] = false;
        }
    }

    /*
     *
     *
     *
     *
     *
     */

    build() {
        this.renderVisible();
    }

    renderVisible() {
        this.calculateElementBounds();

        this.container.style.height = `${this.lastPosition}px`;

        const preload_area_distance = 4 * this.getHeightOfElement(0, 0);

        // get viewport coordinates (topLeft to bottomRight)
        // const viewport_x1 = window.scrollX;
        const viewport_y1 = window.scrollY - preload_area_distance;
        // const viewport_x2 = x1 + innerWidth;
        const viewport_y2 = window.scrollY + window.innerHeight + preload_area_distance;

        // list all items that are in visible area
        for (let i = 0; i < this.headerElementBounds.length; i++) {
            const item_y1 = this.headerElementBounds[i].y1;
            const item_y2 = this.headerElementBounds[i].y2;
            const shouldBeVisible = checkCollusion(item_y1, item_y2, viewport_y1, viewport_y2);
            if (shouldBeVisible) this.showHeaderOnce(i);
            else this.hideHeaderOnce(i);

            for (let j = 0; j < this.itemElementBounds[i].length; j++) {
                const item_y1 = this.itemElementBounds[i][j].y1;
                const item_y2 = this.itemElementBounds[i][j].y2;
                const shouldBeVisible = checkCollusion(item_y1, item_y2, viewport_y1, viewport_y2);
                if (shouldBeVisible) this.showRowOnce(i, j);
                else this.hideRowOnce(i, j);
            }
        }

        // call this.getRowElement(section, item) for those tasks

        // append those elements to root
    }

    /*
     *
     *
     *
     *
     *
     */

}

export default InfiniteSheet;
