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
            />
            <div
              class="w-32 h-10 border border-gray-300 rounded-md flex items-center justify-center bg-gray-50 cursor-pointer overflow-hidden"
              @click="refreshCaptcha"
              :title="captchaLoading ? '加载中...' : '点击刷新验证码'"
            >
              <img
                v-if="captchaImage"
                :src="'data:image/png;base64,' + captchaImage"
                alt="验证码"
                class="w-full h-auto"
              />
              <span v-else class="text-xs text-gray-400">
                {{ captchaLoading ? '加载中...' : '点击获取' }}
              </span>
            </div>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getInitCookie, login } from '@/api/zhjw'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: '',
  captcha: ''
})

const captchaImage = ref('')
const initCookie = ref('')
const loading = ref(false)
const captchaLoading = ref(false)
const errorMsg = ref('')

async function refreshCaptcha() {
  if (captchaLoading.value) return
  
  captchaLoading.value = true
  errorMsg.value = ''
  
  try {
    const res = await getInitCookie()
    console.log('验证码响应:', res)
    if (res.code === 0) {
      captchaImage.value = res.data.captcha_image
      initCookie.value = res.data.cookie
      console.log('验证码图片 base64 长度:', res.data.captcha_image?.length)
    } else {
      errorMsg.value = res.msg || '获取验证码失败'
    }
  } catch (err) {
    console.error('获取验证码失败:', err)
    errorMsg.value = '获取验证码失败，请检查网络'
  } finally {
    captchaLoading.value = false
  }
}

async function handleLogin() {
  if (!initCookie.value) {
    errorMsg.value = '请先获取验证码'
    return
  }

  loading.value = true
  errorMsg.value = ''

  try {
    const res = await login({
      username: form.value.username,
      password: form.value.password,
      captcha: form.value.captcha,
      init_cookie: initCookie.value
    })

    if (res.code === 0) {
      userStore.setCookie(res.data.cookie)
      alert('登录成功！')
    } else {
      errorMsg.value = res.msg || '登录失败'
      refreshCaptcha()
    }
  } catch (err) {
    errorMsg.value = '登录失败，请检查网络'
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
})
</script>
