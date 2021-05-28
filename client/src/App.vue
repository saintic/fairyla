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

export default {
    name: 'App',
    components: { Navbar, Footer },
    created() {
        this.$store.actions.fetchConfig()
        window.addEventListener('beforeunload', (e) => {
            this.$store.actions.saveConfig2Local()
        })
    },
    mounted() {
        let url = '/api/user/event'
        let es = new EventSource(url)
        es.addEventListener('message', (event) => {
            console.log(event.data)
        })
        es.addEventListener('error', (event) => {
            if (event.readyState == EventSource.CLOSED) {
                console.log('event was closed')
            }
        })
        es.addEventListener('close', (event) => {
            console.log(event.type)
            es.close()
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
.el-message--slim {
    min-width: 150px;
}
.el-message-box--slim {
    width: 300px;
}
</style>
