import { defineConfig } from "vite";
import path from "path";

// Node.js built-in modules that should be externalized
const builtins = [
  "electron",
  "path",
  "child_process",
  "fs",
  "http",
  "https",
  "url",
  "os",
  "net",
  "tls",
  "crypto",
  "stream",
  "events",
  "util",
  "buffer",
];

export default defineConfig({
  build: {
    outDir: "dist-electron",
    emptyOutDir: false,
    target: "node18",
    lib: {
      entry: path.resolve(__dirname, "src/preload/preload.ts"),
      formats: ["cjs"],
      fileName: "preload",
    },
    rollupOptions: {
      external: builtins,
      output: {
        entryFileNames: "preload.cjs",
        format: "cjs",
      },
    },
  },
});

