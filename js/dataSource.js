import { symbolizer } from "./bjsl/utilities.js";
import { DelegateRegistry } from "./bjsl/DelegateRegistry.js";
import { InfiniteSheetDataMedium } from "./viewControllers/InfiniteSheetDataMedium.js";

function fetchRetry(url, delay, tries, options) {
    if (tries > 0)
        return fetch(url, options).catch(() => {
            setTimeout(() => {
                fetchRetry(url, delay, tries - 1, options);
            }, delay);
        });
    else return fetch(url, options);
}

/** This event occurs when placement of tasks is changed. */
export const EVENT_PLACEMENT_UPDATE = "EVENT_PLACEMENT_UPDATE";

/** This event occurs when an object's content is changed.
 * Callbacks registered as delegate to this event should
 * accept list of objectIds as argument */
export const EVENT_OBJECT_UPDATE = "EVENT_OBJECT_UPDATE";

function shuffle(array) {
    let currentIndex = array.length,
        randomIndex;

    // While there remain elements to shuffle.
    while (currentIndex != 0) {
        // Pick a remaining element.
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex--;

        // And swap it with the current element.
        [array[currentIndex], array[randomIndex]] = [array[randomIndex], array[currentIndex]];
    }

    return array;
}

export class DataSource {
    constructor() {
        this.medium = new InfiniteSheetDataMedium();
        this.delegates = new DelegateRegistry([EVENT_OBJECT_UPDATE, EVENT_PLACEMENT_UPDATE]);

        this.config = {
            network: {
                apiUrl: "https://localhost:8082",
                endpoints: {
                    task: {
                        fold: "/task/fold",
                        unfold: "/task/unfold",
                        create: "/task/create",
                        delete: "/task/delete",
                        move: "/task/move",
                        type: "/task/type",
                    },
                },
                /** period of time to wait before try again for failed requests */
                delay: 500,
                /** period of time to wait before try again for failed requests */
                retryCount: 10,
            },
        };

        /** Acts like a cache. Remove invalidated data immediately. */
        this.cache = {
            placements: {
                chronological: {
                    /** Incomplete list of placement data.
                     * @type {Array.<string>} */
                    items: [],
                    /** States what is the actual index of items[0]
                     * @type {number} */
                    offset: undefined,
                    headerSymbols: [],
                },
                hierarchical: {
                    /** Incomplete list of placement data.
                     * @type {Array.<string>} */
                    items: [],
                    /** States what is the actual index of items[0]
                     * @type {number} */
                    offset: undefined,
                },
                /** Total number of items in the document. That value is used
                 * for estimation of full height of cell scroller for both
                 * chronological and hierarchical view.
                 * @type {number} */
                totalNumberOfItems: undefined,
            },
            /** Remove invalidated tasks immediatelly after servers confirm the
             * modification. Remove LRU keys to keep memory usage constant.
             * @type {Map.<string,{parentId: string, depth: number, degree: number}>} */
            tasks: new Map(),
        };

        this.computedData = {
            serializedChronologicalPlacement: [], // only indexes -200 <-> +200 scroll position
        };
    }

    /**
     * This function will
     * @param {string} taskId
     * @returns {Promise.<Response>}
     **/
    fold(taskId) {
        return fetchRetry("https://localhost:8080", this.config.network.delay, this.config.network.retryCount, {
            method: "UPDATE",
            headers: { "content-type": "application/json" },
            body: JSON.stringify({
                taskId: taskId,
            }),
        }).then((result) => {
            return result.json();
        });
    }

    /**
     * @param {number} focusedTaskIndex
     * @param {number} offset number of tasks plus/minus focusedTaskIndex
     * @example
     * .getSerializeChronologicalPlacement(0, 100).then((json) => {
     *   console.log(json)
     * });
     * @returns {Promise} JSON in a promise.
     */
    getSerializeChronologicalPlacement(focusedTaskIndex = 0, offset = 100) {
        // TODO: seralize the part requested, attach section headers
        // append to the cache
        // and return

        const lowerBound = focusedTaskIndex > offset ? focusedTaskIndex - offset : 0;
        const url = `${apiRootURL}/document/placement/${documentId}?offset=${lowerBound}&limit=200`;
        fetch(url)
            .then((response) => {
                console.log(response);
                return response.json();
            })
            .then((json) => {
                console.log(json);
            });

        for (let i = lowerBound; i < upperBound; i++) {
            const chunkIndex = i - (i % 100);
            const key = this.upToDateFetchData.placements.chronological.tasks.get();
        }
    }

    loadTestDataset() {
        this.cache.placements.totalNumberOfItems = 100;
        var dayCounter = 1;
        const symbols = new Map();
        const taskIndexes = Array.from(Array(100).keys());
        taskIndexes.forEach((taskIndex, index) => {
            if (index === 0 || Math.random() > 0.8) {
                // add header
                const dayId = "day#" + dayCounter.toString();
                const daySymbol = symbolizer.symbolize(dayId);
                symbols.set(dayId, daySymbol);
                this.cache.placements.chronological.headerSymbols.push(daySymbol);
                this.cache.placements.chronological.items.push(daySymbol);
                dayCounter += 1;
            }
            const taskId = "task#" + taskIndex.toString();
            const taskSymbol = symbolizer.symbolize(taskId);
            symbols.set(taskId, taskSymbol);
            this.cache.placements.hierarchical.items.push(taskSymbol);
            this.cache.placements.chronological.items.push(taskSymbol);
            this.cache.tasks.set(taskSymbol, { content: taskId });
        });
        shuffle(this.cache.placements.hierarchical.items);
        this.delegates.nofify(EVENT_PLACEMENT_UPDATE);
    }

    getTextContent(objectSymbol) {
        // const objectID = symbolizer.desymbolize(objectSymbol);
        // const match = objectID.match(/section/)
        // if (match.length > 0)
        if (
            this.cache.placements.chronological.headerSymbols.findIndex((symbol) => {
                return symbol === objectSymbol;
            }) != -1
        ) {
            return symbolizer.desymbolize(objectSymbol);
        } else {
            return this.cache.tasks.get(objectSymbol).content;
        }
    }
}
