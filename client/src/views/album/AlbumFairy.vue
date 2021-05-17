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
        shareText() {
            // 当前共享名称
            return this.album.ta ? `共享(@${this.album.ta})` : '共享'
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
                ? `/album?fairy=true&user=${user}&album=${name}`
                : `/user/album/${name}?fairy=true`
        if (!name || (this.fromTa === true && !user)) {
            this.$router.go(-1)
        }
        this.$http.get(url).then((res) => {
            this.album = res.data
            this.album.cdate = formatUnixTimestamp(res.data.ctime)
            this.fairies = res.data.fairy
            delete this.album.fairy
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
                        name: this.shareText,
                        type: 'success',
                        click: () => {
                            this.$prompt(
                                '请输入共享给Ta的用户名（覆盖已有共享）',
                                '温馨提示',
                                {
                                    customClass: 'el-message-box--slim'
                                }
                            ).then(({ value }) => {
                                if (!value) {
                                    this.$message.error({
                                        message: '请输入用户名',
                                        customClass: 'el-message--slim'
                                    })
                                    return false
                                }
                                this.$http
                                    .put(
                                        `/user/album/${this.album.id}?action=share`,
                                        {
                                            ta: value
                                        }
                                    )
                                    .then((res) => {
                                        this.$message.success({
                                            message: '已共享',
                                            customClass: 'el-message--slim'
                                        })
                                        this.album.ta = value
                                        this.btns[0].name = this.shareText
                                    })
                            })
                        }
                    },
                    {
                        name: this.statusText,
                        type: 'warning',
                        click: () => {
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
