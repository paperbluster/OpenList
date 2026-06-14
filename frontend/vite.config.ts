import path from "path"
import { defineConfig } from "vite"
import solidPlugin from "vite-plugin-solid"
import { dynamicBase } from "vite-plugin-dynamic-base"
import { viteStaticCopy } from "vite-plugin-static-copy"

export default defineConfig({
  resolve: {
    alias: {
      "~": path.resolve(__dirname, "src"),
      // "@solidjs/router": path.resolve(__dirname, "solid-router/src"),
      "solid-icons": path.resolve(__dirname, "node_modules/solid-icons"),
    },
  },
  plugins: [
    solidPlugin(),
    dynamicBase({
      // dynamic public path var string, default window.__dynamic_base__
      // Normalize: strip trailing "/" to prevent "//" (protocol-relative URL)
      // when __dynamic_base__ is "/" or ends with "/"
      publicPath: "(window.__dynamic_base__||'').replace(/\\/$/,'')",
      // dynamic load resources on index.html, default false. maybe change default true
      transformIndexHtml: true,
      transformIndexHtmlConfig: {
        insertBodyAfter: true,
      },
    }),
    viteStaticCopy({
      targets: [
        {
          src: "node_modules/@jellyfin/libass-wasm/dist/js/subtitles-octopus-worker.{js,wasm}",
          dest: "static/libass-wasm",
        },
        {
          src: "src/components/artplayer-plugin-ass/fonts/*",
          dest: "static/fonts",
        },
      ],
    }),
  ],
  base: process.env.NODE_ENV === "production" ? "/__dynamic_base__/" : "/",
  // base: "/",
  build: {
    // target: "es2015", //next
    // polyfillDynamicImport: false,
  },
  // experimental: {
  //   renderBuiltUrl: (filename, { type, hostId, hostType }) => {
  //     if (type === "asset") {
  //       return { runtime: `window.OPENLIST_CONFIG.cdn/${filename}` };
  //     }
  //     return { relative: true };
  //   },
  // },
  server: {
    host: "0.0.0.0",
    proxy: {
      "/api": {
        target: "http://localhost:5244",
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, ""),
      },
    },
  },
})
