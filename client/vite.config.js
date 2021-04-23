import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve, join } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
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
