import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      "/auth":     { target: "http://localhost:8080", changeOrigin: true },
      "/vehicles": { target: "http://localhost:8080", changeOrigin: true },
      "/sessions": { target: "http://localhost:8080", changeOrigin: true },
      "/slots":    { target: "http://localhost:8080", changeOrigin: true },
    },
  },
});
