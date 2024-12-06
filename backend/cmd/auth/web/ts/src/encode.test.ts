import * as vitest from "vitest";
import { encode } from "./encode";

vitest.describe("Encode", () => {
    vitest.it("should pass", (t) => {
        const input = "Hello world";
        const expected = "SGVsbG8gd29ybGQ";
        const got = encode(input);
        t.expect(got).toBe(expected);
    });
});
