import { AbstractLayoutPipe } from "./AbstractLayoutPipe.js";
import { AbstractLayoutCalculator, AbstractLayoutDecorator, AbstractLayoutMutator } from "./AbstractLayoutPipe.js";
import { Size } from "./Coordinates.js";
import { Layout } from "./Layout.js";

export class LayoutEnvironment {
    constructor() {
        this.internalUsage = {
            /** @type {Array.<AbstractLayoutPipe>} */
            pipeline: [],
            /** @type {Array.<Function>} */
            subscriber: [],
        };

        this.passedThroughPipeline = {
            /** @type {Layout} */
            layout: new Layout(),
            /**  @type {Size} */
            containerSize: undefined,
        };

        this.environmentSymbol = symbolizer.symbolize(`environment#${iota()}`);
    }

    refreshPipeline() {
        for (const pipe of this.internalUsage.pipeline) {
            pipe.passedThroughPipeline = this.passedThroughPipeline;
            pipe.perform();
            this.passedThroughPipeline = pipe.passedThroughPipeline;
        }
    }

    /**
     * @private
     * @param {AbstractLayoutPipe} pipe
     */
    _connectPipeToPipeline(pipe) {
        this.pipeline.push(pipe);
        pipe.controlledByEnvironment.environmentSymbol = this.environmentSymbol;
        pipe.controlledByEnvironment.environmentRef = this;
    }

    /** @param {AbstractLayoutCalculator} pipe */
    connectCalculator(pipe) {
        if (!pipe instanceof AbstractLayoutCalculator)
            console.error("Given pipe is not an instance of AbstractLayoutCalculator class");
        this._connectPipeToPipeline(pipe);
        return this;
    }

    /** @param {AbstractLayoutMutator} pipe */
    connectMutator(pipe) {
        if (!pipe instanceof AbstractLayoutMutator)
            console.error("Given pipe is not an instance of AbstractLayoutMutator class");
        this._connectPipeToPipeline(pipe);
        return this;
    }

    /** @param {AbstractLayoutDecorator} pipe */
    connectDecorator(pipe) {
        if (!pipe instanceof AbstractLayoutDecorator)
            console.error("Given pipe is not an instance of AbstractLayoutDecorator class");
        this._connectPipeToPipeline(pipe);
        return this;
    }

    /**
     *
     */
    start() {}

    /**
     * It is important to call this method when refreshing position calculation
     * for each resizeObserver notification is not needed anymore.
     */
    stop() {}
}
