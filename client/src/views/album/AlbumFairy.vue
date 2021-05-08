<template>
    <Fairy :album="album" :fairies="fairies" :urls="urls" />
</template>

<script>
import Fairy from '@/components/Fairy.vue'
import { formatUnixTimestamp } from '@/libs/util.js'

export default {
    Name: 'AlbumFairy',
    components: { Fairy },
    props: {
        isPublic: Boolean
    },
    data() {
        return { album: {}, fairies: [], urls: [] }
    },
    created() {
        let name = this.$route.params.name,
            url = this.isPublic
                ? `/album?fairy=true&album_name=${name}`
                : `/user/album/${name}?fairy=true`
        this.$http.get(url).then((res) => {
            let data = this.isPublic ? res.data[0] : res.data
            this.album = data
            this.album.cdate = formatUnixTimestamp(data.ctime)
            this.fairies = data.fairy
            for (let f of this.fairies) {
                this.urls.push(f.src)
            }
        })
    }
}
</script>
