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
        fromTa: Boolean // 来源于Ta为True，来源于Home为False
    },
    data() {
        return { album: {}, fairies: [], urls: [], btns: [] }
    },
    computed: {
        statusText() {
            return this.album.public ? '公开' : '私有'
        }
    },
    created() {
        let name = this.$route.params.name,
            url = this.fromTa
                ? `/album?fairy=true&album_name=${name}`
                : `/user/album/${name}?fairy=true`
        if (!name) {
            this.$router.go(-1)
        }
        this.$http.get(url).then((res) => {
            // 专辑数据（单一）
            let data = this.fromTa ? res.data[0] : res.data
            this.album = data
            this.album.cdate = formatUnixTimestamp(data.ctime)
            this.fairies = data.fairy
            for (let f of this.fairies) {
                if (!f.is_video) {
                    this.urls.push(f.src)
                }
            }
            // Add function buttons
            if (this.fromTa) {
                this.btns = [
                    {
                        name: '认领',
                        plain: true,
                        type: 'success',
                        click: () => {
                            console.log('click ta')
                        }
                    }
                ]
            } else {
                this.btns = [
                    {
                        name: '共享',
                        type: 'success',
                        click: () => {
                            console.log('click home')
                        }
                    },
                    {
                        name: this.statusText,
                        type: 'warning',
                        click: () => {
                            this.album.public = !this.album.public
                            this.btns[1].name = this.statusText
                            this.btns[1].disabled = true
                            console.log('click status')
                        }
                    }
                ]
            }
        })
    }
}
</script>
