<template>
    <Fairy :album="album" :fairies="fairies" :urls="urls" :btns="btns" />
</template>

<script>
import Fairy from '@/components/Fairy.vue'
import { formatUnixTimestamp } from '@/libs/util.js'
import { TaLabel, ClaimLabel } from '@/libs/vars.js'
import { mapState } from '@/libs/store.js'

export default {
    Name: 'AlbumFairy',
    components: { Fairy },
    props: {
        source: String // 来源于：Ta、Home(default)、Claim
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
            // 当前分享名称
            return this.album.ta ? `分享(@${this.album.ta})` : '分享'
        },
        isTa() {
            // 来源于Ta为True，其他来源为False
            return this.source === TaLabel
        },
        isClaim() {
            // 来源于Claim为True，其他来源为False
            return this.source === ClaimLabel
        },
        ...mapState(['isLogin', 'user'])
    },
    created() {
        let owner = this.$route.params.owner, // 来源于Ta时，此专辑属主
            name = this.$route.params.name, // 专辑名
            url = `/user/album/${name}?fairy=true`
        if (this.isTa) url = `/album?fairy=true&user=${owner}&album=${name}`
        if (this.isClaim) {
            url = `/user/claim?fairy=true&owner=${owner}&album=${name}`
        }
        if (!name || ((this.isTa || this.isClaim) && !owner)) {
            this.$router.go(-1)
        }
        this.$http.get(url).then((res) => {
            this.album = res.data
            this.album.source = this.source
            this.album.cdate = formatUnixTimestamp(res.data.ctime)
            this.fairies = res.data.fairy
            delete this.album.fairy
            for (let f of this.fairies) {
                if (!f.is_video) {
                    this.urls.push(f.src)
                }
            }
            // Add function buttons
            let taBtns = [
                /*
                {
                    name: '认领',
                    plain: true,
                    type: 'success',
                    click: () => {
                        console.log('click ta')
                    }
                }
                */
            ]
            let claimBtns = []
            let homeBtns = [
                {
                    name: this.shareText,
                    type: 'success',
                    click: () => {
                        this.$prompt(
                            '请输入分享给Ta的用户名（覆盖已有分享）',
                            '温馨提示',
                            {
                                customClass: 'el-message-box--slim'
                            }
                        ).then(({ value }) => {
                            let ta = value
                            if (!ta) {
                                this.$message.error({
                                    message: '请输入用户名',
                                    customClass: 'el-message--slim'
                                })
                                return false
                            }
                            if (ta === this.user) {
                                this.$message.error({
                                    message: '不能分享给自己',
                                    customClass: 'el-message--slim'
                                })
                                return false
                            }
                            this.$http
                                .put(`/user/album/${name}?action=share`, { ta })
                                .then((res) => {
                                    this.$message.success({
                                        message: `专辑已分享给${ta}`,
                                        customClass: 'el-message--slim'
                                    })
                                    this.album.ta = ta
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
                                    message: `专辑已 <b>${this.statusText}</b>`,
                                    dangerouslyUseHTMLString: true,
                                    customClass: 'el-message--slim'
                                })
                                this.album.public = !this.album.public
                                this.btns[1].name = this.statusText
                            })
                    }
                }
            ]
            if (this.isLogin) {
                if (this.source === TaLabel) {
                    this.btns = taBtns
                } else if (this.source === ClaimLabel) {
                    this.btns = claimBtns
                } else {
                    this.btns = homeBtns
                }
            }
        })
    }
}
</script>
