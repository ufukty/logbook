export class AbstractViewController {
    constructor() {
        this.dom = {};
        this.state = {};
        this.config = {
            debug: false,
        };
    }

    _error(...data) {
        if (this.config.debug) console.error(...data);
    }

    _debug(...data) {
        if (this.config.debug) console.log(...data);
    }

    buildView() {
        console.error("Abstract function has called directly.");
    }
}
