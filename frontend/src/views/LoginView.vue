<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center p-4">
    <div class="w-full max-w-md bg-white rounded-lg shadow-md p-6">
      <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">
        QFNU 教务系统登录
      </h1>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">学号</label>
          <input
            v-model="form.username"
            type="text"
            placeholder="请输入学号"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">密码</label>
          <input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            required
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">验证码</label>
          <div class="flex gap-2">
            <input
              v-model="form.captcha"
              type="text"
              placeholder="请输入验证码"
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              required
              autocomplete="off"
            />
            <img
              :src="captchaUrl"
              alt="验证码"
              class="w-32 h-10 border border-gray-300 rounded-md cursor-pointer object-cover bg-gray-50"
              @click="refreshCaptcha"
              :title="captchaLoading ? '加载中...' : '点击刷新验证码'"
            />
          </div>
        </div>

        <div v-if="errorMsg" class="text-red-500 text-sm text-center">
          {{ errorMsg }}
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2 px-4 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { getCaptchaUrl, initSession, login } from '@/api/zhjw'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const form = ref({
  username: '',
  password: '',
  captcha: ''
})

const captchaKey = ref(Date.now())
const loading = ref(false)
const errorMsg = ref('')

const captchaUrl = computed(() => {
  return `${getCaptchaUrl()}?t=${captchaKey.value}`
})

async function refreshCaptcha() {
  captchaKey.value = Date.now()
  form.value.captcha = ''
  errorMsg.value = ''
}

async function handleLogin() {
  if (!form.value.captcha) {
    errorMsg.value = '请输入验证码'
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    await login({
      username: form.value.username,
      password: form.value.password,
      captcha: form.value.captcha
    })
    
    userStore.isLoggedIn.value = true
    alert('登录成功！')
  } catch (err) {
    console.error('登录失败:', err)
    errorMsg.value = err.message || '登录失败，请检查网络'
    form.value.captcha = ''
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    await initSession()
    refreshCaptcha()
  } catch (err) {
    console.error('初始化会话失败:', err)
    errorMsg.value = '连接服务器失败，请检查网络'
  }
})
</script>
