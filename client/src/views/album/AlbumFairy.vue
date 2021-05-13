<template>
    <Fairy :album="album" :fairies="fairies" :urls="urls" :btns="btns" />
</template>

<script>
import Fairy from '@/components/Fairy.vue'
import { formatUnixTimestamp } from '@/libs/util.js'

export default {
    Name: 'AlbumFairy',
    components: { Fairy },
    props: {
        isPublic: Boolean // 来源于Ta为True，来源于Home为False
    },
    data() {
        return { album: {}, fairies: [], urls: [], btns: [] }
    },
    created() {
        let name = this.$route.params.name,
            url = this.isPublic
                ? `/album?fairy=true&album_name=${name}`
                : `/user/album/${name}?fairy=true`
        if (!name) {
            this.$router.go(-1)
        }
        this.$http.get(url).then((res) => {
            let data = this.isPublic ? res.data[0] : res.data
            this.album = data
            this.album.cdate = formatUnixTimestamp(data.ctime)
            this.fairies = data.fairy
            for (let f of this.fairies) {
                if (!f.is_video) {
                    this.urls.push(f.src)
                }
            }
        })
    }
}
</script>
