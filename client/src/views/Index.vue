<template>
    <el-row type="flex" justify="center" align="middle">
        <el-col :md="12" :sm="16" :xs="20">
            <div v-if="isLogin">
                <div class="description">
                    <i
                        class="saintic-icon saintic-icon-goddess saintic-icon-3"
                    ></i>
                    <br />
                    她是小仙女啦
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
                                    v-for="(item, index) in albums"
                                    :key="index"
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
                                    placeholder="请输入照片地址或上传"
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
                                        accept="image/jpg, image/jpeg, image/webp, image/png"
                                        :before-upload="upBefore"
                                        :show-file-list="false"
                                        action="/api/user/upload"
                                        :on-success="upSuccess"
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
                                    placeholder="请输入照片描述"
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
            albums: []
        }
    },
    computed: mapState({
        isLogin: 'isLogin',
        user: 'user',
        headers: (state) => ({ Authorization: 'Bearer ' + state.token })
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
            let isRightSize = file.size / 1024 / 1024 < 10
            if (!isRightSize) {
                this.$message.error('文件大小超过 10MB')
            }
            let isAccept = new RegExp('image/*').test(file.type)
            if (!isAccept) {
                this.$message.error('应该选择image/*类型的文件')
            }
            return isRightSize && isAccept
        },
        upSuccess(res) {
            if (res.success) {
                this.af.src = res.data.src
            } else {
                this.$message.error(res.message)
            }
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