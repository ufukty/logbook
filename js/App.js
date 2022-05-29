import { executeWhenDocumentIsReady } from "./utilities.js";

import ModeSelector from "./viewControllers/ModeSelector.js";
import InfiniteSheet from "./viewControllers/InfiniteSheet.js";

class App {
    constructor() {
        this.state = {
            documentMode: undefined,
        };

        this.modeSelector = new ModeSelector(this.updateMode.bind(this));
        this.infiniteSheet = new InfiniteSheet();
    }

    updateMode(newMode) {
        this.state.documentMode = newMode;
        this.updateView();
    }

    updateView() {
        this.infiniteSheet.build();
    }

    build() {
        this.infiniteSheet.build();
    }
}

executeWhenDocumentIsReady(function () {
    const app = new App();
    app.build();
    // const body = document.getElementsByTagName("body")[0];
    // body.appendChild(app.container)
});
