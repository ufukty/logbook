import * as vitest from "vitest";
import * as solve from "./solve";
import set from "./testdata/set.json";

vitest.describe("Encode", () => {
    vitest.it("should pass", (t) => {
        const input = "Hello world";
        const expected = "SGVsbG8gd29ybGQ";
        const got = solve.encode(input);
        t.expect(got).toBe(expected);
    });
});

vitest.describe("SolveChallenge", async () => {
    // cases are generated with the unit test of server part (Go's) named TestSolve function
    const cases = [
        {
            "Que": "1CWQ9cq2DHZxKhDenM9x8koTp4",
            "Hash": "HNR6LX/LmxHCOv5Q+w+RiHpGmQxST9Ra1jVTCJ+qZcI",
            "Original": "X1CWQ9cq2DHZxKhDenM9x8koTp4",
            "N": 1,
        },
        {
            "Que": "h8jMfgWzkP2QB05DEbSmfWWjY",
            "Hash": "mnmfQ6I74f6vH0h79n3gwjMJy/5KXnAmAkbC2noWi8k",
            "Original": "5Sh8jMfgWzkP2QB05DEbSmfWWjY",
            "N": 2,
        },
        {
            "Que": "5hMxnnXx7tmvrqS+iGJxRzcf2abIPUWH/gF5MF2oZ7HMdk49qN5eiqyJ7AXBPYuAFM",
            "Hash": "g2YwDN3PDiW+SYYll0wCGf9wJOdtCIowqBJ3fvqXg0w",
            "Original": "D5hMxnnXx7tmvrqS+iGJxRzcf2abIPUWH/gF5MF2oZ7HMdk49qN5eiqyJ7AXBPYuAFM",
            "N": 1,
        },
        {
            "Que": "A8ztrr/m2wExSmFw0NQjvPc6EyoXFGaVL03q+SAg6g7u8Y6OAoLKdDe5FkAdgrztc",
            "Hash": "ufe4/0TDcc0n9UGJH/acXBp28xrQUcHPJdiqLIHuY3s",
            "Original": "iuA8ztrr/m2wExSmFw0NQjvPc6EyoXFGaVL03q+SAg6g7u8Y6OAoLKdDe5FkAdgrztc",
            "N": 2,
        },
        {
            "Que": "sR9m1MigEFnPoZrriM3Co5DcAIeey1/XRL4oKXIrEUvMp2LMAFk4pgck30i4oFyo2tED6bBp6srJzTnEet9jk5Z1vWzmeKGEA19NRXQ6slW29Hb0hCuCnO9TQ/FtH6+AfoB0w",
            "Hash": "jMQyLZ3Q/nH2S1Pjs9Nu3TEjg9lvQbkZdsAzjVDioMo",
            "Original":
                "KsR9m1MigEFnPoZrriM3Co5DcAIeey1/XRL4oKXIrEUvMp2LMAFk4pgck30i4oFyo2tED6bBp6srJzTnEet9jk5Z1vWzmeKGEA19NRXQ6slW29Hb0hCuCnO9TQ/FtH6+AfoB0w",
            "N": 1,
        },
        {
            "Que": "l/CX1kpz/UbH/WQYrwSporJBHsduJKUDjnZws1raguE5MvldDpLJV61Dz1thGtwJUvXACYKk/y0krFjt+pVEB1ykDFNRDvtDVPXqF5YQzUZYfpPjzRFIw5DphnMS1vk3Ya3w",
            "Hash": "dKfPxGWPb3Qf5L7jdTisMBDE3x7UukPUmmq0B33kQZo",
            "Original":
                "QMl/CX1kpz/UbH/WQYrwSporJBHsduJKUDjnZws1raguE5MvldDpLJV61Dz1thGtwJUvXACYKk/y0krFjt+pVEB1ykDFNRDvtDVPXqF5YQzUZYfpPjzRFIw5DphnMS1vk3Ya3w",
            "N": 2,
        },
        {
            "Que": "OUJfBGglWVWgiK6XuEqgeeetx1NKxevCo7e46MBzlOYE4XtpJscD9wpZqrAKhygIi0oB8fTMrn+++gPqh1IEKbW2coy5vPYXRB4vps+PhWw9w/wsmpLJLFNT6sXNLLcqUXl1Nn56SWMRnrs17q2ruxWHj5Vmttun7LbehPRyz5fIKP9MWbLCvKU0MlV3cEh295Su7jnjv9huDA1MFYpFVdDinxc7hD0vuSVmSUTb5+CaQvQal6494kC8XcYu97g8edMAGMICIk",
            "Hash": "0pPASSvreRwEe+tDG5Der70Xk7eN+9i+6s6dB/8RI70",
            "Original":
                "KOUJfBGglWVWgiK6XuEqgeeetx1NKxevCo7e46MBzlOYE4XtpJscD9wpZqrAKhygIi0oB8fTMrn+++gPqh1IEKbW2coy5vPYXRB4vps+PhWw9w/wsmpLJLFNT6sXNLLcqUXl1Nn56SWMRnrs17q2ruxWHj5Vmttun7LbehPRyz5fIKP9MWbLCvKU0MlV3cEh295Su7jnjv9huDA1MFYpFVdDinxc7hD0vuSVmSUTb5+CaQvQal6494kC8XcYu97g8edMAGMICIk",
            "N": 1,
        },
        {
            "Que": "ITrhN8CFXsenzEAu+LLKPvVOy53ZaDadbOC52GHjdu9Bfj6a0slG0eCTVYdtI+bmlO1k7fAKV1BCmh+SdkvqCECxdKl+rPYqyuzcOhtq2MeUdb1tTp6Pr6qUnvbn5cQPLfcXHqGwZEcSS1Ig3VeunXnsif8mT88IqvSqyEwDXjqSMT0nSwQiGqeyILr1vF3Rxt1/0z9pZKTkMAHMugDw/Wg7P8N2dqlfEoEhZAR/VTGy45o8Jlh9DMdoOF9Rclc03Nz3+QwyM",
            "Hash": "P1xBPNipzB+Aqa8owC/nuj1K5H5UH0uTvbBQ8zdpP44",
            "Original":
                "/uITrhN8CFXsenzEAu+LLKPvVOy53ZaDadbOC52GHjdu9Bfj6a0slG0eCTVYdtI+bmlO1k7fAKV1BCmh+SdkvqCECxdKl+rPYqyuzcOhtq2MeUdb1tTp6Pr6qUnvbn5cQPLfcXHqGwZEcSS1Ig3VeunXnsif8mT88IqvSqyEwDXjqSMT0nSwQiGqeyILr1vF3Rxt1/0z9pZKTkMAHMugDw/Wg7P8N2dqlfEoEhZAR/VTGy45o8Jlh9DMdoOF9Rclc03Nz3+QwyM",
            "N": 2,
        },
    ];

    for (const c of cases) {
        vitest.it(`should pass for length=${c.Original.length} n=${c.N}`, async (t) => {
            await t.expect(solve.Solve(c.N, c.Que, c.Hash)).resolves.toBe(c.Original);
        });
    }

    vitest.it("should fail with 'invalid challange'", async (t) => {
        await t.expect(solve.Solve(1, "", "")).rejects.toThrow(solve.ErrInvalidChallange);
    });

    vitest.it("should fail with 'not found'", async (t) => {
        await t.expect(solve.Solve(1, "", "definitely not empty")).rejects.toThrow(solve.ErrNotFound);
    });
});

vitest.describe("Measure solution time", async () => {
    const ts: number[] = [];

    vitest.it("should some take time", { timeout: 5 * 100 * 1000 }, async (t) => {
        for (const c of set) {
            const t = performance.now();
            await solve.Solve(c.N, c.Cue, c.Hash);
            const d = performance.now() - t;
            ts.push(d);
            const tot = ts.reduce((a, b) => a + b);
            console.log(`took: ${d / 1000}s (tot: ${tot / 1000}s, avg: ${(tot / 1000 / ts.length).toFixed(0)}s)`);
        }
        const tot = ts.reduce((a, b) => a + b);
        t.expect(tot).to.greaterThan(0);
    });
});
