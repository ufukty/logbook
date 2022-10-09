import { AbstractLayoutCalculator, AbstractLayoutMutator } from "../AbstractLayoutPipe.js";

/** Meant to be used for `Flow.config.direction` */
export const HORIZONTAL = iota();
/** Meant to be used for `Flow.config.direction` */
export const VERTICAL = iota();

export class Wrap extends AbstractLayoutCalculator {
    constructor() {
        super();
        this.config = {
            ...this.config,
            spacing: {
                container: new Spacing(10, 0, 10),
            },
            direction: VERTICAL,
        };
    }
}
