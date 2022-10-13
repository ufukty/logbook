export class UpdateScheduler {
    /**
     * @param  {Function} updateFunction
     * @param {number} hertz - Number of maximum updates allowed each second, that has same period between each
     */
    constructor(updateFunction, hertz = 60) {
        this.updateFunction = updateFunction;
        this.updateIsOngoing = false;
        this.lastUpdateTime = undefined;
        /** @private */
        this.period = 1000 / hertz;
        /** @private */
        this.scheduled = false;
    }

    schedule(...params) {
        console.log("schedule");
        const now = Date.now();

        if (this.lastUpdateTime) {
            const timePassedSinceLastUpdate = now - this.lastUpdateTime;
            const remainingToNextUpdate = this.period - timePassedSinceLastUpdate;

            if (this.updateIsOngoing || timePassedSinceLastUpdate <= this.period) {
                if (!this.scheduled) {
                    this.scheduled = true;
                    setTimeout(() => {
                        this.updateFunction(...params);
                    }, remainingToNextUpdate + 1);
                }
                return;
            }

            if (this.scheduled) {
                this.scheduled = false;
            }
        }

        this.lastUpdateTime = now;
        this.updateFunction(...params);
    }

    finished() {
        this.updateIsOngoing = false;
    }
}
