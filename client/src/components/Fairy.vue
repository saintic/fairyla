<template>
    <section class="container">
        <section class="wrapper">
            <article class="entry">
                <header class="post-head">
                    <h2>{{ album.name }}</h2>
                </header>
                <div class="post-info">
                    &nbsp;&nbsp; 所属：{{ album.owner }}
                </div>
                <div
                    style="text-align: center"
                    v-for="f in fairies"
                    :key="f.id"
                >
                    <img :alt="f.desc" :src="f.src" />
                </div>
                <br />
                <section class="bdshare">
                    <div class="info">
                        <span class="category">
                            <i class="saintic-icon saintic-icon-tags"></i>
                            {{ album.label }}
                        </span>
                        <span class="date">
                            <i class="saintic-icon saintic-icon-time"></i>
                            {{ album.ctime }}
                        </span>
                    </div>
                </section>
            </article>
        </section>
    </section>
</template>

<script setup>
import { defineProps } from 'vue'
import { isObject } from '@/libs/util.js'

defineProps({
    album: {
        type: Object,
        required: true
    },
    faries: {
        type: Array,
        required: true,
        validator: (value) => {
            if (!Array.isArray(value)) return false
            for (let v of value) {
                if (!isObject(v)) return false
                if (!v.hasOwnProperty('album_id') || v.hasOwnProperty('src')) {
                    return false
                }
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
.info .comment:hover i,
.info .more:hover i,
.info .category:hover a,
.info .comment:hover a,
.info .more:hover a {
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
