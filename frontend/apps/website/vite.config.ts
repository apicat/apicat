/* eslint-disable regexp/no-dupe-disjunctions */
/* eslint-disable regexp/no-unused-capturing-group */
import { URL, fileURLToPath } from 'node:url'
import path, { resolve } from 'node:path'
import fs from 'node:fs'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import UnoCSS from 'unocss/vite'
import { createHtmlPlugin } from 'vite-plugin-html'
import copy from 'rollup-plugin-copy'
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'

// https://vitejs.dev/config/
export default ({ mode }: { mode: string }) => {
  const isBuild = mode === 'production'

  // rm dist dir
  if (isBuild) {
    const dir = './dist'
    const appDist = '../../dist'
    if (fs.existsSync(dir))
      fs.rmSync(dir, { recursive: true })

    if (fs.existsSync(appDist))
      fs.rmSync(appDist, { recursive: true })
  }

  return defineConfig({
    envDir: '../../',
    plugins: [
      vue(),
      vueJsx(),
      VueI18nPlugin({
        runtimeOnly: false,
      }),
      UnoCSS(),
      AutoImport({
        imports: ['vue', 'vue-router', '@vueuse/core'],
        dts: './src/typings/auto-imports.d.ts',
      }),
      Components({
        globs: './src/components/*',
        dts: './src/typings/components.d.ts',
        resolvers: [
          IconsResolver({
            prefix: 'ac-icon',
          }),
        ],
      }),
      Icons({
        autoInstall: true,
      }),
      createHtmlPlugin({
        minify: false,
        pages: [
          // main
          {
            entry: path.resolve(__dirname, './src/main.ts'),
            filename: 'index.html',
            template: 'index.html',
          },
        ],
      }),
      copy({
        hook: 'writeBundle',
        verbose: true,
        targets: [
          {
            src: resolve(__dirname, './dist/*.html'),
            rename: (name: string) => `${name}.tmpl`,
            dest: resolve(__dirname, './dist/templates/'),
          },
          {
            src: resolve(__dirname, './public/*'),
            dest: resolve(__dirname, './dist/assets/'),
          },
          {
            src: resolve(__dirname, './dist/'),
            dest: resolve(__dirname, '../../'),
          },
        ],
      }),
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      host: '0.0.0.0',
      open: false,
      proxy: {
        '/api': {
          target: 'http://127.0.0.1:8000',
          changeOrigin: true,
        },
        '/mock': {
          target: 'http://127.0.0.1:8000',
          changeOrigin: true,
        },
      },
    },
    define: {
      __VUE_I18N_FULL_INSTALL__: false,
      __VUE_I18N_LEGACY_API__: false,
      __INTLIFY_PROD_DEVTOOLS__: false,
    },
    build: {
      minify: true,
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes('element-plus'))
              return 'element-plus'

            if (id.includes('/ac-editor/'))
              return 'ac-editor'

            if (id.includes('prosemirror'))
              return 'prosemirror'

            if (/@codemirror\/(view|state|commands|autocomplete|language)/.test(id))
              return 'codemirror'

            if (/node_modules\/(@vue|vue|vue-router|pinia)/.test(id))
              return 'framework'
          },
        },
      },
    },
  })
}
