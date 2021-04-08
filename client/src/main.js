import Vue from 'vue'
import App from './App.vue'
import ElementUI from 'element-ui'
import router from './libs/router.js'
import store from './libs/store.js'
import { http } from './libs/util.js'

Vue.use(ElementUI)
Vue.prototype.$http = http
Vue.prototype.$store = store
Vue.config.productionTip = false

const vm = new Vue({
    router,
    render: (h) => h(App)
}).$mount('#app')
console.log(vm)
