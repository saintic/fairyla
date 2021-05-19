<template>
    <div>
        <p>公测中，欢迎来撩~</p>
        <br />
        <p>分享与记录你心目中的喜欢的、可爱的人或物</p>
        <p>可以是亲人、恋人、偶像，也可以是动漫、宠物、风景等~</p>
        <p>以专辑为单位，上传各种图片和视频：</p>
        <p>或分享给Ta，共同维护；</p>
        <p>或参观并认领公开的专辑~</p>
        <br />
        <p>甚至还可以参与进来设计功能更新~</p>
        <p>
            技术栈：Go + Vue3 + Element-Plus，{ 联系：
            <el-link
                type="info"
                href="mailto:me@tcw.im"
                icon="saintic-icon saintic-icon-email saintic-icon-1-5"
            ></el-link>
            }
        </p>
        <br />
        <p>{{ poesy }}</p>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

const poesy = ref('')

function setPoesy() {
    axios
        .get('https://open.saintic.com/api/sentence/renwu.nvzi.json')
        .then((result) => {
            let res = result.data
            if (res.code === 0) {
                poesy.value = res.data.sentence
            }
        })
}
setPoesy()
setInterval(setPoesy, 10000)
</script>

<style scoped>
a {
    color: #42b983;
}
</style>
