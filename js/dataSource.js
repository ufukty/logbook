import { symbolizer, pick } from "./baja.sl/utilities.js";
import { DelegateRegistry } from "./baja.sl/DelegateRegistry.js";
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
             * @type {Map.<string, {parentId: string,depth: number,degree: number,isCollaborated: bool,isTarget: bool, isCompleted: bool}>}
             */
            tasks: new Map(),
            updateCounts: new Map(),
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
        this.cache.placements.totalNumberOfItems = 50;
        var dayCounter = 1;
        const symbols = new Map();
        const taskIndexes = Array.from(Array(50).keys());
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
            this.cache.tasks.set(taskSymbol, {
                content: taskId,
                isCompleted: pick([false]),
                isCollaborated: pick([true, false]),
                isTarget: pick([true, false]),
                depth: pick(Array.from(Array(10).keys())),
            });
            this.cache.updateCounts.set(taskSymbol, 0);
        });
        shuffle(this.cache.placements.hierarchical.items);
        this.delegates.nofify(EVENT_PLACEMENT_UPDATE);

        const samples = [
            "Lorem ipsum dolor sit amet consectetur adipisicing elit",
            "Unde minus quod eius quas blanditiis incidunt dolores cum quibusdam eos nostrum",
            "Architecto possimus totam deleniti doloribus sint consequuntur inventore nobis perferendis",
            "Facilis debitis consequuntur aliquid distinctio porro ",
            "Iure nam nemo mollitia nesciunt blanditiis corrupti maiores cumque quia autem dolore molestiae doloremque sunt culpa laudantium alias",
            "Quidem debitis ab eligendi adipisci earum",
            "Consectetur",
            "Suscipit eveniet ",
            "Recusandae in eaque laborum blanditiis rem ducimus",
            "Culpa consequuntur ",
            "Excepturi optio qui commodi sequi laboriosam eos nostrum sapiente aliquam",
            "Doloremque vero possimus",
            "Sed deleniti",
            "Distinctio nobis mollitia",
            "Sed quidem fugit",
            "Libero quibusdam laboriosam soluta ipsa quasi magni esse",
            "Iste labore at inventore aut optio dignissimos ipsam quod pariatur",
            "Iure atque distinctio nisi id nihil minima animi ",
            "Dolor",
            "Voluptatem",
            "Porro dolorum qui quos error maxime reprehenderit excepturi vero laboriosam soluta voluptatibus praesentium nam quisquam eaque",
            "Voluptas cupiditate odio",
            "Debitis consectetur sit Mollitia voluptatem",
            "Debitis esse voluptas iste itaque",
            "Esse unde quam ex perspiciatis laboriosam totam numquam ea",
            "Cupiditate repellat voluptas facere accusantium illum quas magnam cumque pariatur maiores dolor voluptatibus et vitae nesciunt nihil ",
            "Laboriosam atque similique quaerat",
            "Consequuntur cupiditate obcaecati ea quisquam facilis",
            "Quae ad non quod quam tempora odit repellat illum eveniet possimus molestiae dolor",
            "Dolorum vitae natus ",
            "Praesentium ipsam a atque recusandae Dolore",
            "Quam cum",
            "Est saepe a at exercitationem vitae",
            "Earum rerum perferendis voluptas accusamus assumenda et asperiores veniam dolore blanditiis ",
            "Saepe",
            "Voluptate est cum quibusdam harum aspernatur quod ",
            "Harum",
            "Expedita Veritatis",
            "Aliquid id",
            "Cupiditate quos impedit quibusdam excepturi Tenetur eius",
            "Enim blanditiis saepe necessitatibus nostrum illum vel ipsum quaerat ab voluptatem",
            "Totam deserunt voluptatibus facere",
            "Sapiente fugiat ",
            "Laudantium doloribus cupiditate fugiat vero iste",
            "Rerum doloremque",
            "Exercitationem totam suscipit placeat necessitatibus error cupiditate ",
            "Sed quia labore nostrum quisquam culpa eligendi velit possimus error enim quis deserunt expedita architecto incidunt",
            "Tenetur officia repellendus ",
            "Qui",
            "Velit",
        ];

        setInterval(() => {
            const itemSymbol = pick(this.cache.placements.hierarchical.items);
            var task = this.cache.tasks.get(itemSymbol);
            const nextUpdateCount = pick([this.cache.updateCounts.get(itemSymbol) + 1, 0]);
            this.cache.updateCounts.set(itemSymbol, nextUpdateCount);
            task.content = `${symbolizer.desymbolize(itemSymbol)} ${pick(samples)}`;
            this.delegates.nofify(EVENT_OBJECT_UPDATE, new Set([itemSymbol]));
        }, 100);
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
