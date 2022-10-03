import { AbstractTableViewController, TRIGGER_REPLACEMENT } from "../baja.sl/AbstactTableViewController.js";
import {
    AbstractScrollViewController,
    SCROLL_MODE_SCROLL,
    SCROLL_MODE_HIDE,
    VERTICAL,
} from "../baja.sl/AbstractScrollViewController.js";
import { AbstractViewController } from "../baja.sl/AbstractViewController.js";
import { iota, adoption, createElement, symbolizer } from "../baja.sl/utilities.js";

const CELL_TYPE_BLUEPRINT_TASK = iota();
const CELL_TYPE_BLUEPRINT_TARGET = iota();

class BlueprintOverviewSheet extends AbstractTableViewController {
    constructor() {
        super();

        this.registerCellIdentifier(CELL_TYPE_BLUEPRINT_TASK, () => {
            // FIXME: don't do this.
            // this controller won't have diffrent type of cells
            // rather, it will stylize imported cells with css (toggle classes)
        });
    }

    getDefaultHeightOfItem() {
        return 38;
    }

    getCellKindForItem() {
        return CELL_TYPE_BLUEPRINT_TARGET;
    }
}

export class BlueprintModal extends AbstractViewController {
    constructor() {
        super();
        this.dom = {
            container: createElement("div", ["bp-container"]),
            leftGridDecoration: createElement("div", ["bp-left-grid-decoration"]),
            rightGridDecoration: createElement("div", ["bp-right-grid-decoration"]),
            // backdrop: createElement("div", ["bp-backdrop"]),
            // modal: createElement("div", ["bp-modal"]),
            // header: createElement("div", ["bp-header"]),
            // overviewBgExpander: createElement("div", ["bp-overview-bg-expander"]),
            // overview: createElement("div", ["bp-overview"]),
            // footer: createElement("div", ["bp-footer"]),
        };

        this.controllers = {
            scrollView: new AbstractScrollViewController(),
            blueprintPreviewSheet: new BlueprintOverviewSheet(),
        };

        adoption(document.getElementById("blueprint-modal-mount"), [
            adoption(this.dom.container, [
                adoption(this.controllers.scrollView.dom.container, [
                    adoption(this.dom.leftGridDecoration),
                    adoption(this.controllers.blueprintPreviewSheet.container),
                    adoption(this.dom.rightGridDecoration),
                ]),
            ]),
        ]);

        this.controllers.blueprintPreviewSheet.config.placement.symbols = [
            symbolizer.symbolize("task1"),
            symbolizer.symbolize("task2"),
            symbolizer.symbolize("task3"),
        ];
        this.controllers.blueprintPreviewSheet.updateView(TRIGGER_REPLACEMENT);

        this.controllers.scrollView.config.overflowMode.vertical = SCROLL_MODE_SCROLL;
        this.controllers.scrollView.config.overflowMode.horizontal = SCROLL_MODE_HIDE;
    }
}
