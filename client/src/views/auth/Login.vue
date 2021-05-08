<template>
    <div class="login-container">
        <el-form
            :model="loginForm"
            :rules="loginRule"
            status-icon
            ref="loginForm"
            label-position="left"
            label-width="0px"
            class="login-area"
        >
            <h3 class="title">登录</h3>
            <el-form-item prop="username">
                <el-input
                    type="text"
                    v-model="loginForm.username"
                    autocomplete="off"
                    placeholder="用户名"
                ></el-input>
            </el-form-item>
            <el-form-item prop="password">
                <el-input
                    type="password"
                    v-model="loginForm.password"
                    autocomplete="off"
                    placeholder="密码"
                    show-password
                    @keyup.enter.native="handleSubmit"
                ></el-input>
            </el-form-item>
            <el-form-item>
                <el-checkbox v-model="checked" class="rememberme"
                    >记住我</el-checkbox
                >
                <el-link type="info" class="forgot">忘记密码？</el-link>
            </el-form-item>
            <el-form-item style="width: 100%">
                <el-button
                    type="primary"
                    style="width: 100%"
                    @click="handleSubmit"
                    :loading="logining"
                    >登录</el-button
                >
            </el-form-item>
        </el-form>
    </div>
</template>

<script>
export default {
    name: 'Login',
    data() {
        return {
            logining: false,
            loginForm: {
                username: '',
                password: ''
            },
            loginRule: {
                username: [
                    {
                        required: true,
                        message: '请输入用户名',
                        trigger: 'blur'
                    }
                ],
                password: [
                    {
                        required: true,
                        message: '请输入密码',
                        trigger: 'blur'
                    }
                ]
            },
            checked: false
        }
    },
    methods: {
        handleSubmit(event) {
            this.$refs.loginForm.validate((valid) => {
                if (valid) {
                    this.logining = true
                    this.$http
                        .post('/auth/signin', {
                            username: this.loginForm.username,
                            password: this.loginForm.password,
                            remember: this.checked
                        })
                        .then((res) => {
                            this.logining = false
                            this.$message.success('登录成功')
                            this.$store.mutations.setLogin(
                                this.loginForm.username,
                                res.data.token
                            )
                            this.$store.actions.saveConfig2Local()
                            this.$store.actions.fetchConfig()
                            let rdt = this.$route.query.redirect || ''
                            if (!rdt.startsWith('/')) {
                                rdt = '/'
                            }
                            this.$router.push({ path: rdt })
                        })
                        .catch((e) => {
                            this.logining = false
                        })
                } else {
                    return false
                }
            })
        }
    }
}
</script>

<style scoped>
.login-area {
    -webkit-border-radius: 5px;
    border-radius: 5px;
    margin: 70px auto 30px;
    width: 350px;
    padding: 35px 35px 15px;
    background: #fff;
    border: 1px solid #eaeaea;
    box-shadow: 0 0 25px #cac6c6;
}
.login-area .forgot {
    float: right;
}
.login-area .rememberme {
    text-align: left;
    float: left;
}
</style>