<template>
  <div class="login-wrap">
    <div class="login-container">
      <!-- LOGO -->
      <div class="logo-box">
        <!-- <img class="logo" src="@/assets/svg/logo.svg" alt="Logo" /> -->
        <h1 class="logo-title">Evobot</h1>
      </div>

      <!-- 登录框 -->
      <div class="login-box">
        <el-form ref="formRef" :model="loginForm" :rules="formRules" class="login-form">
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              type="text"
              placeholder="请输入用户名"
            >
              <template #prefix>
                <el-icon :size="20"><user /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              show-password
              placeholder="请输入密码"
            >
              <template #prefix>
                <el-icon :size="20"><lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              class="login-button"
              :loading="loginLoading"
              @click="loginFun"
            >
              登 录
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 页脚 -->
      <footer class="login-footer">
        <p>
          <a href="https://fit2cloud.com/" target="_blank" class="link">
            FIT2CLOUD 飞致云
          </a>
        </p>
      </footer>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, toRefs } from 'vue'
import { onKeyStroke } from '@vueuse/core'
import { loginApi } from '@/api/login';
import { setToken } from '@/utils'
import router from '@/router'
import type { FormInstance, FormRules } from 'element-plus'

const loginLoading = ref(false)
const activeTab = ref('password')
const formRef = ref<FormInstance>()
const form = reactive({
  loginForm: {
    username: '',
    password: ''
  }
})
const { loginForm } = toRefs(form)
const formRules = reactive<FormRules>({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
})
// 触发登录
const loginFun = () => {
  // 表单校验
  formRef.value?.validate(async (valid) => {
    if (!valid) return;
    try {
      loginLoading.value = true
      const res = await loginApi(form.loginForm);
      const access_token = res.data?.access_token;
      if (access_token) {
        setToken(access_token)
      }
      loginLoading.value = false
      router.push('/home')
    } catch (error) {
      console.error('Login error:', error);
    } finally {
      loginLoading.value = false
    }
  })
}

// 回车事件
onMounted(() => {
  onKeyStroke('Enter', () => {
    if (loginLoading.value) return
    loginFun()
  })
})
</script>

<style lang="scss" scoped>
/* 页面背景 */
.login-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100vh;
  background-color: #f7f8fa; /* 淡灰背景 */
  font-family: 'Inter', Arial, sans-serif; /* 更现代化的字体 */
}

/* 容器 */
.login-container {
  width: 420px; 
  padding: 30px;
  background: #ffffff; 
  border-radius: 10px; 
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1); 
}

/* LOGO */
.logo-box {
  text-align: center;
  margin-bottom: 30px;

  .logo-title {
    margin-top: 10px;
    font-size: 24px; 
    font-weight: 600;
    color: #3e5c96; 
    font-family: 'Inter', Arial, sans-serif;
  }
}

/* 登录框 */
.login-box {
  .login-form {
    .el-input__inner {
      height: 48px; 
      font-size: 18px; 
      border-radius: 6px; 
    }

    .el-input__inner:focus {
      border-color: #606266; 
      box-shadow: 0 0 3px rgba(0, 0, 0, 0.1);
    }

    .el-form-item {
      margin-top: 40px; 
    }

    .login-button {
      width: 100%;
      height: 40px;
      font-size: 16px;
      font-weight: 500; 
      background-color: #3e5c96; 
      color: #ffffff;
      border-radius: 6px;
      font-family: 'Inter', Arial, sans-serif;
      transition: background-color 0.2s;
    }

    .login-button:hover {
      background-color: #19284d;
    }
  }
}

/* 页脚 */
.login-footer {
  text-align: center;
  margin-top: 15px; 
  font-size: 14px; 
  color: #aaaaaa;

  .link {
    color: #3e5c96;
    text-decoration: none;
    font-family: 'Inter', Arial, sans-serif;
  }
}
</style>
