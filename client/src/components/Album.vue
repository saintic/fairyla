<template>
    <section class="container">
        <section class="wrapper">
            <section class="content one_fourth">
                <article
                    class="post effect3"
                    v-for="album in albums"
                    :key="album.id"
                >
                    <router-link
                        :to="{ name: 'Album', params: { name: album.name } }"
                        :title="album.label"
                        class="post-main"
                    >
                        <el-image
                            :src="album.fairy.src"
                            :title="album.fairy.desc"
                            fit="cover"
                        >
                        </el-image>
                        <div class="post-content">
                            <h3>{{ album.name }}</h3>
                            <p>&nbsp;&nbsp; 所属：{{ album.owner }}</p>
                        </div>
                    </router-link>
                </article>
            </section>
        </section>
    </section>
</template>

<script setup>
import { defineProps } from 'vue'
import { isObject } from '@/libs/util.js'

defineProps({
    albums: {
        type: Array,
        required: true,
        validator: (value) => {
            console.log(value)
            if (!Array.isArray(value)) return false
            for (let v of value) {
                if (!isObject(v)) return false
                if (!v.hasOwnProperty('name') || !v.hasOwnProperty('fairy')) {
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

article,
img {
    margin: 0;
    padding: 0;
    border: 0;
    background: 0 0;
}

article {
    display: block;
}

img {
    width: auto;
    height: auto;
    display: block;
    margin: 0 auto;
}

.post {
    -webkit-box-sizing: border-box;
    -moz-box-sizing: border-box;
    box-sizing: border-box;
}

img,
.post-main > div,
.post-content p,
.post-content h3,
*:before,
*:after {
    -moz-transition: all ease-in-out 0.35s;
    -webkit-transition: all ease-in-out 0.35s;
    -o-transition: all ease-in-out 0.35s;
    -ms-transition: all ease-in-out 0.35s;
    transition: all ease-in-out 0.35s;
}

.content {
    margin: 20px 0;
    font-size: 0;
}

.post {
    font-size: 12px;
    display: inline-block;
    vertical-align: top;
}

.post-main img {
    width: 100%;
    height: auto;
}

.post-main {
    position: relative;
    display: block;
    overflow: hidden;
}

.post-content h3 {
    font-size: 16px;
    color: #fff;
}

.post-content p {
    color: #ddd;
}

.post-content {
    background-color: #f46;
    position: absolute;
    top: 0;
    bottom: 0;
    right: 0;
    left: 0;
    visibility: hidden;
    opacity: 0;
    text-align: center;
}

.post:hover .post-content {
    visibility: visible;
    opacity: 1;
}

.effect3:hover img {
    -webkit-transform: scale(1.2);
    -moz-transform: scale(1.2);
    -ms-transform: scale(1.2);
    -o-transform: scale(1.2);
    transform: scale(1.2);
}

.effect3 .post-content {
    background-color: rgba(0, 0, 0, 0.6);
}

.effect3 .post-content h3 {
    -webkit-transform: translateY(-100%);
    -moz-transform: translateY(-100%);
    -ms-transform: translateY(-100%);
    -o-transform: translateY(-100%);
    transform: translateY(-100%);
    background-color: #e02b57;
    margin-top: 16%;
}

.effect3 .post-content p {
    -webkit-transform: translateY(100%);
    -moz-transform: translateY(100%);
    -ms-transform: translateY(100%);
    -o-transform: translateY(100%);
    transform: translateY(100%);
}

.effect3:hover .post-content h3,
.effect3:hover .post-content p {
    -webkit-transform: translateY(0);
    -moz-transform: translateY(0);
    -ms-transform: translateY(0);
    -o-transform: translateY(0);
    transform: translateY(0);
}

.one_fourth .post-content h3,
.one_fourth .post-content p {
    padding: 8px;
}

.one_fourth {
    margin-right: -24px;
}

.one_fourth .post {
    width: 25%;
    padding-right: 24px;
    margin-bottom: 24px;
}

.one_fourth > div {
    margin-right: 24px;
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

@media only screen and (max-width: 1280px) {
    .wrapper {
        width: 96% !important;
    }
}

@media only screen and (max-width: 767px) {
    .one_fourth .post {
        width: 50%;
    }
}

@media only screen and (max-width: 479px) {
    .one_fourth .post {
        width: 50%;
    }

    .one_fourth {
        margin-right: -8px;
    }

    .one_fourth .post {
        margin-bottom: 8px;
        padding-right: 8px;
    }

    .one_fourth > div {
        margin-right: 8px;
    }

    .post-head,
    .post-info {
        margin-bottom: 12px;
    }
}
</style>
