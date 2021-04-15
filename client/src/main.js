import { createApp } from 'vue'
import App from './App.vue'
import ElementPlus from 'element-plus'
import 'element-plus/lib/theme-chalk/index.css'
//import router from './libs/router.js'
import store from './libs/store.js'
import { http } from './libs/util.js'

const app = createApp(App)
    .use(ElementPlus)
    .mount('#app')

app.config.globalProperties.$http = http
app.config.globalProperties.$store = store

console.log(app)
