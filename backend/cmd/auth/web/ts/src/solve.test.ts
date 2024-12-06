import * as vitest from "vitest";
import * as solve from "./solve";

// cases are generated with the unit test of server part (Go's) named TestSolve function
import d10 from "./testdata/d10.json";
import d20 from "./testdata/d20.json";
import d30 from "./testdata/d30.json";
import d40 from "./testdata/d40.json";
import d50 from "./testdata/d50.json";
import d60 from "./testdata/d60.json";

vitest.describe("SolveChallenge", async () => {
    vitest.it(`should fail with '${solve.ErrEmptyMasked}'`, async (t) => {
        await t.expect(solve.Solve(2, "", "not empty")).rejects.toThrow(solve.ErrEmptyMasked);
    });
    vitest.it(`should fail with '${solve.ErrEmptyHashed}'`, async (t) => {
        await t.expect(solve.Solve(2, "not empty", "")).rejects.toThrow(solve.ErrEmptyHashed);
    });
    vitest.it(`should fail with '${solve.ErrMinDifficulty}'`, async (t) => {
        await t.expect(solve.Solve(1, "not empty", "not empty")).rejects.toThrow(solve.ErrMinDifficulty);
    });
    vitest.it(`should fail with '${solve.ErrMaxDifficulty}'`, async (t) => {
        await t.expect(solve.Solve(63, "not empty", "not empty")).rejects.toThrow(solve.ErrMaxDifficulty);
    });
    vitest.it("should fail with 'not found'", async (t) => {
        await t.expect(solve.Solve(2, "not empty", "definitely not empty")).rejects.toThrow(solve.ErrNotFound);
    });
});

const ms = 1;
const second = 1000 * ms;

function subSlice<T>(array: T[], start: number, end: number): T[] {
    const adjustedStart = Math.max(0, start);
    const adjustedEnd = Math.min(array.length, end);
    const result: T[] = [];
    for (let i = adjustedStart; i < adjustedEnd; i++) {
        result.push(array[i]);
    }
    return result;
}

vitest.describe("Measure solution time", async () => {
    const dts = {
        "d10": d10,
        "d20": d20,
        "d30": d30,
        "d40": d40,
        "d50": d50,
        "d60": d60,
    };

    const samplesize = 10;
    Object.entries(dts).forEach(([dtn, dt]) => {
        vitest.it(dtn, { timeout: 60 * second }, async (t) => {
            const ts: number[] = [];
            for (const c of subSlice(dt, 0, samplesize)) {
                const t = performance.now();
                await solve.Solve(c.D, c.M, c.H);
                const d = performance.now() - t;
                ts.push(d);
                // const tot = ts.reduce((a, b) => a + b);
                // console.log(`took: ${d}ms (tot: ${tot}ms, avg: ${(tot / ts.length).toFixed(0)}ms)`);
            }
            const tot = ts.reduce((a, b) => a + b);
            t.expect(tot).to.greaterThan(0);

            const min = Math.min(...ts).toFixed(0);
            const max = Math.max(...ts).toFixed(0);
            const avg = (tot / samplesize).toFixed(0);
            console.log(`=> stats for ${dtn} (min: ${min}ms) (max: ${max}ms) (avg: ${avg}ms)`);
        });
    });
});
