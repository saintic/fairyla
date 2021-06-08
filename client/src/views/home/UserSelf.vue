<template>
    <el-row type="flex" justify="center" align="middle">
        <el-col :md="12" :sm="16" :xs="20">
            <el-tabs tab-position="left">
                <el-tab-pane label="个人资料">
                    <el-form
                        ref="userProfile"
                        :model="userProfile"
                        label-width="60px"
                        size="mini"
                    >
                        <el-form-item label="昵称" prop="alias">
                            <el-input
                                v-model="userProfile.alias"
                                placeholder="用户别名、昵称"
                                clearable
                            ></el-input>
                        </el-form-item>
                        <el-form-item label="邮箱" prop="email">
                            <el-input
                                v-model="userProfile.email"
                                placeholder="个人邮箱（找回密码）"
                                clearable
                            ></el-input>
                        </el-form-item>
                        <el-form-item label="简介" prop="bio">
                            <el-input
                                type="textarea"
                                v-model="userProfile.bio"
                                placeholder="个人介绍、交友宣言"
                                clearable
                                :rows="3"
                            ></el-input>
                        </el-form-item>
                        <el-form-item>
                            <el-button type="primary" @click="submitProfileForm"
                                >保存</el-button
                            >
                            <el-button @click="resetProfileForm"
                                >重置</el-button
                            >
                        </el-form-item>
                    </el-form>
                </el-tab-pane>
                <el-tab-pane label="账号安全">
                    <el-row type="flex">
                        <el-form
                            ref="userPasswd"
                            :model="userPasswd"
                            :rules="passwdRules"
                            label-width="80px"
                            size="mini"
                        >
                            <el-form-item label="当前密码" prop="old_passwd">
                                <el-input
                                    v-model="userPasswd.old_passwd"
                                    placeholder="请输入当前密码"
                                    clearable
                                    show-password
                                ></el-input>
                            </el-form-item>
                            <el-form-item label="新密码" prop="new_passwd">
                                <el-input
                                    v-model="userPasswd.new_passwd"
                                    placeholder="请输入新密码"
                                    clearable
                                    show-password
                                ></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button
                                    type="primary"
                                    @click="submitPasswdForm"
                                    >修改密码</el-button
                                >
                                <el-button @click="resetPasswdForm"
                                    >重置</el-button
                                >
                            </el-form-item>
                        </el-form>
                    </el-row>
                </el-tab-pane>
                <el-tab-pane label="应用设置">
                    <el-row type="flex">
                        <el-form
                            ref="appSetting"
                            :model="appSetting"
                            label-width="auto"
                            size="mini"
                        >
                            <el-form-item label="专辑默认">
                                <el-switch
                                    v-model="appSetting.album_default_public"
                                    active-text="公开"
                                    inactive-text="私有"
                                ></el-switch>
                            </el-form-item>
                            <el-form-item label="Slogan" prop="slogan">
                                <el-input
                                    v-model="appSetting.slogan"
                                    placeholder="首页上传标语"
                                    clearable
                                ></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button
                                    type="primary"
                                    @click="submitSettingForm"
                                    >保存</el-button
                                >
                                <el-button @click="resetSettingForm"
                                    >重置</el-button
                                >
                            </el-form-item>
                        </el-form>
                    </el-row>
                </el-tab-pane>
            </el-tabs>
        </el-col>
    </el-row>
</template>

<script>
import { mutations } from '@/libs/store.js'

export default {
    name: 'UserSelf',
    data() {
        let ui = this.$store.state.userinfo || {}
        return {
            userProfile: {
                alias: ui.alias,
                email: ui.email,
                bio: ui.bio
            },
            userPasswd: {
                new_passwd: '',
                old_passwd: ''
            },
            passwdRules: {
                new_passwd: [
                    {
                        required: true,
                        message: '请输入新密码',
                        trigger: 'blur'
                    },
                    {
                        pattern: /^[\S]{6,32}$/,
                        message: '密码要求6-32位且无空格',
                        trigger: 'blur'
                    }
                ],
                old_passwd: [
                    {
                        required: true,
                        message: '请输入当前密码',
                        trigger: 'blur'
                    }
                ]
            },
            appSetting: {
                album_default_public: ui.album_default_public,
                slogan: ui.slogan
            }
        }
    },
    methods: {
        submitProfileForm() {
            this.$refs['userProfile'].validate((valid) => {
                if (!valid) return
                this.$http
                    .put('/user/profile', this.userProfile)
                    .then((res) => {
                        this.$message.success({
                            message: '资料已更新',
                            customClass: 'el-message--slim'
                        })
                        for (let [k, v] of Object.entries(res.data)) {
                            mutations.commitNested('userinfo', k, v)
                        }
                        this.userProfile = res.data
                    })
            })
        },
        resetProfileForm() {
            this.$refs['userProfile'].resetFields()
        },
        submitPasswdForm() {
            this.$refs['userPasswd'].validate((valid) => {
                if (!valid) return
                this.$http.put('/user/passwd', this.userPasswd).then(() => {
                    this.$message.success({
                        message: '密码已更新，请重新登录',
                        customClass: 'el-message--slim'
                    })
                    this.$router.push({ name: 'Logout' })
                    this.$router.push({ name: 'Login' })
                })
            })
        },
        resetPasswdForm() {
            this.$refs['userPasswd'].resetFields()
        },
        submitSettingForm() {
            this.$refs['appSetting'].validate((valid) => {
                if (!valid) return
                this.$http.put('/user/setting', this.appSetting).then((res) => {
                    this.$message.success({
                        message: '设置已更新',
                        customClass: 'el-message--slim'
                    })
                    for (let [k, v] of Object.entries(res.data)) {
                        mutations.commitNested('userinfo', k, v)
                    }
                    this.appSetting = res.data
                })
            })
        },
        resetSettingForm() {
            this.$refs['appSetting'].resetFields()
        }
    }
}
</script>
