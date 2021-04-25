import { createRouter, createWebHistory } from 'vue-router'
import Index from '@/views/Index.vue'

import { mutations, actions } from './store.js'

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
        meta: { requiresAuth: false }
    },
    {
        path: '/register',
        name: 'Register',
        component: () => import('@/views/auth/Register.vue'),
        meta: { requiresAuth: false }
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
    }
    /*
    {
        path: '/control',
        name: 'Control',
        component: () => import('@/views/admin/Control.vue'),
        children: [
            {
                path: 'setting',
                component: () => import('@/views/admin/AdminSetting.vue')
            },
            {
                path: 'hook',
                component: () => import('@/views/admin/AdminHook.vue')
            },
            {
                path: 'user',
                component: () => import('@/views/admin/AdminUser.vue')
            }
        ]
    },
    {
        path: '/user',
        name: 'Home',
        component: () => import('@/views/user/Home.vue'),
        children: [
            {
                path: 'profile',
                component: () => import('@/views/user/UserProfile.vue')
            },
            {
                path: 'setting',
                component: () => import('@/views/user/UserSetting.vue')
            },
            {
                path: 'image',
                component: () => import('@/views/user/UserImage.vue')
            }
        ]
    }
    */
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from) => {
    // 而不是去检查每条路由记录
    // to.matched.some(record => record.meta.requiresAuth)
    if (to.meta.requiresAuth === false && mutations.isLogged) {
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
