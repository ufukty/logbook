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
        this.cache.placements.totalNumberOfItems = 1000;

        const symbols = {
            "task#123": symbolizer.symbolize("task#123"),
            "task#143": symbolizer.symbolize("task#143"),
            "task#133": symbolizer.symbolize("task#133"),
            "task#124": symbolizer.symbolize("task#124"),
            "task#144": symbolizer.symbolize("task#144"),
            "task#134": symbolizer.symbolize("task#134"),
            "task#125": symbolizer.symbolize("task#125"),
            "task#145": symbolizer.symbolize("task#145"),
            "task#126": symbolizer.symbolize("task#126"),
            "task#146": symbolizer.symbolize("task#146"),
            "task#135": symbolizer.symbolize("task#135"),
            "task#127": symbolizer.symbolize("task#127"),
            "task#147": symbolizer.symbolize("task#147"),
            "task#136": symbolizer.symbolize("task#136"),
            "task#137": symbolizer.symbolize("task#137"),
            "day#0": symbolizer.symbolize("day#0"),
            "day#1": symbolizer.symbolize("day#1"),
            "day#2": symbolizer.symbolize("day#2"),
            "day#3": symbolizer.symbolize("day#3"),
        };

        this.cache.placements.chronological = {
            headerSymbols: [symbols["day#0"], symbols["day#1"], symbols["day#2"], symbols["day#3"]],
            offset: 0,
            items: [
                symbols["day#0"],
                symbols["task#123"],
                symbols["task#124"],
                symbols["task#125"],
                symbols["task#126"],
                symbols["day#1"],
                symbols["task#127"],
                symbols["task#133"],
                symbols["task#134"],
                symbols["task#135"],
                symbols["day#2"],
                symbols["task#136"],
                symbols["task#137"],
                symbols["task#143"],
                symbols["task#144"],
                symbols["day#3"],
                symbols["task#145"],
                symbols["task#146"],
                symbols["task#147"],
            ],
        };

        this.cache.placements.hierarchical = {
            offset: 0,
            items: [
                symbols["task#123"],
                symbols["task#143"],
                symbols["task#133"],
                symbols["task#124"],
                symbols["task#144"],
                symbols["task#134"],
                symbols["task#125"],
                symbols["task#145"],
                symbols["task#126"],
                symbols["task#146"],
                symbols["task#135"],
                symbols["task#127"],
                symbols["task#147"],
                symbols["task#136"],
                symbols["task#137"],
            ],
        };

        // prettier-ignore
        this.cache.tasks.set(symbols["task#123"], { content: "task#123", parentId: symbols["-1"] });
        this.cache.tasks.set(symbols["task#143"], { content: "task#143", parentId: symbols["task#123"] });
        this.cache.tasks.set(symbols["task#144"], { content: "task#144", parentId: symbols["task#143"] });
        this.cache.tasks.set(symbols["task#126"], { content: "task#126", parentId: symbols["task#144"] });
        this.cache.tasks.set(symbols["task#146"], { content: "task#146", parentId: symbols["task#144"] });
        this.cache.tasks.set(symbols["task#134"], { content: "task#134", parentId: symbols["task#143"] });
        this.cache.tasks.set(symbols["task#125"], { content: "task#125", parentId: symbols["task#143"] });
        this.cache.tasks.set(symbols["task#137"], { content: "task#137", parentId: symbols["task#125"] });
        this.cache.tasks.set(symbols["task#145"], { content: "task#145", parentId: symbols["task#143"] });
        this.cache.tasks.set(symbols["task#135"], { content: "task#135", parentId: symbols["task#145"] });
        this.cache.tasks.set(symbols["task#133"], { content: "task#133", parentId: symbols["task#123"] });
        this.cache.tasks.set(symbols["task#127"], { content: "task#127", parentId: symbols["task#133"] });
        this.cache.tasks.set(symbols["task#147"], { content: "task#147", parentId: symbols["task#133"] });
        this.cache.tasks.set(symbols["task#136"], { content: "task#136", parentId: symbols["task#133"] });
        this.cache.tasks.set(symbols["task#124"], { content: "task#124", parentId: symbols["task#123"] });

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
