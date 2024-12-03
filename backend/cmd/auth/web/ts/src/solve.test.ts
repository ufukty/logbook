import * as vitest from "vitest";
import * as solve from "./solve";

vitest.describe("Encode", () => {
    vitest.it("should pass", (t) => {
        const input = "Hello world";
        const expected = "SGVsbG8gd29ybGQ";
        const got = solve.encode(input);
        t.expect(got).toBe(expected);
    });
});

vitest.describe("SolveChallange", async () => {
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
            "Que": "NfFN7uD6EX5bbu6FsiUHwDFw",
            "Hash": "add5dX+CQB/Ce/f3Od+PDD7Ta8gv63HO8A/1dhHx2Bo",
            "Original": "LOhNfFN7uD6EX5bbu6FsiUHwDFw",
            "N": 3,
        },
        {
            "Que": "UtV/vQpN76eHj8BNYIWxccg",
            "Hash": "APn0ekN6DcjSel8yduzF7GpBBP71lXecXmO1dVs++LM",
            "Original": "YaDFUtV/vQpN76eHj8BNYIWxccg",
            "N": 4,
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
            "Que": "EZGbyeFba4B/OHW5tQZl+fRc07i9Uml72dLBTtRml7CkQgEJkYBtN3jTyj+5sQbQ",
            "Hash": "ePV02hEH/0DBsMaTsQ2vYJ9Kv3R5XMTq9uwzFH8vDhM",
            "Original": "U3eEZGbyeFba4B/OHW5tQZl+fRc07i9Uml72dLBTtRml7CkQgEJkYBtN3jTyj+5sQbQ",
            "N": 3,
        },
        {
            "Que": "+xENwVkbqwbj8zXFMeQVoLSjIMRSCgRZU6P3l9u2YFrg6EAVnWo20Y1/02akpVM",
            "Hash": "vIEppgzDP3hhTHQw6GLJuoWJd636l/CWDRlBMHQtrQY",
            "Original": "fHSX+xENwVkbqwbj8zXFMeQVoLSjIMRSCgRZU6P3l9u2YFrg6EAVnWo20Y1/02akpVM",
            "N": 4,
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
            "Que": "yx7fjEKWCAPUl54kC9mQI8j7eV+Ak8GoBl6PrHVA8yV7NqGvMt0PA9qFMB2vg7iEbRmpNmgQrWk7zhq1FS4LiooeUr79cuRBZRmu7fBQGWrAB1FH49qJVoW29pEx8KRR7Vw",
            "Hash": "GQAtt9pWjXTF9HT2m9DSOX0iD23snhFtdqeQGfrwc/M",
            "Original":
                "NwEyx7fjEKWCAPUl54kC9mQI8j7eV+Ak8GoBl6PrHVA8yV7NqGvMt0PA9qFMB2vg7iEbRmpNmgQrWk7zhq1FS4LiooeUr79cuRBZRmu7fBQGWrAB1FH49qJVoW29pEx8KRR7Vw",
            "N": 3,
        },
        {
            "Que": "P0+Y4mKZg7PUGTHvSOEC2caESkoo/l/ZZOpJh3TU2Z7uMjk3noTLhcsOUStBoOEXcdZuM56YIDwGjY4AyGnKzxyZQbx0mIsN+DEb2rF6oWiujmwJ3LPdoqxiv1HZLLA9bw",
            "Hash": "X6D3hTeBntZrc5zgU+ijlhTY1zLN81Pc3RzQMuNu0i0",
            "Original":
                "YOF4P0+Y4mKZg7PUGTHvSOEC2caESkoo/l/ZZOpJh3TU2Z7uMjk3noTLhcsOUStBoOEXcdZuM56YIDwGjY4AyGnKzxyZQbx0mIsN+DEb2rF6oWiujmwJ3LPdoqxiv1HZLLA9bw",
            "N": 4,
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
        {
            "Que": "oX6VmdmSVb3KuWzyta74Ka7/GZZTqhmR6OBMvd1Xh/Id0Cz5nEKJs3J20y++u7cEPtU1KglRGH7RzVJQYLfxyZ5K2antSY6JUpx80T1Uw+NzuYQwJlOyi5QTqya3cffJ7YOYwqr2c2AG4tKYVpMWp03BX88QmIJ3Ly0k/+PZS3O1xxMDGgIocO622+xCU5uvAUuTmDTc4YM75zS+uN6jhxjnXR/QlmvngHIG4SPEvVm705WYfhffthC1CCLCJdnB5Ev+YBtc",
            "Hash": "eaieIVSU8ZpU352wjBZVbsY46Srzuutl/yVx2HEeQ0E",
            "Original":
                "bI2oX6VmdmSVb3KuWzyta74Ka7/GZZTqhmR6OBMvd1Xh/Id0Cz5nEKJs3J20y++u7cEPtU1KglRGH7RzVJQYLfxyZ5K2antSY6JUpx80T1Uw+NzuYQwJlOyi5QTqya3cffJ7YOYwqr2c2AG4tKYVpMWp03BX88QmIJ3Ly0k/+PZS3O1xxMDGgIocO622+xCU5uvAUuTmDTc4YM75zS+uN6jhxjnXR/QlmvngHIG4SPEvVm705WYfhffthC1CCLCJdnB5Ev+YBtc",
            "N": 3,
        },
        {
            "Que": "QcJpBgaO+2+DzwH56H7iK1qKOSlEKTb/iMr10YMGrEQOdMKp8vv8OyB445BXMDiWu1lxD++Z7HvnqF54lLyAZ2cStE6rT86Mk+LPa0Izu1Tt+9dKehjYXLS6Hp593pSSoYgpbq5mWhuPw+IcqWSpcKhKtsKhlH/7acooLNRHUzN3Onf3n0DLJvCoYidyX+kQqiR4RhBQgO9lVECIGVnyBPmPNrbaYAxcBTeCg6YklUs9OwLd30dkVG4MLRR5zOAzZ09xHgI",
            "Hash": "/ANUxKfislyljRQL4ncKQ3wn5zcMCq5Mx8YpFZbA4Vo",
            "Original":
                "uKffQcJpBgaO+2+DzwH56H7iK1qKOSlEKTb/iMr10YMGrEQOdMKp8vv8OyB445BXMDiWu1lxD++Z7HvnqF54lLyAZ2cStE6rT86Mk+LPa0Izu1Tt+9dKehjYXLS6Hp593pSSoYgpbq5mWhuPw+IcqWSpcKhKtsKhlH/7acooLNRHUzN3Onf3n0DLJvCoYidyX+kQqiR4RhBQgO9lVECIGVnyBPmPNrbaYAxcBTeCg6YklUs9OwLd30dkVG4MLRR5zOAzZ09xHgI",
            "N": 4,
        },
    ];

    for (const c of cases) {
        vitest.it(`should pass for length=${c.Original.length} n=${c.N}`, { timeout: 10 * 60 * 1000 }, async (t) => {
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
