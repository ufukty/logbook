/*
    Near up-to-date mirror of data on the server. 
    Only-source-of-truth for any data consumer component in the frontend.
*/
class LocalSourceOfTruth {
    constructor() {
        this.delegates = {
            linearizedHierarchicalOrdering: () => {}, // declared to be assigned later by user
        };

        this.linearizedHierarchicalOrdering = {
            "1000": [
                "5627716c-fb7c-50f7-b8bc-c459311b471e",
                "c6130fc3-5589-53e7-adde-f6f9d1d80524",
                "333d0c95-5c29-57cd-b091-21d1e238901e",
                "792ad984-e09f-556f-8316-5c85f6dc6a05",
                "7ad104bf-ffc8-558f-8173-a8971cb3c2fc",
                "86db66fd-b43f-5029-96da-706f168a9fa0",
                "4d6a4812-5f9b-57b0-b04c-42f1477493f7",
                "4c092684-5dbe-55a2-a945-e88deeff3fa1",
            ],
            "1100": [
                "d530e9e5-c1b9-5804-a9ad-891f820ea9c0",
                "d987e8fc-a26c-5884-9a3e-1b1bb958553a",
                "fa8614c7-6165-5ece-a620-d4a1e4ca612b",
                "ae7a2be2-4851-5542-94d9-9e141ccf6678",
                "ede7eead-c094-5db4-9e5c-21663412405f",
                "422754a3-c815-55a9-8057-be84161fd7fb",
                "d02eb51a-64a5-57b5-b257-15ea614a3701",
                "26c9a035-f222-5350-8c2e-23019f2b74f6",
                "ad42bc65-4aae-5631-8dfd-3ac97cf3ed13",
            ],
        };
    }

    _invalidateCacheforLinearizedHierarchicalOrdering(taskOrder) {
        const cacheGroupKey = (Math.floor(taskOrder / 100) * 100).toString();
        if (this.linearizedHierarchicalOrdering.hasOwnProperty(cacheGroupKey)) {
            delete this.linearizedHierarchicalOrdering[cacheGroupKey];
        }
    }

    _notifyAboutUpdateInLinearizedHierarchicalOrdering() {
        this.delegates.linearizedHierarchicalOrdering();
    }

    getTaskDetails(taskId, callback) {
        // [if] data is not available
        //         [OR] if data is not up-to-date (not invalidated)
        //               fetch data from server
        //               append to the memory block
        // [return] data from memory block
    }

    updateTaskIndex(taskId, newIndex) {
        // invalidate task on memory
        // send modifications to the server
        //             get list of invalidated tasks from server
        //             mark the in-memory representations as invalidated too
        // return list of updated tasks details
    }

    updateTaskParent(taskId, newParentId, index) {
        // same with updateTaskIndex
    }
}

export const localSourceOfTruth = new LocalSourceOfTruth();
