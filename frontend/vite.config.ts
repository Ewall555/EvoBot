import { resolve } from "node:path";
import { URL, fileURLToPath } from "node:url";
import UnoCSS from "unocss/vite";
import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { createSvgIconsPlugin } from "vite-plugin-svg-icons";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";

const prefix = `monaco-editor/esm/vs`;

// https://vitejs.dev/config/
export default defineConfig({
  base: "/ui/",
  server: {
    port: 24121,
    open: false,
    host: "0.0.0.0",
    proxy: {
      "/api/v1": {
        target: "http://localhost:24916",
        changeOrigin: true,
      },
    },
  },
  plugins: [
    vue(),
    UnoCSS(),
    AutoImport({
      resolvers: [
        ElementPlusResolver({
          importStyle: "sass",
        }),
      ],
    }),
    Components({
      resolvers: [
        ElementPlusResolver({
          importStyle: "sass",
        }),
      ],
    }),
    createSvgIconsPlugin({
      iconDirs: [resolve(process.cwd(), "src/icons")],
      inject: "body-first",
    }),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@use "@/assets/styles/element-theme.scss" as *;`,
      },
    },
  },
  build: {
    minify: "esbuild",
    target: "esnext",
    rollupOptions: {
      output: {
        chunkFileNames: "assets/js/[name]-[hash].js",
        entryFileNames: "assets/js/[name]-[hash].js",
        assetFileNames: "assets/[ext]/[name]-[hash].[ext]",
        manualChunks: {
          jsonWorker: [`${prefix}/language/json/json.worker`],
          cssWorker: [`${prefix}/language/css/css.worker`],
          htmlWorker: [`${prefix}/language/html/html.worker`],
          tsWorker: [`${prefix}/language/typescript/ts.worker`],
          editorWorker: [`${prefix}/editor/editor.worker`],
        },
      },
    },
  },
});
