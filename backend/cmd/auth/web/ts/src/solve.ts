const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";

const encoder = new TextEncoder();

async function hash(input: string): Promise<string> {
    try {
        const buf = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(input));
        const binary = String.fromCharCode(...new Uint8Array(buf));
        return window.btoa(binary);
    } catch (error) {
        console.error("Hashing failed:", error);
        throw error;
    }
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

export const ErrNotFound = new Error("not found");
export const ErrInvalidChallange = new Error("invalid challenge: que or hash is empty");

export async function Solve(n: number, que: string, hash_: string): Promise<string> {
    if (hash_.length === 0 || n === 0) {
        throw ErrInvalidChallange;
    }
    let p = start(n);
    do {
        let cand = p.value + que;
        if ((await hash(cand)) === hash_) {
            return cand;
        }
    } while (p.iterate());
    throw ErrNotFound;
}
