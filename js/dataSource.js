import { pSymbol } from "./bjsl/utilities.js";
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

export class DataSource {
    constructor() {
        this.medium = new InfiniteSheetDataMedium();

        /** @type { Object.<string, Array.<function>> } */
        this.delegates = {
            placementUpdate: [],
            objectUpdate: [],
        };

        this.config = {
            network: {
                apiUrl: "https://localhost:8080",
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

    notifyDelegateFor(event, ...args) {
        this.delegates[event].forEach((delegate) => {
            delegate(...args);
        });
    }

    loadTestDataset() {
        this.medium.addSection("sectionID#123");
        this.medium.addSection("sectionID#124");
        this.medium.addSection("sectionID#125");
        this.medium.addSection("sectionID#126");

        this.rowSections = new Map([
            ["taskID#1", "sectionID#123"],
            ["taskID#2", "sectionID#123"],
            ["taskID#3", "sectionID#123"],
            ["taskID#4", "sectionID#123"],
            ["taskID#5", "sectionID#123"],
            ["taskID#6", "sectionID#123"],
            ["taskID#7", "sectionID#123"],
            ["taskID#8", "sectionID#123"],
            ["taskID#9", "sectionID#123"],
            ["taskID#10", "sectionID#123"],
            ["taskID#11", "sectionID#123"],
            ["taskID#12", "sectionID#123"],
            ["taskID#13", "sectionID#123"],
            ["taskID#14", "sectionID#123"],
            ["taskID#15", "sectionID#123"],
            ["taskID#16", "sectionID#123"],
            ["taskID#17", "sectionID#123"],
            ["taskID#18", "sectionID#123"],
            ["taskID#19", "sectionID#123"],
            ["taskID#20", "sectionID#123"],
            ["taskID#21", "sectionID#123"],
            ["taskID#22", "sectionID#123"],
            ["taskID#23", "sectionID#123"],
            ["taskID#24", "sectionID#123"],
            ["taskID#25", "sectionID#123"],
            ["taskID#26", "sectionID#123"],
            ["taskID#27", "sectionID#124"],
            ["taskID#28", "sectionID#124"],
            ["taskID#29", "sectionID#124"],
            ["taskID#30", "sectionID#124"],
            ["taskID#31", "sectionID#124"],
            ["taskID#32", "sectionID#124"],
            ["taskID#33", "sectionID#124"],
            ["taskID#34", "sectionID#124"],
            ["taskID#35", "sectionID#124"],
            ["taskID#36", "sectionID#124"],
            ["taskID#37", "sectionID#124"],
            ["taskID#38", "sectionID#124"],
            ["taskID#39", "sectionID#124"],
            ["taskID#40", "sectionID#124"],
            ["taskID#41", "sectionID#124"],
            ["taskID#42", "sectionID#124"],
            ["taskID#43", "sectionID#124"],
            ["taskID#44", "sectionID#124"],
            ["taskID#45", "sectionID#124"],
            ["taskID#46", "sectionID#124"],
            ["taskID#47", "sectionID#124"],
            ["taskID#48", "sectionID#124"],
            ["taskID#49", "sectionID#124"],
            ["taskID#50", "sectionID#124"],
            ["taskID#51", "sectionID#124"],
            ["taskID#52", "sectionID#124"],
            ["taskID#53", "sectionID#124"],
            ["taskID#54", "sectionID#124"],
            ["taskID#55", "sectionID#124"],
            ["taskID#56", "sectionID#124"],
            ["taskID#57", "sectionID#124"],
            ["taskID#58", "sectionID#124"],
            ["taskID#59", "sectionID#124"],
            ["taskID#60", "sectionID#125"],
            ["taskID#61", "sectionID#125"],
            ["taskID#62", "sectionID#125"],
            ["taskID#63", "sectionID#125"],
            ["taskID#64", "sectionID#125"],
            ["taskID#65", "sectionID#125"],
            ["taskID#66", "sectionID#125"],
            ["taskID#67", "sectionID#125"],
            ["taskID#68", "sectionID#125"],
            ["taskID#69", "sectionID#125"],
            ["taskID#70", "sectionID#125"],
            ["taskID#71", "sectionID#125"],
            ["taskID#72", "sectionID#125"],
            ["taskID#73", "sectionID#125"],
            ["taskID#74", "sectionID#125"],
            ["taskID#75", "sectionID#125"],
            ["taskID#76", "sectionID#125"],
            ["taskID#77", "sectionID#125"],
            ["taskID#78", "sectionID#125"],
            ["taskID#79", "sectionID#125"],
            ["taskID#80", "sectionID#125"],
            ["taskID#81", "sectionID#125"],
            ["taskID#82", "sectionID#125"],
            ["taskID#83", "sectionID#125"],
            ["taskID#84", "sectionID#125"],
            ["taskID#85", "sectionID#125"],
            ["taskID#86", "sectionID#125"],
            ["taskID#87", "sectionID#125"],
            ["taskID#88", "sectionID#125"],
            ["taskID#89", "sectionID#125"],
            ["taskID#90", "sectionID#125"],
            ["taskID#91", "sectionID#125"],
            ["taskID#92", "sectionID#125"],
            ["taskID#93", "sectionID#125"],
            ["taskID#94", "sectionID#125"],
            ["taskID#95", "sectionID#125"],
            ["taskID#96", "sectionID#125"],
            ["taskID#97", "sectionID#125"],
            ["taskID#98", "sectionID#125"],
            ["taskID#99", "sectionID#125"],
            ["taskID#00", "sectionID#125"],
        ]);

        // prettier-ignore
        this.objectContents = new Map([
            ["sectionID#123", "August 21, 2022"],
            ["sectionID#124", "August 22, 2022"],
            ["sectionID#125", "August 23, 2022"],
            // ["taskID#1", "Lorem ipsum dolor sit amet consectetur adipisicing elit. "],
            // ["taskID#2", "Accusantium voluptatem excepturi suscipit quibusdam, pariatur deleniti ex provident, quaerat fuga earum quasi architecto aliquam natus dolores consequatur repellendus, quis exercitationem quod?"],
            // ["taskID#3", "Assumenda sit repudiandae voluptatum ipsum nulla facilis eligendi aspernatur commodi asperiores, aperiam hic corporis aliquam sint. "],
            // ["taskID#4", "Dolore autem, architecto neque recusandae, voluptatum esse accusantium repellendus corrupti adipisci molestiae culpa tenetur?"],
            // ["taskID#5", "Temporibus autem quia nam dolorum, officiis debitis rem, ipsam quisquam at esse maiores, itaque pariatur nisi voluptate illum rerum laboriosam doloribus corporis. "],
            // ["taskID#6", "Numquam porro soluta quaerat doloremque aspernatur voluptas minus!"],
            // ["taskID#7", "Tenetur iure at voluptates quaerat, illum quae omnis quidem numquam consectetur maxime porro placeat eligendi ut, doloremque, recusandae magni. "],
            // ["taskID#8", "Quo explicabo assumenda pariatur esse, ratione consequuntur perspiciatis ipsam similique blanditiis."],
            // ["taskID#9", "Error distinctio fuga veritatis nisi! Iure quam harum quas ipsum voluptas deserunt. "],
            // ["taskID#10", "Necessitatibus ad vero, voluptate reprehenderit ex odio quod architecto quibusdam, culpa officia mollitia tempora accusamus, consequuntur porro repudiandae?"],
            // ["taskID#11", "Nisi, sunt vel. "],
            // ["taskID#12", "Tempora numquam dolore earum tenetur animi cumque incidunt placeat, velit commodi, totam rerum! Nobis, consectetur eligendi assumenda nihil corporis praesentium maxime id, quidem amet aperiam nostrum? Voluptate."],
            // ["taskID#13", "Nulla quos consectetur aspernatur odio magnam repellendus dolores quae possimus perferendis voluptates inventore, est exercitationem nihil blanditiis error. "],
            // ["taskID#14", "Aspernatur at amet eaque accusantium atque cum molestias recusandae repudiandae velit necessitatibus?"],
            // ["taskID#15", "Doloribus sunt, debitis necessitatibus ratione commodi, labore at, odit cum consectetur accusamus eligendi beatae sit natus. "],
            // ["taskID#16", "Dolores delectus a veniam quam at cupiditate commodi magni, velit, voluptas dolorem reprehenderit accusamus."],
            // ["taskID#17", "Ratione, nulla quibusdam. "],
            // ["taskID#18", "Quidem nihil et repellat! Voluptatem vero natus aliquam nihil, quae quaerat accusamus quidem suscipit quasi debitis, perferendis voluptatum totam ratione nulla non ipsum. "],
            // ["taskID#19", "Modi aliquid asperiores necessitatibus."],
            // ["taskID#20", "Nisi incidunt magnam possimus quam. "],
            // ["taskID#21", "Neque unde minima, accusamus minus asperiores iusto soluta harum ullam rem assumenda suscipit, alias ipsam, sunt atque amet dolorum quo. "],
            // ["taskID#22", "Laudantium in repudiandae nostrum sunt."],
            ["taskID#1", "taskID#1"],
            ["taskID#2", "taskID#2"],
            ["taskID#3", "taskID#3"],
            ["taskID#4", "taskID#4"],
            ["taskID#5", "taskID#5"],
            ["taskID#6", "taskID#6"],
            ["taskID#7", "taskID#7"],
            ["taskID#8", "taskID#8"],
            ["taskID#9", "taskID#9"],
            ["taskID#10", "taskID#10"],
            ["taskID#11", "taskID#11"],
            ["taskID#12", "taskID#12"],
            ["taskID#13", "taskID#13"],
            ["taskID#14", "taskID#14"],
            ["taskID#15", "taskID#15"],
            ["taskID#16", "taskID#16"],
            ["taskID#17", "taskID#17"],
            ["taskID#18", "taskID#18"],
            ["taskID#19", "taskID#19"],
            ["taskID#20", "taskID#20"],
            ["taskID#21", "taskID#21"],
            ["taskID#22", "taskID#22"],
            ["taskID#23", "text content for taskID#23"],
            ["taskID#24", "text content for taskID#24"],
            ["taskID#25", "text content for taskID#25"],
            ["taskID#26", "text content for taskID#26"],

            ["taskID#27", "text content for taskID#27"],
            ["taskID#28", "text content for taskID#28"],
            ["taskID#29", "text content for taskID#29"],
            ["taskID#30", "text content for taskID#30"],
            ["taskID#31", "text content for taskID#31"],
            ["taskID#32", "text content for taskID#32"],
            ["taskID#33", "text content for taskID#33"],
            ["taskID#34", "text content for taskID#34"],
            ["taskID#35", "text content for taskID#35"],
            ["taskID#36", "text content for taskID#36"],
            ["taskID#37", "text content for taskID#37"],
            ["taskID#38", "text content for taskID#38"],
            ["taskID#39", "text content for taskID#39"],
            ["taskID#40", "text content for taskID#40"],
            ["taskID#41", "text content for taskID#41"],
            ["taskID#42", "text content for taskID#42"],
            ["taskID#43", "text content for taskID#43"],
            ["taskID#44", "text content for taskID#44"],
            ["taskID#45", "text content for taskID#45"],
            ["taskID#46", "text content for taskID#46"],
            ["taskID#47", "text content for taskID#47"],
            ["taskID#48", "text content for taskID#48"],
            ["taskID#49", "text content for taskID#49"],
            ["taskID#50", "text content for taskID#50"],
            ["taskID#51", "text content for taskID#51"],
            ["taskID#52", "text content for taskID#52"],
            ["taskID#53", "text content for taskID#53"],
            ["taskID#54", "text content for taskID#54"],
            ["taskID#55", "text content for taskID#55"],
            ["taskID#56", "text content for taskID#56"],
            ["taskID#57", "text content for taskID#57"],
            ["taskID#58", "text content for taskID#58"],
            ["taskID#59", "text content for taskID#59"],

            ["taskID#60", "text content for taskID#60"],
            ["taskID#61", "text content for taskID#61"],
            ["taskID#62", "text content for taskID#62"],
            ["taskID#63", "text content for taskID#63"],
            ["taskID#64", "text content for taskID#64"],
            ["taskID#65", "text content for taskID#65"],
            ["taskID#66", "text content for taskID#66"],
            ["taskID#67", "text content for taskID#67"],
            ["taskID#68", "text content for taskID#68"],
            ["taskID#69", "text content for taskID#69"],
            ["taskID#70", "text content for taskID#70"],
            ["taskID#71", "text content for taskID#71"],
            ["taskID#72", "text content for taskID#72"],
            ["taskID#73", "text content for taskID#73"],
            ["taskID#74", "text content for taskID#74"],
            ["taskID#75", "text content for taskID#75"],
            ["taskID#76", "text content for taskID#76"],
            ["taskID#77", "text content for taskID#77"],
            ["taskID#78", "text content for taskID#78"],
            ["taskID#79", "text content for taskID#79"],
            ["taskID#80", "text content for taskID#80"],
            ["taskID#81", "text content for taskID#81"],
            ["taskID#82", "text content for taskID#82"],
            ["taskID#83", "text content for taskID#83"],
            ["taskID#84", "text content for taskID#84"],
            ["taskID#85", "text content for taskID#85"],
            ["taskID#86", "text content for taskID#86"],
            ["taskID#87", "text content for taskID#87"],
            ["taskID#88", "text content for taskID#88"],
            ["taskID#89", "text content for taskID#89"],
            ["taskID#90", "text content for taskID#90"],
            ["taskID#91", "text content for taskID#91"],
            ["taskID#92", "text content for taskID#92"],
            ["taskID#93", "text content for taskID#93"],
            ["taskID#94", "text content for taskID#94"],
            ["taskID#95", "text content for taskID#95"],
            ["taskID#96", "text content for taskID#96"],
            ["taskID#97", "text content for taskID#97"],
            ["taskID#98", "text content for taskID#98"],
            ["taskID#99", "text content for taskID#99"],
            ["taskID#00", "text content for taskID#00"],
        ]);

        for (const [rowID, sectionID] of this.rowSections.entries()) {
            this.medium.addRowToSection(sectionID, rowID);
        }

        for (let i = 101; i < 1000; i++) {
            this.medium.addRowToSection("sectionID#126", `taskID#${i.toString()}`);
        }

        this.notifyDelegateFor("placementUpdate");

        setTimeout(() => {
            this.objectContents.set(
                "taskID#1",
                "Lorem ipsum dolor sit amet consectetur adipisicing elit. Omnis voluptatum labore in hic possimus dolor. Aliquam tempore unde quia natus hic optio modi excepturi. Reprehenderit natus recusandae dolores rerum omnis?"
            );
            this.notifyDelegateFor("objectUpdate", new Set([pSymbol.get("taskID#1")]));
        }, 1000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 3);
            this.notifyDelegateFor("placementUpdate");
        }, 2000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 2);
            this.medium.moveRow("taskID#1", 2);
            this.medium.moveRow("taskID#7", 2);
            this.medium.moveRow("taskID#7", 2);
            this.medium.moveRow("taskID#8", 2);
            this.medium.moveRow("taskID#1", 2);
            this.medium.moveRow("taskID#1", 2);
            this.medium.moveRow("taskID#7", 2);
            this.medium.moveRow("taskID#2", 2);
            this.medium.moveRow("taskID#7", 2);
            this.medium.moveRow("taskID#4", 2);
            this.notifyDelegateFor("placementUpdate");
        }, 3000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 20);
            this.notifyDelegateFor("placementUpdate");
        }, 4000);

        setTimeout(() => {
            this.medium.moveRowToAnotherSection("taskID#3", "sectionID#124", 0);
            this.notifyDelegateFor("placementUpdate");
        }, 5000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 1);
            this.notifyDelegateFor("placementUpdate");
        }, 6000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 2);
            this.notifyDelegateFor("placementUpdate");
        }, 7000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 3);
            this.notifyDelegateFor("placementUpdate");
        }, 8000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 4);
            this.notifyDelegateFor("placementUpdate");
        }, 9000);

        setTimeout(() => {
            this.medium.moveRow("taskID#3", 10);
            this.notifyDelegateFor("placementUpdate");
        }, 10000);

        setTimeout(() => {
            this.objectContents.set(
                "taskID#3",
                "Lorem ipsum dolor sit amet consectetur adipisicing elit. Omnis voluptatum labore in hic possimus dolor. Aliquam tempore unde quia natus hic optio modi excepturi. Reprehenderit natus recusandae dolores rerum omnis?"
            );
            this.notifyDelegateFor("objectUpdate", new Set([pSymbol.get("taskID#3")]));
        }, 11000);
    }

    loadTestDataset2() {
        setTimeout(() => {
            this.cache.placements.totalNumberOfItems = 1000;

            this.cache.placements.chronological = {
                offset: 0,
                items: [
                    "task#123",
                    "task#124",
                    "task#125",
                    "task#126",
                    "task#127",
                    "task#133",
                    "task#134",
                    "task#135",
                    "task#136",
                    "task#137",
                    "task#143",
                    "task#144",
                    "task#145",
                    "task#146",
                    "task#147",
                ],
            };

            this.cache.placements.hierarchical = {
                offset: 0,
                items: [
                    "task#123",
                    "task#143",
                    "task#133",
                    "task#124",
                    "task#144",
                    "task#134",
                    "task#125",
                    "task#145",
                    "task#126",
                    "task#146",
                    "task#135",
                    "task#127",
                    "task#147",
                    "task#136",
                    "task#137",
                ],
            };

            this.cache.tasks.set("task#123", { content: "task#123", parentId: "-1" });
            this.cache.tasks.set("task#124", { content: "task#124", parentId: "-1" });
            this.cache.tasks.set("task#125", { content: "task#125", parentId: "-1" });
            this.cache.tasks.set("task#126", { content: "task#126", parentId: "-1" });
            this.cache.tasks.set("task#127", { content: "task#127", parentId: "-1" });
            this.cache.tasks.set("task#133", { content: "task#133", parentId: "-1" });
            this.cache.tasks.set("task#134", { content: "task#134", parentId: "-1" });
            this.cache.tasks.set("task#135", { content: "task#135", parentId: "-1" });
            this.cache.tasks.set("task#136", { content: "task#136", parentId: "-1" });
            this.cache.tasks.set("task#137", { content: "task#137", parentId: "-1" });
            this.cache.tasks.set("task#143", { content: "task#143", parentId: "-1" });
            this.cache.tasks.set("task#144", { content: "task#144", parentId: "-1" });
            this.cache.tasks.set("task#145", { content: "task#145", parentId: "-1" });
            this.cache.tasks.set("task#146", { content: "task#146", parentId: "-1" });
            this.cache.tasks.set("task#147", { content: "task#147", parentId: "-1" });

            this.notifyDelegateFor("placementUpdate");
            console.log("test database 2 is loaded");
        }, 2);
    }

    getTextContent(objectSymbol) {
        const objectID = pSymbol.reverse(objectSymbol);
        // const match = objectID.match(/section/)
        // if (match.length > 0)
        if (this.objectContents.has(objectID)) return this.objectContents.get(objectID);
        return `loop generated task content for ${objectID}`;
    }
}
