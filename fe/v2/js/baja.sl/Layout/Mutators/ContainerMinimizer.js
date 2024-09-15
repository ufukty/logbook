import { AbstractLayoutMutator } from "../AbstractLayoutPipe.js";
import { Position } from "../Coordinates.js";

export class ContainerMinimizer extends AbstractLayoutMutator {
    /**
     * @param {Position} notStartAfter
     * @param {Position} notEndBefore
     */
    constructor(notStartAfter = undefined, notEndBefore = undefined) {
        this.notStartAfter = notStartAfter;
        this.notEndBefore = notEndBefore;
    }

    perform() {
        var container = new Area(Infinity, Infinity, -Infinity, -Infinity);

        if (this.notStartAfter) {
            container.x0 = this.notStartAfter.x;
            container.y0 = this.notStartAfter.y;
        }
        if (this.notEndBefore) {
            container.x1 = this.notEndBefore.x;
            container.y1 = this.notEndBefore.y;
        }

        for (const [itemSymbol, item] of this.passedThroughPipeline.layout.positions) {
            if (item.x0 < container.x0) container.x0 = item.x0;
            if (item.y0 < container.y0) container.y0 = item.y0;
            if (container.x1 < item.x1) container.x1 = item.x1;
            if (container.y1 < item.y1) container.y1 = item.y1;
        }

        this.passedThroughPipeline.containerSize = container;
    }
}
