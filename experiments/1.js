import { AbstractManagedLayoutViewController } from "../js/baja.sl/AbstractManagedLayoutViewController.js";
import { Flow, VERTICAL } from "../js/baja.sl/Layout/Calculators/Flow.js";
import { AbstractManagedLayoutCellViewController } from "../js/baja.sl/AbstractManagedLayoutCellViewController.js";
import { itemCellPairing } from "../js/baja.sl/ItemCellPairing.js";

import { Layout } from "../js/baja.sl/Layout/Layout.js";
import { adoption, iota, symbolizer } from "../js/baja.sl/utilities.js";
import { itemMeasurer } from "../js/baja.sl/ItemMeasurer.js";
import { Size, Spacing } from "../js/baja.sl/Layout/Coordinates.js";

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

    prepareForUse() {
        super.prepareForUse();
        this.dom.container.innerText = "0";
    }

    firstLevelOfPresentation() {
        this.dom.container.innerText = "1";
    }

    secondLevelOfPresentation() {
        this.dom.container.innerText = "2";
    }

    thirdLevelOfPresentation() {
        this.dom.container.innerText = "3";
    }
}

const VIEW_CONTROLLER_SYMBOL_TASK = symbolizer.symbolize(iota());

class ConcreteLayoutPresenterViewController extends AbstractManagedLayoutViewController {
    constructor() {
        super();

        itemCellPairing.registerViewControllerConstructor(VIEW_CONTROLLER_SYMBOL_TASK, () => {
            return new BasicViewController();
        });

        this._setupContainer();
        this._setupMainLayout();

        const mainLayout = this.config.layoutEnvironment;
        const mainEnvironmentSymbol = mainLayout.environmentSymbol;

        const itemSymbols = [Symbol("1"), Symbol("2"), Symbol("3"), Symbol("4")];

        this.layoutPipes.flow.config.spacing.set(VIEW_CONTROLLER_SYMBOL_TASK, new Spacing(200, 100, 200));

        itemMeasurer.setAverageSize(mainEnvironmentSymbol, new Size(100, 100));
        
        itemSymbols.forEach((itemSymbol) => {
            itemMeasurer.setDefaultSize(itemSymbol, mainEnvironmentSymbol, new Size(100, 100));
        });

        this.layoutPipes.flow.newPlacement(itemSymbols);
    }

    /** @private */
    _setupContainer() {
        this.dom.container.style.width = "100%";
        this.dom.container.style.height = "100%";
        this.dom.container.style.overflowX = "hidden";
        this.dom.container.style.overflowY = "scroll";
    }

    /** @private */
    _setupMainLayout() {
        this.layoutPipes = {
            flow: new Flow(),
            // indentation: new Indentation(),
            // align: new Align(),
            // focusStabilizer: new FocusStabilizer(),
            // counterShift: new CounterShift(),
            // avatars: new AvatarLayout(),
            // panes: new Panes(),
            // padding: new Padding(20, 20, 20, 20),
            // measure: new MeasureContainer(),
        };

        this.layoutPipes.flow.config.direction = VERTICAL;
        // align.config.alignOn = HORIZONTAL_CENTER;

        this.config.layoutEnvironment = new Layout();

        // prettier-ignore
        this.config.layoutEnvironment
            .connectCalculator(this.layoutPipes.flow)
        // .connectMutator(measure);
        // .connectMutator(align)
        // .connectMutator(indentation)
        // .connectMutator(counterShift)
        // .connectMutator(focusStabilizer)
        // .connectDecorator(foldedItems)
        // .connectDecorator(avatars)
        // .connectDecorator(panes)
        // .connectMutator(padding);

        // autoFocus(); // TODO:

        var setContainerSize = () => {
            const computedStyle = getComputedStyle(this.dom.container);
            const containerSize = this.config.layoutEnvironment.passedThroughPipeline.containerSize;
            containerSize.width = parseFloat(computedStyle.getPropertyValue("width"));
            containerSize.height = parseFloat(computedStyle.getPropertyValue("height"));
        };
        setContainerSize();

        this.config.layoutEnvironment.start();
    }

    _playTest() {
        // setTimeout(() => {
        //     const selectedItemSymbol = "task#1";
        //     focusStabilizer.stabilize(selectedItemSymbol);
        //     infiniteSheetPositions.refreshPipeline();
        // }, 1000);
        // setTimeout(() => {
        //     avatars.config.anchors.set();
        // }, 2000);
    }
}

function main() {
    const managedLayoutView = new ConcreteLayoutPresenterViewController();
    adoption(document.body, [managedLayoutView.dom.container]);
}

main();
