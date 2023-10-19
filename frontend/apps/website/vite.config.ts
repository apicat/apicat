import { fileURLToPath, URL } from 'url'
import { resolve } from 'node:path'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import UnoCSS from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import Icons from 'unplugin-icons/vite'
import IconsResolver from 'unplugin-icons/resolver'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import vueJsx from '@vitejs/plugin-vue-jsx'
import { createHtmlPlugin } from 'vite-plugin-html'
import copy from 'rollup-plugin-copy'

export default ({ mode }) => {
  const isProd = mode === 'production'

  return defineConfig({
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
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        globs: './src/components/*',
        dts: './src/typings/components.d.ts',
        resolvers: [
          ElementPlusResolver({
            importStyle: 'sass',
          }),
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
          {
            entry: 'src/main.ts',
            filename: 'index.html',
            template: 'index.html',
          },
          {
            entry: 'src/dbConfigMain.ts',
            filename: 'db-config.html',
            template: 'db-config.html',
            injectOptions: {
              data: {
                injectDBConfig: isProd ? '<script type="text/javascript">var DB_CONFIG = {{.db_config}}</script>' : ''
              }
            }
          },
        ]
      }),
      copy({
        hook: 'writeBundle',
        verbose: true,
        targets: [
          {
            src: resolve(__dirname, '../../dist/*.html'),
            rename: (name) => `${name}.tmpl`,
            dest: resolve(__dirname, '../../dist/templates/')
          }
        ]
      })
    ],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    css: {
      preprocessorOptions: {
        scss: {
          charset: false,
          additionalData: `@use "@/styles/element/index.scss" as *;`,
        },
      },
    },
    server: {
      open: true,
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
      outDir: '../../dist',
      emptyOutDir: true,
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (id.includes('element-plus')) {
              return 'element-plus'
            }

            if (id.includes('@codemirror')) {
              return 'codemirror'
            }
          },
        },
      },
    },
  })
}
