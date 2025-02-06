import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import commonjs from "vite-plugin-commonjs";
import { compression } from "vite-plugin-compression2";

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), commonjs(), compression()],
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes("node_modules")) {
            const modulePath = id.split("node_modules/")[1];
            const topLevelFolder = modulePath.split("/")[0];
            if (topLevelFolder !== ".pnpm") {
              return topLevelFolder;
            }
            const scopedPackageName = modulePath.split("/")[1];
            const chunkName =
              scopedPackageName.split("@")[
              scopedPackageName.startsWith("@") ? 1 : 0
              ];
            return chunkName;
          }
        },
      },
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://backend:8000",
        changeOrigin: true,
        secure: process.env.NODE_ENV === "production",
      },
    },
  },
});
