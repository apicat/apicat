import { fileURLToPath, URL } from 'url'
import { defineConfig, splitVendorChunkPlugin } from 'vite'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import copy from 'rollup-plugin-copy'
import del from 'rollup-plugin-delete'
import { visualizer } from 'rollup-plugin-visualizer'

export default defineConfig(({ mode }) => {
    const plugins = [
        splitVendorChunkPlugin(),
        vue(),
        vueJsx(),
        AutoImport({
            resolvers: [ElementPlusResolver()],
        }),
        Components({
            resolvers: [ElementPlusResolver()],
        }),
        visualizer(),
    ]

    if (mode === 'production') {
        plugins.push(
            del({ targets: '../../../public/assets/*', verbose: true, force: true }),
            del({ targets: '../../../public/static/*', verbose: true, force: true }),
            copy({
                targets: [
                    {
                        src: 'dist/index.html',
                        dest: '../../views/',
                        rename: 'index.blade.php',
                    },
                    {
                        src: ['dist/*', '!*/*.html'],
                        dest: '../../../public/',
                    },
                ],
                verbose: true,
                hook: 'closeBundle',
            })
        )
    }

    return {
        plugins,
        resolve: {
            alias: {
                '@': fileURLToPath(new URL('./src', import.meta.url)),
            },
        },

        css: {
            preprocessorOptions: {
                // https://github.com/vitejs/vite/discussions/5079
                css: { charset: false },
                scss: { charset: false },
            },
        },

        build: {
            rollupOptions: {
                output: {
                    manualChunks(id) {
                        if (id.includes('element-plus')) {
                            return 'element-plus'
                        }
                    },
                },
            },
        },

        server: {
            proxy: {
                '/api': {
                    target: 'http://192.168.50.61:8000',
                    changeOrigin: true,
                },
            },
        },
    }
})
