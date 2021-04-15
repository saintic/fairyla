import vue from '@vitejs/plugin-vue'
import path from 'path'

/**
 * @type {import('vite').UserConfig}
 */
export default {
    plugins: [vue()],
    server: {
        proxy: {
            '/api': 'http://127.0.0.1:10210/api'
        }
    },
    resolve: {
        alias: {
            '@': path.resolve(__dirname, 'src')
        }
    }
}
