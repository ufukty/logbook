import { adoption, domElementReuseCollector, createElement, toggleAnimationWithClass } from "../utilities.js";
import { InfiniteSheetTableCellViewController } from "./InfiniteSheetTableCellViewController.js";
import { AbstractViewController } from "./AbstractViewController.js";
import { AbstractTableCellViewController } from "./AbstractTableCellViewController.js";
import { AbstractTableViewController } from "./AbstactTableViewController.js";
import InfiniteSheetHeader from "./InfiniteSheetHeader.js";
import InfiniteSheetTask from "./InfiniteSheetTask.js";

export class InfiniteSheet extends AbstractTableViewController {
    constructor() {
        super();

        const regularCellId = Symbol("regularCellViewContainer");

        // const tasks = ["taskID#234", "taskID#235", "taskID#236", "taskID#237", "taskID#238", "taskID#239"];
        // const symbols = {
        //     "sectionID#123": Symbol("sectionID#123"),
        //     "sectionID#124": Symbol("sectionID#124"),
        //     "sectionID#125": Symbol("sectionID#125"),
        //     "taskID#234": Symbol("taskID#234"),
        //     "taskID#235": Symbol("taskID#235"),
        //     "taskID#236": Symbol("taskID#236"),
        //     "taskID#237": Symbol("taskID#237"),
        //     "taskID#238": Symbol("taskID#238"),
        //     "taskID#239": Symbol("taskID#239"),
        // };

        this.config.margins = {
            pageContent: {
                before: 100,
                after: 300,
            },
            section: {
                before: 100,
                between: 100,
            },
            row: {
                before: 100,
                between: 20,
            },
        };

        // this.config.placement.sections = [symbols["sectionID#123"], symbols["sectionID#124"], symbols["sectionID#125"]];

        // for (const [i, sectionIDSymbol] of this.config.placement.sections.slice(0, 3).entries()) {
        //     this.config.objectReuseIdentifiers.set(sectionIDSymbol, regularCellId);

        //     this.config.placement.rows.set(sectionIDSymbol, []);
        //     const arr = this.config.placement.rows.get(sectionIDSymbol);
        //     const task1 = symbols[tasks[2 * i]];
        //     const task2 = symbols[tasks[2 * i + 1]];
        //     arr.push(task1);
        //     arr.push(task2);
        //     this.config.objectReuseIdentifiers.set(task1, regularCellId);
        //     this.config.objectReuseIdentifiers.set(task2, regularCellId);
        // }

        // this.findSection = {
        //     "taskID#234": "sectionID#123",
        //     "taskID#235": "sectionID#123",
        //     "taskID#236": "sectionID#124",
        //     "taskID#237": "sectionID#124",
        //     "taskID#238": "sectionID#125",
        //     "taskID#239": "sectionID#125",
        // };

        // this.objectComputedHeights = {
        //     "sectionID#123": 25,
        //     "sectionID#123": 25,
        //     "sectionID#124": 25,
        //     "sectionID#124": 25,
        //     "sectionID#125": 25,
        //     "sectionID#125": 25,
        //     "taskID#234": 25,
        //     "taskID#235": 25,
        //     "taskID#236": 25,
        //     "taskID#237": 25,
        //     "taskID#238": 25,
        //     "taskID#239": 25,
        // };

        // this.objectTextContent = {
        //     "sectionID#123": "section 123",
        //     "sectionID#123": "section 123",
        //     "sectionID#124": "section 124",
        //     "sectionID#124": "section 124",
        //     "sectionID#125": "section 125",
        //     "sectionID#125": "section 125",
        //     "taskID#234": "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
        //     "taskID#235": "Vivamus vitae nibh nec tortor porta congue quis eu ante",
        //     "taskID#236": "Aliquam rhoncus tortor nec elit molestie, a rutrum odio hendrerit",
        //     "taskID#237": "Donec varius feugiat purus id sagittis",
        //     "taskID#238": "Nunc quis fringilla tellus, sed aliquam dui",
        //     "taskID#239": "Proin molestie dolor eget purus molestie, a cursus mi iaculis",
        // };

        // this.config.this.objectReuseIdentifiers = {
        //     "sectionID#123": regularCellId,
        //     "sectionID#123": regularCellId,
        //     "sectionID#124": regularCellId,
        //     "sectionID#124": regularCellId,
        //     "sectionID#125": regularCellId,
        //     "sectionID#125": regularCellId,
        //     "taskID#234": regularCellId,
        //     "taskID#235": regularCellId,
        //     "taskID#236": regularCellId,
        //     "taskID#237": regularCellId,
        //     "taskID#238": regularCellId,
        //     "taskID#239": regularCellId,
        // };

        this.config.defaultHeightForReuseId.set(regularCellId, 25);

        domElementReuseCollector.registerItemIdentifier(regularCellId, () => {
            const cell = new InfiniteSheetTableCellViewController();
            adoption(this.anchorPosition, [cell.container]);
            return cell;
        });
    }

    deleteTask(taskId) {
        const ref = this.getReferenceOfAllocatedRowElement(0, 0);
        this.hideRowOnce(0, 1);
        this.state.effectiveOrdering[0].splice(1);
        this.calculateElementBounds();
        this.rePosition();
        debugger;
    }
}
