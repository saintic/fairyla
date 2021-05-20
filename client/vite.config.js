import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import styleImport from 'vite-plugin-style-import'
import { resolve, join } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        styleImport({
            libs: [
                {
                    libraryName: 'element-plus',
                    esModule: true,
                    ensureStyleFile: true,
                    resolveStyle: (name) => {
                        return `element-plus/lib/theme-chalk/${name}.css`
                    },
                    resolveComponent: (name) => {
                        return `element-plus/lib/${name}`
                    }
                }
            ]
        })
    ],
    server: {
        proxy: {
            '/api': 'http://127.0.0.1:10210'
        }
    },
    resolve: {
        alias: {
            '@': resolve(__dirname, './src')
        }
    },
    build: {
        outDir: join(__dirname, '../server/ui')
    }
})
