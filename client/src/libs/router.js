import { createRouter, createWebHistory } from 'vue-router'
import Index from '@/views/Index.vue'
import { mutations, actions } from './store.js'
import { TitleSep, TitleSuffix } from './vars.js'

console.log('init router')

const routes = [
    {
        path: '/',
        name: 'Index',
        component: Index
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/views/auth/Login.vue'),
        meta: { requiresAuth: false, title: '登录' }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/auth/Register.vue'),
        meta: { requiresAuth: false, title: '注册' }
    },
    {
        path: '/logout',
        name: 'Logout',
        redirect() {
            mutations.clearLogin()
            actions.removeConfig()
            return '/'
        },
        meta: { requiresAuth: true }
    },
    {
        path: '/ta',
        name: 'Ta',
        component: () => import('@/views/ta/Ta.vue'),
        meta: { requiresAuth: true, title: '她是' }
    },
    {
        path: '/my',
        name: 'Home',
        component: () => import('@/views/home/Home.vue'),
        children: [
            {
                path: 'profile',
                component: () => import('@/views/home/UserProfile.vue')
            },
            {
                path: 'setting',
                component: () => import('@/views/home/UserSetting.vue')
            },
            {
                path: 'image',
                component: () => import('@/views/home/UserImage.vue')
            }
        ],
        meta: { requiresAuth: true }
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'NotFound',
        component: () => import('@/views/public/NotFound.vue')
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from) => {
    if (to.meta.title) {
        document.title = to.meta.title + TitleSep + TitleSuffix
    } else {
        document.title = TitleSuffix
    }
    // 而不是去检查每条路由记录
    // to.matched.some(record => record.meta.requiresAuth)
    if (to.meta.requiresAuth === false && mutations.isLogged()) {
        return { path: '/' }
    }
    if (to.meta.requiresAuth && !mutations.isLogged()) {
        // 此路由需要授权，请检查是否已登录
        // 如果没有，则重定向到登录页面
        return {
            path: '/login',
            // 保存我们所在的位置，以便以后再来
            query: { redirect: to.fullPath }
        }
    }
})

export default router
