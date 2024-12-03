export function EqualizerChallange(n: number, que: string, hash_: string): Promise<string> {
    return new Promise<string>((resolve, reject) => {
        const worker: Worker = new Worker("solve.worker.js");

        worker.onmessage = function (e) {
            if (e.data.success) {
                resolve(e.data.result);
            } else {
                reject(e.data.error);
            }
        };

        worker.onerror = function (err) {
            reject(err);
        };

        // Send data to the worker to start the computation
        worker.postMessage({ n, que, hash_ });
    });
}
