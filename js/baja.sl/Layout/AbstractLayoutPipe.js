import { DelegateRegistry } from "../DelegateRegistry.js";
import { Size } from "./Coordinates.js";
import { Layout } from "./Layout.js";
import { LayoutEnvironment } from "./LayoutEnvironment.js";

export class AbstractLayoutPipe {
    constructor() {
        this.config = {};

        this.passedThroughPipeline = {
            /** @type {Layout} - Only needed if the solid class mutates the existing layout or takes it as reference. */
            layout: undefined,
            /**  @type {Size} */
            containerSize: undefined,
        };

        /**
         * This will be automatically assigned and internally used by the
         * Environment class.
         */
        this.controlledByEnvironment = {
            /** @type {LayoutEnvironment} */
            environmentRef: undefined,
            /**  @type {Symbol} */
            environmentSymbol: undefined,
        };
    }

    _getTemplateForComputedValues() {
        return {
            itemChangedPosition: new Set(),
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
            /** @type {Array.<Symbol>} */
            placement: [],
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
        this.placement = placement;
        this.controlledByEnvironment.environmentRef.refreshPipeline();
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
        this.controlledByEnvironment.environmentRef.refreshPipeline();
    }
}
