import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

// The dev server proxies API and WebSocket traffic to the Go backend so the
// frontend can use same-origin relative URLs (/api, /ws).
//
// FRONTEND_PORT  — port the Vite dev server listens on (default 5174)
// BACKEND_URL    — backend origin the proxy forwards to (default http://localhost:8081)
//
// Both values are read from the .env file in this directory (copy .env.example
// to .env to customise). Real environment variables take precedence because
// Vite's loadEnv prefers process.env over the file.
export default defineConfig(({ mode }) => {
  // "." resolves to the project root (cwd) without needing Node typings.
  const env = loadEnv(mode, ".", "");

  const port = env.FRONTEND_PORT ? parseInt(env.FRONTEND_PORT, 10) : 5174;
  const backendUrl = env.BACKEND_URL || "http://localhost:8081";

  return {
    plugins: [react()],
    server: {
      port,
      proxy: {
        "/api": {
          target: backendUrl,
          changeOrigin: true,
        },
        "/ws": {
          target: backendUrl,
          ws: true,
          changeOrigin: true,
        },
      },
    },
  };
});
