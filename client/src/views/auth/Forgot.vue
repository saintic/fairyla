<template>
    <div class="forgot-container">
        <el-form
            :model="forgotForm"
            :rules="forgotRule"
            status-icon
            ref="forgotForm"
            label-position="left"
            label-width="0px"
            class="forgot-area"
            @submit.native.prevent
        >
            <h3 class="title">忘记密码</h3>
            <el-form-item prop="username">
                <el-input
                    type="text"
                    v-model="forgotForm.username"
                    autocomplete="off"
                    placeholder="用户名"
                    @keyup.enter.native="handleSubmit"
                ></el-input>
            </el-form-item>
            <el-form-item>
                <router-link :to="loginRoute">
                    <el-button type="text" class="login">返回登录</el-button>
                </router-link>
            </el-form-item>
            <el-form-item style="width: 100%">
                <el-button
                    type="primary"
                    style="width: 100%"
                    @click="handleSubmit"
                    :loading="forgoting"
                    >提交</el-button
                >
            </el-form-item>
        </el-form>
    </div>
</template>

<script>
export default {
    name: 'Forgot',
    data() {
        return {
            forgoting: false,
            forgotForm: {
                username: ''
            },
            forgotRule: {
                username: [
                    {
                        required: true,
                        message: '请输入用户名',
                        trigger: 'blur'
                    }
                ]
            },
            loginRoute: {
                name: 'Login'
            }
        }
    },
    methods: {
        handleSubmit(event) {
            this.$refs.forgotForm.validate((valid) => {
                if (!valid) return
                this.forgoting = true
                this.$http
                    .post('/auth/forgot', this.forgotForm)
                    .then(() => {
                        this.forgoting = false
                        this.$message.success({
                            message: '已发送验证邮件，请查收',
                            customClass: 'el-message--slim'
                        })
                    })
                    .catch((e) => {
                        this.forgoting = false
                    })
            })
        }
    }
}
</script>

<style scoped>
.forgot-area {
    -webkit-border-radius: 5px;
    border-radius: 5px;
    margin: 70px auto 30px;
    width: 350px;
    padding: 35px 35px 15px;
    background: #fff;
    border: 1px solid #eaeaea;
    box-shadow: 0 0 25px #cac6c6;
}
.forgot-area .login {
    text-align: left;
    float: left;
}
</style>