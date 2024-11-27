const worker = new Worker("solver.js");

const n = 2;
const que =
    "4QYD5OASOT7NAPFKN4DXWRJSSPSIN7GUJUS3LL2XP33XJ5XJTUTTEX6GIDS6G4YUZ5YWUWSG2DIMPOZ2WXHYRUZ3GJZXIZWSI2NL3NRC4NT3YOQWX2Q4THRV5EKEXQDMVBQWJ4ZHSXP45BKPHBIJ7YPC674BNUFTKX7HAPEKSLSAFOIECVGQFE6GS5EAGPICWAVANKJBR4BBWTBYIAIUS2IAHKK7FBMUMVJV2J3FOW2JSRF";
const hash_ = "CZZVTBQASPQJLUGY4LTE7GKCVFBDD6NFYHHSFINMUNF55DPF6MYQ";

worker.onmessage = function (e) {
    if (e.data.success) {
        console.log("Found:", e.data.result);
    } else {
        console.error("Error:", e.data.error);
    }
};

// Send data to the worker to start the computation
worker.postMessage({ n, que, hash_ });
