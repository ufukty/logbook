import { DelegateRegistry } from "../DelegateRegistry.js";
import { symbolizer } from "../utilities.js";
import { Size } from "./Coordinates.js";
// import { Layout } from "./Layout.js";
import { Layout } from "./LayoutEnvironment.js";

export class AbstractLayoutPipe {
    constructor() {
        this.config = {};

        this.passedThroughPipeline = {
            /**
             * Only needed if the solid class mutates the existing layout or
             *   takes it as reference.
             */
            layout: {
                /** @type {Map.<Symbol, Area>} */
                positions: new Map(),
                /** @type {Map.<Symbol, number>} */
                scaling: new Map(),
            },
            /**  @type {Size} */
            containerSize: undefined,
        };

        /**
         * This will be automatically assigned and internally used by the
         * Environment class.
         */
        this.controlledByEnvironment = {
            /** @type {Layout} */
            environmentRef: undefined,
            /**  @type {Symbol} */
            environmentSymbol: undefined,
        };
    }

    /** @abstract */
    perform() {
        console.error("abstract function is called directly");
    }
}

export class AbstractLayoutCalculator extends AbstractLayoutPipe {
    /**
     * @param {Array.<Symbol>} placement
     * @param {Map.<Symbol, Size>} itemSizes
     */
    constructor(placement, itemSizes) {
        super();

        this.config = {
            ...this.config,
            /**
             * @type {Array.<Symbol>}
             * Incomplete-ordered list of placement data, height-ignored items
             *   should also be in this array.
             */
            placement: [],
            /**
             * States what is the actual index of items[0]
             * @type {number}
             */
            offset: 0,
            /**
             * Total number of items in the document. That value is used
             *   for estimation of full height of cell scroller for both
             *   chronological and hierarchical view.
             * @type {number}
             */
            totalNumberOfItems: undefined,
            averageSizeForUnplacedItem: new Size(0, 0),
        };

        this.computedValues = {
            ...this.computedValues,
            originItemSymbol: symbolizer.symbolize("origin"),
        };

        this._delegates = new DelegateRegistry();
    }

    perform() {
        console.error("abstract function is called directly");
    }

    update() {
        this.perform();
        if (this.computedValues.update.changed) this._delegates.nofify(POSITION_CHANGE);
    }

    /** @param {Array.<Symbol>} placement */
    newPlacement(placement) {
        console.log("newplacement", placement);
        this.config.placement = placement;
        this.controlledByEnvironment.environmentRef.scheduleRecalculation();
    }
}

export class AbstractLayoutMutator extends AbstractLayoutPipe {
    constructor() {
        this.config = {
            ...this.config,
            /** @type {Layout} */
            layout: undefined,
        };
    }

    perform() {
        console.error("abstract function is called directly.");
    }
}

export class AbstractLayoutDecorator extends AbstractLayoutPipe {
    constructor() {
        this.config = {
            items: [],
            totalNumberOfItems: 0,
            /** @type {Map.<Symbol, Anchor>} */
            anchors: new Map(),
        };
    }

    perform() {
        console.error("abstract function is called directly.");
    }

    /** @param {Array.<Symbol>} placement */
    newPlacement(placement) {
        this.placement = placement;
        this.controlledByEnvironment.environmentRef.scheduleRecalculation();
    }
}
