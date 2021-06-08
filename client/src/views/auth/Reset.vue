<template>
    <div class="reset-container">
        <el-form
            :model="resetForm"
            :rules="resetRule"
            status-icon
            ref="resetForm"
            label-position="left"
            label-width="0px"
            class="reset-area"
            @submit.native.prevent
        >
            <h3 class="title">重置密码</h3>
            <el-form-item prop="password">
                <el-input
                    type="password"
                    v-model="resetForm.password"
                    autocomplete="off"
                    placeholder="新密码"
                    show-password
                    @keyup.enter.native="handleSubmit"
                ></el-input>
            </el-form-item>
            <el-form-item style="width: 100%">
                <el-button
                    type="primary"
                    style="width: 100%"
                    @click="handleSubmit"
                    :loading="reseting"
                    >提交</el-button
                >
            </el-form-item>
        </el-form>
    </div>
</template>

<script>
export default {
    name: 'Reset',
    data() {
        return {
            reseting: false,
            resetForm: {
                jwt: this.$route.query.jwt,
                password: ''
            },
            resetRule: {
                password: [
                    {
                        required: true,
                        message: '请输入新密码',
                        trigger: 'blur'
                    }
                ]
            }
        }
    },
    methods: {
        handleSubmit(event) {
            this.$refs.resetForm.validate((valid) => {
                if (!valid) return
                this.reseting = true
                this.$http
                    .post('/auth/reset_passwd', this.resetForm)
                    .then(() => {
                        this.reseting = false
                        this.$message.success({
                            message: '已重置密码，请重新登录',
                            customClass: 'el-message--slim'
                        })
                        this.$router.push({ name: 'Login' })
                    })
                    .catch((e) => {
                        this.reseting = false
                    })
            })
        }
    }
}
</script>

<style scoped>
.reset-area {
    -webkit-border-radius: 5px;
    border-radius: 5px;
    margin: 70px auto 30px;
    width: 350px;
    padding: 35px 35px 15px;
    background: #fff;
    border: 1px solid #eaeaea;
    box-shadow: 0 0 25px #cac6c6;
}
</style>