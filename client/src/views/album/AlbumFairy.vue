<template>
    <div>
        <Fairy :album="album" :fairies="fairies" :urls="urls" :btns="btns" />
    </div>
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
        return {
            album: {},
            fairies: [],
            urls: [],
            btns: {}
        }
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
        let owner = this.$route.params.owner, // 来源于Ta、Claim时的专辑属主
            name = this.$route.params.name, // 专辑名
            url = `/user/album/${name}?fairy=true`
        if (this.isTa) url = `/album/${owner}/${name}?fairy=true`
        if (this.isClaim) url = `/user/claim/${owner}/${name}?fairy=true`
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
            let taBtns = this.rendeTaBtns(owner, name),
                claimBtns = this.rendeClaimBtns(owner, name),
                homeBtns = this.rendeHomeBtns(name)
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
    },
    methods: {
        rendeTaBtns(owner, albumName) {
            let taBtns = {}
            if (this.isLogin && this.user !== owner) {
                taBtns['taClaim'] = {
                    name: '认领',
                    plain: true,
                    type: 'success',
                    click: () => {
                        this.$http
                            .post('/user/claim', { owner, album: albumName })
                            .then(() => {
                                this.$message.success({
                                    message: '已申请，等待批准',
                                    customClass: 'el-message--slim'
                                })
                                this.btns['taClaim'].name = '待批准'
                                this.btns['taClaim'].disabled = true
                            })
                    }
                }
                let by = this.album.opt.claiming_by
                if (Array.isArray(by) && by.includes(this.user)) {
                    taBtns['taClaim'].name = '待批准'
                    taBtns['taClaim'].disabled = true
                }
                if (this.album.ta && this.album.ta === this.user) {
                    taBtns['taClaim'].name = '已认领'
                    taBtns['taClaim'].disabled = true
                }
            }
            return taBtns
        },
        rendeClaimBtns(owner, albumName) {
            let claimBtns = {}
            return claimBtns
        },
        rendeHomeBtns(name) {
            let homeBtns = {
                Share: {
                    name: this.shareText,
                    type: 'success',
                    click: this.shareTo
                },
                Status: {
                    name: this.statusText,
                    type: 'warning',
                    click: () => {
                        this.$http
                            .put(`/user/album/${name}?action=status`)
                            .then(() => {
                                this.$message.success({
                                    message: `专辑已 <b>${this.statusText}</b>`,
                                    dangerouslyUseHTMLString: true,
                                    customClass: 'el-message--slim'
                                })
                                this.album.public = !this.album.public
                                this.btns['Status'].name = this.statusText
                            })
                    }
                },
                Delete: {
                    name: '删除',
                    type: 'danger',
                    icon: 'el-icon-delete',
                    click: () => {
                        this.$confirm(
                            '此操作将永久删除该专辑（及照片）, 是否继续?',
                            '温馨提示',
                            {
                                type: 'error',
                                confirmButtonText: '确定',
                                cancelButtonText: '取消',
                                customClass: 'el-message-box--slim'
                            }
                        ).then(() => {
                            this.$http
                                .delete(`/user/album/${name}`)
                                .then(() => {
                                    this.$message.success({
                                        message: '已删除专辑',
                                        customClass: 'el-message--slim'
                                    })
                                    this.$router.push({ path: '/' })
                                })
                        })
                    }
                }
            }
            let by = this.album.opt.claiming_by
            if (Array.isArray(by) && by.length > 0) {
                homeBtns['Share']['badge'] = {
                    dot: true,
                    type: 'success'
                }
            }
            return homeBtns
        },
        shareTo() {
            let msg = '请输入分享给Ta的用户名（覆盖已有分享）'
            let by = this.album.opt.claiming_by
            if (Array.isArray(by) && by.length > 0) {
                let plus = '<br>此专辑有认领者：' + by.join(', ')
                msg += plus
            }
            this.$prompt(msg, '温馨提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                dangerouslyUseHTMLString: true,
                customClass: 'el-message-box--slim',
                inputPlaceholder: '分享给Ta（或在认领者中选择）'
            }).then(({ value }) => {
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
                let name = this.album.name
                this.$http
                    .put(`/user/album/${name}?action=share`, { ta })
                    .then(() => {
                        this.$message.success({
                            message: `专辑已分享给${ta}`,
                            customClass: 'el-message--slim'
                        })
                        this.album.ta = ta
                        this.album.opt.claiming_by = []
                        this.btns['Share'].name = this.shareText
                        this.btns['Share'].badge = null
                    })
            })
        }
    }
}
</script>
