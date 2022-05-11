/* eslint-disable no-undef */
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import autoExternal from 'rollup-plugin-auto-external'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vitejs.dev/config/
export default defineConfig({
    build: {
        lib: {
            entry: path.resolve(__dirname, './src/index.js'),
            formats: ['es'],
            name: 'AcEditor',
        },
        rollupOptions: {
            plugins: [
                autoExternal({
                    dependencies: false,
                    packagePath: path.resolve(__dirname, './package.json'),
                }),
            ],
        },
    },
    plugins: [vue(), visualizer()],
})
