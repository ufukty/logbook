# Logbook Javascript Library

## CellScroller a.k.a `AbstractCellScrollerViewController`

**Summary**

- This VC, keeps the memory consumption constant regardless of how many items should be displayed through user scrolls on page. So, a page with million cells will be as fast as a page with thousand cell in terms of page loading and navigation speed.
- Supports custom cell view controllers.
- Supports custom margins for before/between/after multiple custom cell view controllers. 
- CellScroller supports dynamic cell resizing via ResizeObserver, and it works automatically. You only have to provide default and average heights for scroll height assumption for not-loaded-yet cells.

**Usage**

- Initialize instances of `AbstractCellScrollerViewController` through inheriting it with a concrete class.
  ```js
  class MyScroller extends AbstractCellScrollerViewController {}; // keep reading to learn how to implement
  
  const myScrollerInstance = new MyScroller();
  document.appendChild(myScrollerInstance.container);
  ```

- Add the `css` file to your web page.
  ```html
  <link rel="stylesheet" href="./path/to/cell-scroller-view.css" type="text/css" />
  ```

**Configuration**

```javascript
class MyScroller extends AbstractCellScrollerViewController {
  constructor() {
    super()
    this.config.margins = {}; // see below for valid values
  }
  
  getDefaultCellHeightForObjectId(objectId) { 
    if (objectId === "specificId") return 60; // in pixels
    else return 25; 
  } 
  
  getAverageCellHeight() { return 30; }
};
```

**Placement**

- You can push initial placement by assigning an array of object identifiers (in type of `Symbol`) to `this.config.placement.items`
  ```javascript
  class MyScroller extends AbstractCellScrollerCellViewController {};
  
  const objectSymbols = [
    Symbol("cell-001"), Symbol("cell-002"), Symbol("cell-003"),
  ];
  myScrollerInstance.config.placement = {
    totalNumberOfItems: 6,
    offset: 0,
    items: objectsSymbols,
  };
  myScrollerInstance.updateView();
  ```

- You can update the placement of cells by modifying `items` array and calling `.updateView()` method after. After CellScroller recalculates cell positions, all items will move to their new position with transition.
  ```javascript
  objectSymbols.push(objectSymbols.shift()) // moving first item to the end
  myScrollerInstance.updateView()
  ```

**Custom cells**

```javascript
// Create custom cell from abstract cell view controller class
class CustomCellVC extends AbstractCellScrollerCellViewController {
  prepareForFree() {}
	prepareForUse() {}
	setContent() {}
}

// Register custom cell
class MyScroller extends CellScrollerViewController {
  constructor() {
    // Keep that symbol for using later
    const myCustomCellVC = Symbol("myCustomCellVC"); 
		this.registerCellIdentifier(myCustomCellVC, () => {
			// this is how CellScroller will know how 
      // to initialize your custom cells
      return new CustomCellVC(); 
    });
  }
};
```

**Staging features**

- CellScroller can notify you when a cell just came in to viewport and left from it. 
- To receive notifications implement those two methods:
  ```js
  class MyScroller extends CellScrollerViewController {
    cellAppears(objectSymbol, cellPositioner) {}
    cellDisappears(objectSymbol, cellPositioner) {}
  };
  ```

**Section headers**

* Section headers are just custom cells that are different than main content. 

**Cell folding**

- CellScroller supports animating the folding process to provide better user-experience. 
  
- To fold some child nodes over a parent node, give this command: 
  
  ```js
  myScrollerInstance.fold("cell-001", ["cell-002", "cell-003"])
  //                      ^parent      ^children
  ```
  
- To unfold previously folded nodes:
  ```js
  myScrollerInstance.unfold("cell-001")
  ```

**Full code for essential features**

```javascript
```

