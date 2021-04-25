<template>
    <div class="footer">
        <a href="http://www.beian.gov.cn" target="_blank" v-if="beian">
            <img class="beian" :src="beianSrc" />
            {{ beian }} </a
        >&nbsp;<a href="http://beian.miit.gov.cn" target="_blank" v-if="icp">
            {{ icp }}
        </a>
        Copyright &copy; {{ year }} All rights reserved. Powered by
        <a
            href="https://github.com/staugur/fairyla"
            target="_blank"
            ref="nofollow"
        >
            fairyla
        </a>
    </div>
</template>

<script>
import { mapState } from '@/libs/store.js'
import beianSrc from '@/assets/img/beian.png'

export default {
    name: 'Footer',
    data() {
        return {
            beianSrc: beianSrc,
            year: new Date().getFullYear()
        }
    },
    computed: mapState(['icp', 'beian'])
}
</script>

<style>
.footer {
    text-align: center;
}
.footer a {
    color: #666;
    text-decoration: none;
}
.footer a:hover {
    color: #777;
}
.footer .beian {
    vertical-align: middle;
    margin-top: -4px;
}
</style>