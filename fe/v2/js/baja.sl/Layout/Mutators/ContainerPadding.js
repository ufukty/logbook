import { AbstractLayoutMutator } from "../AbstractLayoutPipe.js";

export class ContainerPadding extends AbstractLayoutMutator {
    /**
     * @param {number} top
     * @param {number} left
     * @param {number} right
     * @param {number} bottom
     */
    constructor(top = undefined, left = undefined, right = undefined, bottom = undefined) {
        super();

        this.top = top ?? 0;
        this.left = left ?? 0;
        this.right = right ?? 0;
        this.bottom = bottom ?? 0;
    }

    perform() {
        for (const position of this.passedThroughPipeline.layout.positions.values()) {
            position.x += this.left;
            position.y += this.top;
        }
        this.passedThroughPipeline.containerSize.height += this.top + this.bottom;
        this.passedThroughPipeline.containerSize.width += this.left + this.right;
    }
}
