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
import App from './App.vue'
import router from './libs/router.js'
import store from './libs/store.js'
import { http } from './libs/util.js'

import {
    ElAvatar,
    ElBadge,
    ElButton,
    ElButtonGroup,
    ElCol,
    ElContainer,
    ElCheckbox,
    ElCheckboxButton,
    ElEmpty,
    ElFooter,
    ElForm,
    ElFormItem,
    ElHeader,
    ElIcon,
    ElImage,
    ElInput,
    ElLink,
    ElMain,
    ElMenu,
    ElMenuItem,
    ElMenuItemGroup,
    ElOption,
    ElOptionGroup,
    ElRow,
    ElResult,
    ElSelect,
    ElSubmenu,
    ElTag,
    ElTooltip,
    ElUpload,
    ElMessage,
    ElMessageBox,
    ElNotification
} from 'element-plus'

const components = [
    ElAvatar,
    ElBadge,
    ElButton,
    ElButtonGroup,
    ElCol,
    ElContainer,
    ElCheckbox,
    ElCheckboxButton,
    ElEmpty,
    ElFooter,
    ElForm,
    ElFormItem,
    ElHeader,
    ElIcon,
    ElImage,
    ElInput,
    ElLink,
    ElMain,
    ElMenu,
    ElMenuItem,
    ElMenuItemGroup,
    ElOption,
    ElOptionGroup,
    ElRow,
    ElResult,
    ElSelect,
    ElSubmenu,
    ElTag,
    ElTooltip,
    ElUpload
]

const plugins = [ElMessage, ElMessageBox, ElNotification]

const app = createApp(App)

components.forEach((component) => {
    app.component(component.name, component)
})

plugins.forEach((plugin) => {
    app.use(plugin)
})

app.use(router)
app.config.globalProperties.$http = http
app.config.globalProperties.$store = store
app.config.globalProperties.$ELEMENT = { size: 'mini' }
let vm = app.mount('#app')
//window.app = app
//window.vm = vm
