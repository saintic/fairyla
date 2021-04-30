<template>
    <Fairy :album="album" :fairies="fairies" :urls="urls" />
</template>

<script>
import Fairy from '@/components/Fairy.vue'
import { formatUnixTimestamp } from '@/libs/util.js'
export default {
    Name: 'Home',
    components: { Fairy },
    data() {
        return { album: {}, fairies: [], urls: [] }
    },
    created() {
        let name = this.$route.params.name
        this.$http
            .get(`/user/album/${name}?is_name=true&fairy=true`)
            .then((res) => {
                this.album = res.data
                this.album.cdate = formatUnixTimestamp(res.data.ctime)
                // TODO 排序
                this.fairies = res.data.fairy
                for (let f of this.fairies) {
                    this.urls.push(f.src)
                }
            })
    }
}
</script>
