import * as config from "vitest/config";

export default config.defineConfig({
    test: {
        environment: "happy-dom",
        testTimeout: 30 * 1000,
    },
});
