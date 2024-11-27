const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

const encoder = new TextEncoder();

async function hash(input: string): Promise<string> {
    const buf = await crypto.subtle.digest("SHA-256", encoder.encode(input));
    const arr = new Uint8Array(buf); // .toString();
    let binary = "";
    for (let i = 0; i < arr.length; i++) {
        binary += String.fromCharCode(arr[i]);
    }
    return window.btoa(binary);
}

class Prefix {
    value: string;

    constructor(value: string) {
        this.value = value;
    }

    iterate() {
        let pb = this.value.split("");
        for (let i = pb.length - 1; i >= 0; i--) {
            if (i === 0 && pb[i] === alphabet[alphabet.length - 1]) {
                return false;
            }
            let j = alphabet.indexOf(pb[i]);
            pb[i] = alphabet[(j + 1) % alphabet.length];
            this.value = pb.join("");
            if (pb[i] !== alphabet[0]) {
                return true;
            }
        }
        return true;
    }
}

function start(n: number) {
    return new Prefix(alphabet[0].repeat(n));
}

async function Solve(n: number, que: string, hash_: string): Promise<string> {
    if (hash_.length === 0 || n === 0) {
        throw new Error("invalid challenge: que or hash is empty");
    }
    let p = start(n);
    do {
        let cand = p.value + que;
        if ((await hash(cand)) === hash_) {
            return cand;
        }
    } while (p.iterate());
    throw new Error("not found");
}

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
