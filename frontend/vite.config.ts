import react from "@vitejs/plugin-react";
import path from "path";
import { defineConfig } from "vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";

export default defineConfig({
  plugins: [
    tanstackRouter({
      routesDirectory: path.resolve(__dirname, "./src/renderer/routes"),
      generatedRouteTree: path.resolve(__dirname, "./src/renderer/routeTree.gen.ts"),
      routeFileIgnorePrefix: "-",
    }),
    react(),
  ],
  base: "./",
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src/renderer"),
    },
  },
  server: {
    port: 3000,
    host: '127.0.0.1',
  },
});

