<template>
    <el-row type="flex" justify="center" align="middle">
        <el-col :md="12" :sm="16" :xs="20">
            <div v-if="isLogin">
                <div class="description">
                    <i
                        class="saintic-icon saintic-icon-goddess saintic-icon-3"
                    ></i>
                    <br />
                    {{ slogan }}
                </div>
                <el-form
                    ref="fairy"
                    :model="af"
                    :rules="rules"
                    size="small"
                    label-width="70px"
                >
                    <el-row>
                        <el-form-item label="专辑" prop="album">
                            <el-select
                                v-model="af.album"
                                placeholder="请选择或新建专辑"
                                allow-create
                                filterable
                            >
                                <el-option
                                    v-for="item in albums"
                                    :key="item"
                                    :value="item"
                                ></el-option>
                            </el-select>
                        </el-form-item>
                    </el-row>
                    <el-row :gutter="1">
                        <el-col :sm="17" :xs="17">
                            <el-form-item label="照片" prop="src">
                                <el-input
                                    v-model="af.src"
                                    placeholder="请输入照片（视频）地址或上传"
                                    show-word-limit
                                    clearable
                                    prefix-icon="el-icon-link"
                                ></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :sm="6" :xs="6">
                            <el-form-item label-width="0" prop="">
                                <el-tooltip
                                    effect="dark"
                                    :content="upTip"
                                    placement="top"
                                >
                                    <el-upload
                                        :accept="acceptMimes"
                                        :before-upload="upBefore"
                                        :show-file-list="false"
                                        action="/api/user/upload"
                                        :on-success="upSuccess"
                                        :on-error="upErr"
                                        :headers="headers"
                                    >
                                        <el-button
                                            size="mini"
                                            type="primary"
                                            icon="el-icon-upload"
                                            >上传</el-button
                                        >
                                    </el-upload>
                                </el-tooltip>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :sm="24">
                            <el-form-item label="描述" prop="desc">
                                <el-input
                                    v-model="af.desc"
                                    placeholder="请输入照片（视频）描述"
                                    clearable
                                >
                                </el-input>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-form-item size="small">
                            <el-button type="primary" @click="submitForm"
                                >提交</el-button
                            >
                            <el-button @click="resetForm">重置</el-button>
                        </el-form-item>
                    </el-row>
                </el-form>
            </div>
            <div v-else><Welcome /></div>
        </el-col>
    </el-row>
</template>

<script>
import { mapState } from '@/libs/store.js'
import { IndexSlogan } from '@/libs/vars.js'
import Welcome from '@/components/Welcome.vue'

export default {
    name: 'Index',
    components: { Welcome },
    data() {
        return {
            upTip: '支持上传 jpg/jpeg/png/webp 类型图片（不超过10MB）',
            af: {
                album: '',
                desc: '',
                src: ''
            },
            rules: {
                album: [
                    {
                        required: true,
                        message: '请选择或新建专辑',
                        trigger: 'change'
                    }
                ]
            },
            albums: [],
            acceptMimes: 'image/*,video/*'
        }
    },
    computed: mapState({
        isLogin: 'isLogin',
        user: 'user',
        headers: (state) => ({ Authorization: 'Bearer ' + state.token }),
        slogan: (state) => {
            return state.slogan || IndexSlogan
        },
        upLimit: (state) => {
            return state.upload_limit || 10
        }
    }),
    methods: {
        submitForm() {
            this.$refs['fairy'].validate((valid) => {
                if (!valid) return
                this.$http.post('/user/fairy', this.af).then((res) => {
                    this.$message.success('已提交')
                    // try update albums
                    if (!this.albums.includes(this.af.album)) {
                        this.albums.push(this.af.album)
                    }
                    this.resetForm()
                })
            })
        },
        resetForm() {
            this.$refs['fairy'].resetFields()
        },
        upBefore(file) {
            let isRightSize = file.size / 1024 / 1024 < this.upLimit
            if (!isRightSize) {
                this.$message.error('文件大小超过限制')
                return
            }
            let isImage = new RegExp('image/*').test(file.type)
            let isVideo = new RegExp('video/*').test(file.type)
            if (!isImage && !isVideo) {
                this.$message.error('不支持上传的文件类型')
            }
            return isRightSize && (isImage || isVideo)
        },
        upSuccess(res) {
            if (res.success) {
                this.af.src = res.data.src
            } else {
                this.$message.error(res.message)
            }
        },
        upErr(err) {
            let msg = JSON.parse(err.message)
            this.$message.error(msg)
        }
    },
    created() {
        if (this.isLogin) {
            this.$http.get('/user/album').then((res) => {
                res.data.map((a) => {
                    this.albums.push(a.name)
                })
            })
        }
    }
}
</script>

<style scoped>
.description {
    margin-bottom: 10px;
}
</style>