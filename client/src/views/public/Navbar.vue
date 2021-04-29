<template>
    <div class="nav">
        <div class="nav-left">
            <router-link to="/">
                <el-image class="logo" :src="logo"></el-image>
            </router-link>
        </div>
        <div class="nav-right">
            <el-menu
                mode="horizontal"
                router
                :default-active="this.$route.path"
            >
                <el-menu-item index="/ta">
                    <i class="saintic-icon saintic-icon-goddess"></i> 她是
                </el-menu-item>
                <el-menu-item index v-if="isLogin">
                    <el-avatar
                        class="avatar"
                        :src="avatar"
                        :size="30"
                    ></el-avatar>
                </el-menu-item>
                <el-submenu v-if="isLogin">
                    <template #title>{{ user }}</template>
                    <el-menu-item index="/my/profile">
                        <i class="saintic-icon saintic-icon-user"></i>
                        个人资料
                    </el-menu-item>
                    <el-menu-item index="/my/setting">
                        <i class="saintic-icon saintic-icon-setting"></i>
                        用户设置
                    </el-menu-item>
                    <el-menu-item index="/my/album">
                        <i class="saintic-icon saintic-icon-user-album"></i>
                        我的专辑
                    </el-menu-item>
                    <el-menu-item>-------------</el-menu-item>
                    <el-menu-item index="/logout">
                        <i class="saintic-icon saintic-icon-logoff"></i>
                        登出
                    </el-menu-item>
                </el-submenu>
                <el-menu-item index="/login" v-if="!isLogin">
                    <i class="saintic-icon saintic-icon-login"></i> 登录
                </el-menu-item>
                <el-menu-item index="/register" v-if="!isLogin">
                    <i class="saintic-icon saintic-icon-register"></i> 注册
                </el-menu-item>
            </el-menu>
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
    computed: mapState(['isLogin', 'user'])
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
.el-menu--collapse .el-menu .el-submenu,
.el-menu--popup {
    min-width: 120px !important;
}
</style>