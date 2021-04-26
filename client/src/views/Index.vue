<template>
    <el-row type="flex" class="row-bg" justify="center">
        <el-col :span="6">
            <div v-if="isLogin">
                <div class="description">
                    <i
                        class="saintic-icon saintic-icon-goddess saintic-icon-2-5"
                    ></i>
                    <br />
                    她是小仙女啦
                </div>

                <el-row :gutter="15">
                    <el-form
                        ref="fairy"
                        :model="albumfairy"
                        :rules="rules"
                        size="medium"
                        label-width="100px"
                    >
                        <el-col :span="24">
                            <el-form-item label="专辑" prop="album">
                                <el-select
                                    v-model="albumfairy.album"
                                    placeholder="请选择或新建专辑"
                                    allow-create
                                    filterable
                                >
                                    <el-option
                                        v-for="(item, index) in albumOptions"
                                        :key="index"
                                        :label="item.label"
                                        :value="item.value"
                                        :disabled="item.disabled"
                                    ></el-option>
                                </el-select>
                            </el-form-item>
                        </el-col>
                        <el-col :span="12">
                            <el-form-item label="描述" prop="desc">
                                <el-input
                                    v-model="albumfairy.desc"
                                    placeholder="请输入照片描述"
                                    clearable
                                >
                                </el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="7">
                            <el-form-item label-width="0" prop="">
                                <el-upload
                                    ref=""
                                    :file-list="fileList"
                                    :action="Action"
                                    :before-upload="BeforeUpload"
                                    accept="image/*"
                                    name="picbed"
                                >
                                    <el-button
                                        size="small"
                                        type="primary"
                                        icon="el-icon-upload"
                                        >点击上传</el-button
                                    >
                                </el-upload>
                            </el-form-item>
                        </el-col>
                        <el-col :span="24">
                            <el-form-item label="照片" prop="src">
                                <el-input
                                    v-model="albumfairy.src"
                                    placeholder="请输入照片地址或上传"
                                    show-word-limit
                                    clearable
                                    prefix-icon="el-icon-link"
                                ></el-input>
                            </el-form-item>
                        </el-col>
                        <el-col :span="24">
                            <el-form-item size="large">
                                <el-button type="primary" @click="submitForm"
                                    >提交</el-button
                                >
                                <el-button @click="resetForm">重置</el-button>
                            </el-form-item>
                        </el-col>
                    </el-form>
                </el-row>

                <div class="upload-area">
                    <el-upload
                        class="upload"
                        drag
                        :action="api"
                        :headers="headers"
                        :data="data"
                        :name="field"
                        accept="images"
                        :limit="10"
                        multiple
                    >
                        <i class="el-icon-upload"></i>
                        <div class="el-upload__text">
                            将文件拖到此处，或<em>点击上传</em>
                        </div>
                        <template #tip>
                            <div class="el-upload__tip">
                                可以上传 jpg/jpeg/png/webp 文件，且不超过 10MB
                                <br />
                                也可以在此处粘贴图片地址上传（不支持粘贴图片本身哦）
                            </div>
                        </template>
                    </el-upload>
                </div>
            </div>
            <div v-else>你好，INDEX</div>
        </el-col>
    </el-row>
</template>

<script>
import { mapState } from '@/libs/store.js'

export default {
    name: 'Index',
    data() {
        return {
            data: { picbed: this.field, album: '' },
            headers: { Authorization: 'LinkToken ' + this.linkToken },
            albumfairy: {
                album: 2,
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
                ],
                desc: [
                    {
                        required: true,
                        message: '请输入照片描述',
                        trigger: 'blur'
                    }
                ],
                src: []
            },
            Action: 'https://jsonplaceholder.typicode.com/posts/',
            fileList: [],
            albumOptions: [
                {
                    label: '选项一',
                    value: 1
                },
                {
                    label: '选项二',
                    value: 2
                }
            ]
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
        BeforeUpload(file) {
            let isRightSize = file.size / 1024 / 1024 < 10
            if (!isRightSize) {
                this.$message.error('文件大小超过 10MB')
            }
            let isAccept = new RegExp('image/*').test(file.type)
            if (!isAccept) {
                this.$message.error('应该选择image/*类型的文件')
            }
            return isRightSize && isAccept
        }
    }
}
</script>

<style scoped>
.description {
    margin-bottom: 10px;
}
</style>