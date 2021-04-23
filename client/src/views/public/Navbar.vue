<template>
    <div class="nav">
        <div class="nav-left">
            <router-link to="/">
                <el-image class="logo" :src="logo"></el-image>
            </router-link>
        </div>
        <div class="nav-right">
            <div v-if="isLogin">
                <el-menu
                    mode="horizontal"
                    router
                    :default-active="this.$route.path"
                >
                    <el-menu-item index>
                        <el-avatar
                            class="avatar"
                            :src="avatar"
                            :size="30"
                        ></el-avatar>
                    </el-menu-item>
                    <el-submenu index="/user">
                        <template slot="title">{{ user }}</template>
                        <el-menu-item index="/user/profile">
                            <i class="saintic-icon saintic-icon-user"></i>
                            个人资料
                        </el-menu-item>
                        <el-menu-item index="/user/setting">
                            <i class="saintic-icon saintic-icon-setting"></i>
                            用户设置
                        </el-menu-item>
                        <el-menu-item>-------------</el-menu-item>
                        <el-menu-item index="/logout">
                            <i class="saintic-icon saintic-icon-logoff"></i>
                            登出
                        </el-menu-item>
                    </el-submenu>
                </el-menu>
            </div>
            <div v-else>
                <el-menu
                    mode="horizontal"
                    router
                    :default-active="this.$route.path"
                >
                    <el-menu-item index="/login">
                        <i class="saintic-icon saintic-icon-login"></i> 登录
                    </el-menu-item>
                    <el-menu-item index="/register">
                        <i class="saintic-icon saintic-icon-register"></i> 注册
                    </el-menu-item>
                </el-menu>
            </div>
        </div>
    </div>
</template>

<script>
import { mapState } from '@/libs/store.js'
import defaultLogo from '@/assets/img/logo.png'
import defaultAvatar from '@/assets/img/defaultAvatar.png'

export default {
    name: 'Navbar',
    data() {
        return {
            logo: defaultLogo,
            avatar: defaultAvatar
        }
    },
    computed: {
        ...mapState(['isLogin', 'user'])
    }
}
</script>

<style>
.nav {
    height: 60px;
    position: relative;
    border-radius: 2px;
    font-size: 15px;
    box-sizing: border-box;
}
.nav .nav-left {
    position: absolute;
    left: 0;
    top: 0;
    width: 200px;
}
.nav .nav-right {
    position: absolute;
    right: 0;
    top: 0;
}
.nav .nav-left .logo {
    max-height: 100%;
}
.nav .nav-right .avatar {
    background: #fff;
}
</style>