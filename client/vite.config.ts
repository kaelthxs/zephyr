import { defineConfig } from "vite";

export default defineConfig({
  optimizeDeps: {
    include: ["cesium"],
  },
  build: {
    rollupOptions: {
      external: [],
    },
  },
});
