/*
   Copyright 2021 Hiroshi.tao

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/lib/theme-chalk/index.css'
import locale from 'element-plus/lib/locale/lang/zh-cn'
import App from './App.vue'
import router from './libs/router.js'
import store from './libs/store.js'
import { http } from './libs/util.js'

var app = createApp(App)
app.use(router)
app.use(ElementPlus, { locale })
app.config.globalProperties.$http = http
app.config.globalProperties.$store = store
var vm = app.mount('#app')
window.app = app
window.vm = vm
