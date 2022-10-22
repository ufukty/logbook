import {
    AbstractManagedLayoutViewController,
    TRIGGER_REPLACEMENT,
    TRIGGER_SCROLL_LISTENER,
} from "../js/baja.sl/AbstractManagedLayoutViewController.js";
import { Flow, VERTICAL } from "../js/baja.sl/Layout/Calculators/Flow.js";
import { AbstractManagedLayoutCellViewController } from "../js/baja.sl/AbstractManagedLayoutCellViewController.js";
import { itemCellPairing } from "../js/baja.sl/ItemCellPairing.js";

import { Layout } from "../js/baja.sl/Layout/Layout.js";
import { adoption, createElement, iota, symbolizer } from "../js/baja.sl/utilities.js";
import { itemMeasurer } from "../js/baja.sl/ItemMeasurer.js";
import { Area, Size, Spacing } from "../js/baja.sl/Layout/Coordinates.js";
import { resizeObserverWrapper } from "../js/baja.sl/ResizeObserverWrapper.js";

import { Align, HORIZONTAL_CENTER, HORIZONTAL_LEFT, HORIZONTAL_RIGHT } from "../js/baja.sl/Layout/Mutators/Align.js";

class BasicViewController extends AbstractManagedLayoutCellViewController {
    constructor() {
        super();

        this.dom.container.style.width = "100px";
        this.dom.container.style.height = "100px";
    }

    async prepareForFreeAsync() {
        await super.prepareForFreeAsync();
        this.dom.container.innerText = "";
    }

    firstLevelOfPresentation() {
        this.dom.container.innerText += "1";
        this.dom.container.style.backgroundColor = "lightgray";
    }

    secondLevelOfPresentation() {
        this.dom.container.innerText += "2";
        this.dom.container.style.backgroundColor = "gray";
    }

    thirdLevelOfPresentation() {
        this.dom.container.innerText += "3";
        this.dom.container.style.backgroundColor = "lightblue";
    }
}

const VIEW_CONTROLLER_SYMBOL_TASK = symbolizer.symbolize(iota());

class CustomManagedLayoutViewController extends AbstractManagedLayoutViewController {
    constructor() {
        super();

        this.config = {
            ...this.config,
            zoneOffsets: {
                preload: 0.6,
                parking: 0.6,
            },
        };

        // this._setupContainer();
    }

    /** @private */
    _setupContainer() {
        this.dom.container.style.width = "100%";
        this.dom.container.style.height = "100%";
    }

    playTestCase() {
        const layoutPipes = {
            flow: new Flow(VERTICAL),
            // indentation: new Indentation(),
            align: new Align(HORIZONTAL_LEFT),
            // focusStabilizer: new FocusStabilizer(),
            // counterShift: new CounterShift(),
            // avatars: new AvatarLayout(),
            // panes: new Panes(),
            // padding: new Padding(20, 20, 20, 20),
            // measure: new MeasureContainer(),
        };

        this.config.layout = new Layout()
            .connectCalculator(layoutPipes.flow)
            // .connectMutator(measure)
            .connectMutator(layoutPipes.align);
        // .connectMutator(indentation)
        // .connectMutator(counterShift)
        // .connectMutator(focusStabilizer)
        // .connectDecorator(foldedItems)
        // .connectDecorator(avatars)
        // .connectDecorator(panes)
        // .connectMutator(padding);

        // autoFocus(); // TODO:
        const mainEnvironmentSymbol = this.config.layout.environmentSymbol;

        const itemSymbols = [
            symbolizer.symbolize("1"),
            symbolizer.symbolize("2"),
            symbolizer.symbolize("3"),
            symbolizer.symbolize("4"),
            symbolizer.symbolize("5"),
            symbolizer.symbolize("6"),
            symbolizer.symbolize("7"),
            symbolizer.symbolize("8"),
            symbolizer.symbolize("9"),
            symbolizer.symbolize("10"),
            symbolizer.symbolize("11"),
            symbolizer.symbolize("12"),
            symbolizer.symbolize("13"),
            symbolizer.symbolize("14"),
            symbolizer.symbolize("15"),
            symbolizer.symbolize("16"),
            symbolizer.symbolize("17"),
            symbolizer.symbolize("18"),
            symbolizer.symbolize("19"),
        ];

        itemMeasurer.setAverageSize(mainEnvironmentSymbol, new Size(100, 100));
        itemSymbols.forEach((itemSymbol) => {
            itemMeasurer.setDefaultSize(itemSymbol, mainEnvironmentSymbol, new Size(100, 100));
            // itemMeasurer.setSize(itemSymbol, mainEnvironmentSymbol, new Size(Math.floor(Math.random() * 100), 100));
            itemCellPairing.setCellKindForItem(itemSymbol, mainEnvironmentSymbol, VIEW_CONTROLLER_SYMBOL_TASK);
        });

        this.registerCellViewControllerConstructor(VIEW_CONTROLLER_SYMBOL_TASK, () => {
            return new BasicViewController();
        });

        this.config.layout.subscribe(() => {
            this.updateView(TRIGGER_REPLACEMENT);
        });

        layoutPipes.flow.config.spacing.set(VIEW_CONTROLLER_SYMBOL_TASK, new Spacing(200, 1, 200));
        layoutPipes.flow.setPlacement(itemSymbols);
        this.config.layout.scheduleRecalculation();

        // setTimeout(() => {
        //     layoutPipes.align.updateWithNewAlignment(HORIZONTAL_RIGHT);

        //     setTimeout(() => {
        //         layoutPipes.align.updateWithNewAlignment(HORIZONTAL_CENTER);

        //         setTimeout(() => {
        //             itemSymbols.push(itemSymbols.shift());
        //             // itemSymbols.push(itemSymbols.shift());
        //             // itemSymbols.push(itemSymbols.shift());
        //             layoutPipes.flow.updateWithNewPlacement(itemSymbols);
        //         }, 1000);
        //     }, 1000);
        // }, 1000);

        // setTimeout(() => {
        //     const selectedItemSymbol = "task#1";
        //     focusStabilizer.stabilize(selectedItemSymbol);
        //     infiniteSheetPositions.refreshPipeline();
        // }, 1000);
        // setTimeout(() => {
        //     avatars.config.anchors.set();
        // }, 2000);
    }

    /**
     * @param {AbstractManagedLayoutCellViewController} managedLayoutCellViewController
     * @param {ItemSymbol} itemSymbol
     * Populate content of managedLayoutCellViewController according to itemSymbol
     */
    populateCellForItem(managedLayoutCellViewController, itemSymbol) {
        if (managedLayoutCellViewController instanceof AbstractManagedLayoutCellViewController) {
            managedLayoutCellViewController.dom.container.innerText = `${symbolizer.desymbolize(itemSymbol)} 1`;
        }
    }
}

function main() {
    const managedLayoutViewController = new CustomManagedLayoutViewController();
    adoption(document.body, adoption(managedLayoutViewController.dom.container));

    const passNewViewport = () => {
        managedLayoutViewController.setViewport(
            new Area(
                window.scrollX,
                window.scrollY,
                window.scrollX + window.innerWidth,
                window.scrollY + window.innerHeight
            )
        );
    };
    const managedLayoutUpdater = () => {
        passNewViewport();
        managedLayoutViewController.updateView(TRIGGER_SCROLL_LISTENER);
    };
    window.addEventListener("scroll", managedLayoutUpdater);
    window.addEventListener("resize", managedLayoutUpdater);
    passNewViewport();

    managedLayoutViewController.playTestCase();
}

window.addEventListener("load", main);
