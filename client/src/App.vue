<template>
    <el-container>
        <el-header>
            <Navbar />
        </el-header>
        <el-main>
            <router-view />
        </el-main>
        <el-footer>
            <Footer />
        </el-footer>
    </el-container>
</template>

<script>
import Navbar from './views/public/Navbar.vue'
import Footer from './views/public/Footer.vue'
import { setStorage } from './libs/util.js'

export default {
    name: 'App',
    components: { Navbar, Footer },
    created() {
        //获取全局基本配置
        this.$store.actions.fetchConfig()
        //在页面刷新时持久化状态数据
        window.addEventListener('beforeunload', (e) => {
            setStorage({ ...this.$store.state })
            e.returnValue = ''
        })
    }
}
</script>

<style>
#app {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
}
body > .el-container {
    margin-bottom: 40px;
}
.el-header,
.el-footer {
    line-height: 60px;
}
</style>
