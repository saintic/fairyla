<template>
    <section class="container">
        <section class="wrapper">
            <article class="entry" v-if="album.id">
                <header class="post-head">
                    <h2>{{ album.name }}</h2>
                </header>
                <div class="post-info">
                    &nbsp;&nbsp; 所属：{{ album.owner }}
                </div>
                <div class="post-info" v-if="Object.keys(btns).length > 0">
                    <fairy-btn
                        v-for="(btn, index) in Object.values(btns)"
                        :key="index"
                        :btn="btn"
                    ></fairy-btn>
                </div>
                <div class="post-main" v-for="f in fairies" :key="f.id">
                    <video
                        :title="f.desc"
                        :src="f.src"
                        preload="metadata"
                        controls
                        style="top: 0; left: 0; max-width: 90%"
                        v-if="f.is_video"
                    >
                        抱歉，您的浏览器不支持内嵌视频！
                    </video>
                    <el-image
                        :title="f.desc"
                        :src="f.src"
                        :lazy="true"
                        :preview-src-list="urls"
                        :hide-on-click-modal="true"
                        v-else
                    >
                    </el-image>
                </div>
                <el-empty
                    :image-size="150"
                    v-if="fairies.length === 0"
                ></el-empty>
                <br />
                <section class="bdshare">
                    <div class="info">
                        <span class="category">
                            <i class="saintic-icon saintic-icon-catalog"></i>
                            <span v-if="album.source === 'Ta'">公开</span>
                            <span v-else-if="album.source === 'Claim'"
                                >共享</span
                            >
                            <span v-else>专属</span>
                        </span>
                        <span class="label" v-if="album.label">
                            <i class="saintic-icon saintic-icon-tags"></i>&nbsp;
                            <span v-for="l in album.label" :key="l"
                                >{{ l }}
                            </span>
                        </span>
                        <span class="date">
                            <i class="saintic-icon saintic-icon-time"></i>
                            {{ album.cdate || album.ctime }}
                        </span>
                    </div>
                </section>
            </article>
        </section>
        <Backtop />
    </section>
</template>

<script setup>
import { defineProps } from 'vue'
import { isObject } from '@/libs/util.js'
import Backtop from './Backtop.vue'
import FairyBtn from './FairyBtn.vue'

defineProps({
    album: {
        type: Object,
        required: true
    },
    fairies: {
        type: Array,
        required: true,
        validator: (value) => {
            if (!Array.isArray(value)) return false
            for (let v of value) {
                if (!isObject(v)) return false
                if (!v.hasOwnProperty('src')) {
                    return false
                }
            }
            return true
        }
    },
    urls: Array,
    btns: {
        type: Object,
        validator: (value) => {
            if (!isObject(value)) return false
            for (let v of Object.values(value)) {
                if (!isObject(v)) return false
            }
            return true
        }
    }
})
</script>

<style scoped>
.wrapper {
    width: 960px;
    margin: 0 auto;
}
.container .wrapper {
    overflow: hidden;
}
.post-head {
    margin-bottom: 24px;
    padding: 8px 0 4px;
    position: relative;
}
.post-head:before {
    position: absolute;
    content: '';
    height: 1px;
    width: 32px;
    background-color: #f46;
    bottom: 0;
    left: 50%;
    margin-left: -16px;
}
.post-info {
    margin-bottom: 32px;
}
.post-main {
    text-align: center;
    margin: 5px auto;
}

.entry {
    text-align: center;
}
.entry p {
    font-size: 14px;
}
.entry br {
    display: none;
}
.entry img {
    margin: 4px auto 8px;
    max-width: 100%;
}
.info {
    overflow: hidden;
    color: #999;
    height: auto;
    font-size: 12px;
    margin-bottom: 8px;
}
.info span {
    margin-right: 4px;
    display: inline-block;
}
.info span i {
    font-size: 14px;
}
.entry .info a {
    color: #999;
}

.info a:hover {
    color: #e02b57;
}
.info i {
    color: #f46;
}
.info .category:hover i,
.info .label:hover i,
.info .date:hover i,
.info .category:hover a,
.info .label:hover a,
.info .date:hover a {
    color: #e02b57;
}

.bdshare {
    margin-top: 32px;
}
.bdshare .bdsharebuttonbox a {
    color: #999 !important;
    display: inline-block !important;
    width: 24px !important;
    height: 24px !important;
    border-radius: 24px;
    border: 1px #999 solid;
    line-height: 22px !important;
    text-align: center;
    background-image: none !important;
    padding-left: 0 !important;
    float: none;
    margin: 6px 1px;
}
.bdshare .bdsharebuttonbox a:hover {
    background-color: #e02b57;
    border: 1px #e02b57 solid;
    color: #fff !important;
}
@media only screen and (max-width: 1280px) {
    .wrapper {
        width: 96% !important;
    }
}
@media only screen and (max-width: 479px) {
    .post-head,
    .post-info {
        margin-bottom: 12px;
    }
}
</style>
