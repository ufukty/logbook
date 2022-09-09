export class AbstractViewController {
    constructor() {
        this.dom = {};
        this.state = {};
        this.config = {};
    }

    buildView() {
        console.error("Abstract function has called directly.");
    }
}
