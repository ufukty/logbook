export function findFirstGreaterOrClosestItem(orderedList, searchItem) {
    var lastItemIndex = orderedList.length - 1;
    if (searchItem <= orderedList[0]) {
        // if searchItem is smaller than the smallest item on orderedList
        return 0;
    } else if (orderedList[lastItemIndex] <= searchItem) {
        // if searchItem is bigger than the biggest item on orderedList
        return lastItemIndex;
    } else {
        // if searchItem is in between first and last item of
        // orderedList, perform below instructions based on
        // binary search.
        var lo = 0,
            hi = lastItemIndex,
            mid = undefined;
        while (hi - lo > 1) {
            mid = Math.floor((lo + hi) / 2);
            if (orderedList[mid] <= searchItem) {
                lo = mid;
            } else {
                hi = mid;
            }
        }
        return lo;
    }
}

export function averageInt(listOfValues) {
    var total = 0;
    listOfValues.forEach((value) => {
        total += value;
    });
    return Math.floor(total / listOfValues.length);
}
