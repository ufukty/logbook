import { Solve } from "./solve";

onmessage = async function (ev: MessageEvent) {
    const { n, que, hash_ } = ev.data;

    try {
        const result = await Solve(n, que, hash_);
        postMessage({ success: true, result });
    } catch (err) {
        if (err instanceof Error) {
            postMessage({ success: false, error: err.message });
        }
    }
};
