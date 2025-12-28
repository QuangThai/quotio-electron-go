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
    emptyOutDir: true,
    target: "node18",
    lib: {
      entry: path.resolve(__dirname, "src/main/main.ts"),
      formats: ["cjs"],
      fileName: "main",
    },
    rollupOptions: {
      external: builtins,
      output: {
        entryFileNames: "[name].cjs",
        format: "cjs",
      },
    },
  },
});

