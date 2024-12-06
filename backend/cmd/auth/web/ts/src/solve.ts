import { hash } from "./hash";
import { alphabet } from "./alphabet";

const ML = 3;

class Prefix {
    value: string;
    alphabetSize: number;

    constructor(alphabetSize: number) {
        this.value = alphabet[0].repeat(ML);
        this.alphabetSize = alphabetSize;
    }

    iterate() {
        let pb = this.value.split("");
        for (let i = pb.length - 1; i >= 0; i--) {
            if (i === 0 && pb[i] === alphabet[this.alphabetSize - 1]) {
                return false;
            }
            let j = alphabet.indexOf(pb[i]);
            pb[i] = alphabet[(j + 1) % this.alphabetSize];
            this.value = pb.join("");
            if (pb[i] !== alphabet[0]) {
                return true;
            }
        }
        return true;
    }
}

export const ErrEmptyHashed = new Error("'hashed' is empty");
export const ErrEmptyMasked = new Error("'masked' is empty");
export const ErrMaxDifficulty = new Error("difficulty needs to be smaller than the alphabet size");
export const ErrMinDifficulty = new Error("difficulty needs to be bigger than 1");
export const ErrNotFound = new Error("not found");

export async function Solve(difficulty: number, masked: string, hashed: string): Promise<string> {
    if (difficulty < 2) {
        throw ErrMinDifficulty;
    }
    if (alphabet.length <= difficulty) {
        throw ErrMaxDifficulty;
    }
    if (hashed.length === 0) {
        throw ErrEmptyHashed;
    }
    if (masked.length === 0) {
        throw ErrEmptyMasked;
    }
    let p = new Prefix(difficulty);
    do {
        let cand = p.value + masked;
        if ((await hash(cand)) === hashed) {
            return cand;
        }
    } while (p.iterate());
    throw ErrNotFound;
}
