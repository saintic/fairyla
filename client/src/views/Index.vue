<template>
    <el-row type="flex" justify="center" align="middle">
        <el-col :span="12">
            <div v-if="isLogin">
                <div class="description">
                    <i
                        class="saintic-icon saintic-icon-goddess saintic-icon-2-5"
                    ></i>
                    <br />
                    她是小仙女啦
                </div>

                <el-form
                    ref="fairy"
                    :model="af"
                    :rules="rules"
                    size="small"
                    label-width="100px"
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
                                    :label="item"
                                    :value="item"
                                ></el-option>
                            </el-select>
                        </el-form-item>
                    </el-row>
                    <el-row>
                        <el-col :span="17">
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
                        <el-col :span="7">
                            <el-form-item label-width="0" prop="">
                                <el-tooltip
                                    effect="dark"
                                    :content="upTip"
                                    placement="top"
                                >
                                    <el-upload
                                        accept="image/jpg, image/jpeg, image/webp, image/png"
                                        :before-upload="upBefore"
                                        :http-request="upload"
                                        :show-file-list="false"
                                        action=""
                                    >
                                        <el-button
                                            size="small"
                                            type="primary"
                                            icon="el-icon-upload"
                                            >点击上传</el-button
                                        >
                                    </el-upload>
                                </el-tooltip>
                            </el-form-item>
                        </el-col>
                    </el-row>
                    <el-row>
                        <el-col :span="24">
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
            <div v-else>你好，INDEX</div>
        </el-col>
    </el-row>
</template>

<script>
import axios from 'axios'
import { mapState } from '@/libs/store.js'

export default {
    name: 'Index',
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
        api: (state) => {
            return state.sapic.api
        },
        field: (state) => {
            return state.sapic.field || 'picbed'
        },
        linkToken: (state) => {
            return state.sapic.token
        }
    }),
    methods: {
        submitForm() {
            this.$refs['fairy'].validate((valid) => {
                if (!valid) return
                // TODO 提交表单
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
        upload(opt) {
            console.log(opt)
            let data = new FormData()
            data.append(this.field, opt.file)
            data.append('album', this.af.album)
            data.append('title', this.af.desc)
            axios
                .post(this.api, data, {
                    headers: {
                        'Content-Type': 'multipart/form-data',
                        Authorization: 'LinkToken ' + this.linkToken
                    }
                })
                .then((res) => {
                    console.log('上传图片接口-数据', res)
                    let data = res.data
                    if (data.code === 0) {
                        this.af.src = data.src
                    } else {
                        this.$message.warning(data.msg)
                    }
                })
                .catch((err) => {
                    this.$message.error(err)
                })
        }
    },
    created() {
        if (this.isLogin) {
            this.$http.get('/user/album').then((res) => {
                console.log(res)
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