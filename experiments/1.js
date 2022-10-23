import {
    AbstractManagedLayoutViewController,
    TRIGGER_REPLACEMENT,
    TRIGGER_SCROLL_LISTENER,
} from "../js/baja.sl/AbstractManagedLayoutViewController.js";
import { Flow, VERTICAL } from "../js/baja.sl/Layout/Calculators/Flow.js";
import { AbstractManagedLayoutCellViewController } from "../js/baja.sl/AbstractManagedLayoutCellViewController.js";
import { itemCellPairing } from "../js/baja.sl/ItemCellPairing.js";

import { Layout } from "../js/baja.sl/Layout/Layout.js";
import { adoption, createElement, iota, pick, symbolizer } from "../js/baja.sl/utilities.js";
import { itemMeasurer } from "../js/baja.sl/ItemMeasurer.js";
import { Area, Size, Spacing } from "../js/baja.sl/Layout/Coordinates.js";
import { resizeObserverWrapper } from "../js/baja.sl/ResizeObserverWrapper.js";

import { Align, HORIZONTAL_CENTER, HORIZONTAL_LEFT, HORIZONTAL_RIGHT } from "../js/baja.sl/Layout/Mutators/Align.js";

class BasicViewController extends AbstractManagedLayoutCellViewController {
    constructor() {
        super();

        // this.dom.container.style.width = "100px";
        // this.dom.container.style.height = "100px";
    }

    async prepareForFreeAsync() {
        await super.prepareForFreeAsync();
        this.dom.container.innerHTML = "";
    }

    firstLevelOfPresentation() {
        console.log(this.config.itemSymbol);
        this.dom.container.innerHTML = ` 1-${symbolizer.desymbolize(this.config.itemSymbol)}<br>`;
        this.dom.container.style.backgroundColor = "lightgray";
    }

    secondLevelOfPresentation() {
        this.dom.container.innerHTML += ` 2-${symbolizer.desymbolize(this.config.itemSymbol)}<br>`;
        this.dom.container.style.backgroundColor = "gray";
    }

    thirdLevelOfPresentation() {
        this.dom.container.innerHTML += ` 3-${symbolizer.desymbolize(this.config.itemSymbol)}<br>`;
        this.dom.container.style.backgroundColor = "lightblue";
        this.dom.container.innerHTML += pick([
            "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
            "Donec gravida consequat orci, sed luctus arcu lacinia eget.",
            "Cras interdum nibh nunc, in ornare risus feugiat sed.",
            "Maecenas porta, lectus quis consectetur hendrerit, orci massa eleifend arcu, mollis molestie mauris felis nec ligula.",
            "Nulla semper tempus sagittis.",
            "Donec semper vel dolor vel porta.",
            "Nam vel placerat tellus.",
            "Quisque venenatis non felis sed hendrerit.",
        ]);
    }
}

const VIEW_CONTROLLER_SYMBOL_TASK = symbolizer.symbolize(iota());

class CustomManagedLayoutViewController extends AbstractManagedLayoutViewController {
    constructor() {
        super();

        this.config = {
            ...this.config,
            zoneOffsets: {
                preload: 1.2,
                parking: 1.3,
            },
            updateMaxFrequency: 20,
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
            symbolizer.symbolize("20"),
            symbolizer.symbolize("21"),
            symbolizer.symbolize("22"),
            symbolizer.symbolize("23"),
            symbolizer.symbolize("24"),
            symbolizer.symbolize("25"),
            symbolizer.symbolize("26"),
            symbolizer.symbolize("27"),
            symbolizer.symbolize("28"),
            symbolizer.symbolize("29"),
            symbolizer.symbolize("30"),
            symbolizer.symbolize("31"),
            symbolizer.symbolize("32"),
            symbolizer.symbolize("33"),
            symbolizer.symbolize("34"),
            symbolizer.symbolize("35"),
            symbolizer.symbolize("36"),
            symbolizer.symbolize("37"),
            symbolizer.symbolize("38"),
            symbolizer.symbolize("39"),
            symbolizer.symbolize("40"),
            symbolizer.symbolize("41"),
            symbolizer.symbolize("42"),
            symbolizer.symbolize("43"),
            symbolizer.symbolize("44"),
            symbolizer.symbolize("45"),
            symbolizer.symbolize("46"),
            symbolizer.symbolize("47"),
            symbolizer.symbolize("48"),
            symbolizer.symbolize("49"),
            symbolizer.symbolize("50"),
            symbolizer.symbolize("51"),
            symbolizer.symbolize("52"),
            symbolizer.symbolize("53"),
            symbolizer.symbolize("54"),
            symbolizer.symbolize("55"),
            symbolizer.symbolize("56"),
            symbolizer.symbolize("57"),
            symbolizer.symbolize("58"),
            symbolizer.symbolize("59"),
            symbolizer.symbolize("60"),
            symbolizer.symbolize("61"),
            symbolizer.symbolize("62"),
            symbolizer.symbolize("63"),
            symbolizer.symbolize("64"),
            symbolizer.symbolize("65"),
            symbolizer.symbolize("66"),
            symbolizer.symbolize("67"),
            symbolizer.symbolize("68"),
            symbolizer.symbolize("69"),
            symbolizer.symbolize("70"),
            symbolizer.symbolize("71"),
            symbolizer.symbolize("72"),
            symbolizer.symbolize("73"),
            symbolizer.symbolize("74"),
            symbolizer.symbolize("75"),
            symbolizer.symbolize("76"),
            symbolizer.symbolize("77"),
            symbolizer.symbolize("78"),
            symbolizer.symbolize("79"),
            symbolizer.symbolize("80"),
            symbolizer.symbolize("81"),
            symbolizer.symbolize("82"),
            symbolizer.symbolize("83"),
            symbolizer.symbolize("84"),
            symbolizer.symbolize("85"),
            symbolizer.symbolize("86"),
            symbolizer.symbolize("87"),
            symbolizer.symbolize("88"),
            symbolizer.symbolize("89"),
            symbolizer.symbolize("90"),
            symbolizer.symbolize("91"),
            symbolizer.symbolize("92"),
            symbolizer.symbolize("93"),
            symbolizer.symbolize("94"),
            symbolizer.symbolize("95"),
            symbolizer.symbolize("96"),
            symbolizer.symbolize("97"),
            symbolizer.symbolize("98"),
            symbolizer.symbolize("99"),
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
    // populateCellForItem(managedLayoutCellViewController, itemSymbol) {
    //     if (managedLayoutCellViewController instanceof AbstractManagedLayoutCellViewController) {
    //         managedLayoutCellViewController.dom.container.innerText = `${symbolizer.desymbolize(itemSymbol)} 1`;
    //     }
    // }
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
