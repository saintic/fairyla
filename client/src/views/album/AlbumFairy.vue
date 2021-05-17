<template>
    <Fairy :album="album" :fairies="fairies" :urls="urls" :btns="btns" />
</template>

<script>
import Fairy from '@/components/Fairy.vue'
import { formatUnixTimestamp } from '@/libs/util.js'
import { TaLabel } from '@/libs/vars.js'

export default {
    Name: 'AlbumFairy',
    components: { Fairy },
    props: {
        source: String // 来源于：Ta、Home(default)
    },
    data() {
        return { album: {}, fairies: [], urls: [], btns: [] }
    },
    computed: {
        statusText() {
            // 当前专辑反向状态
            return this.album.public ? '私有' : '公开'
        },
        fromTa() {
            // 来源于Ta为True，其他来源为False
            return this.source === TaLabel
        }
    },
    created() {
        let user = this.$route.params.user,
            name = this.$route.params.name,
            url = this.fromTa
                ? `/album?fairy=true&user=${user}&album_name=${name}`
                : `/user/album/${name}?fairy=true`
        if (!name || (this.fromTa === true && !user)) {
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
            if (this.source === TaLabel) {
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
                            console.log('click status')
                            this.$http
                                .put(`/user/album/${name}?action=status`)
                                .then((res) => {
                                    this.$message.success({
                                        message: `已 <b>${this.statusText}</b> 状态`,
                                        dangerouslyUseHTMLString: true,
                                        customClass: 'el-message--slim'
                                    })
                                    this.album.public = !this.album.public
                                    this.btns[1].name = this.statusText
                                })
                        }
                    }
                ]
            }
        })
    }
}
</script>
