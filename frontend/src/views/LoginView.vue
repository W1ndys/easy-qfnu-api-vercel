<template>
  <div class="relative min-h-screen overflow-hidden bg-surface px-3 py-4 sm:px-4 sm:py-6">
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="hero-orbit absolute left-[-3rem] top-12 h-24 w-24 opacity-50 sm:h-32 sm:w-32 sm:opacity-60"></div>
      <div class="hero-orbit absolute right-[-3rem] top-20 h-36 w-36 opacity-60 sm:h-48 sm:w-48 sm:opacity-68"></div>
      <div class="hero-orbit absolute bottom-8 left-1/3 h-16 w-16 opacity-40 [animation-delay:0.5s] sm:h-20 sm:w-20"></div>
    </div>

    <div class="page-shell relative flex min-h-[calc(100vh-2rem)] items-center justify-center px-0">
      <section class="surface-panel w-full p-4 sm:max-w-md sm:p-5">
        <div class="surface-deep-well p-4 text-center sm:p-5">
          <div class="surface-pill mx-auto flex h-12 w-12 items-center justify-center text-accent">
            <AppIcon name="lock" class="h-5 w-5" />
          </div>
          <p class="mt-4 text-[11px] font-semibold uppercase tracking-[0.2em] text-muted">Cookie Assistant</p>
          <h1 class="mt-2 font-display text-2xl font-extrabold tracking-tight text-ink">曲奇助手</h1>
          <p class="mt-3 text-sm leading-6 text-muted">
            输入学号和密码后进入功能入口仪表盘。
          </p>
        </div>

        <form class="mt-4 space-y-3.5" @submit.prevent="handleLogin">
          <label>
            <span class="surface-field-label">学号</span>
            <input
              v-model.trim="form.username"
              type="text"
              placeholder="请输入学号"
              class="surface-input"
              autocomplete="username"
              required
            />
          </label>

          <label>
            <span class="surface-field-label">密码</span>
            <input
              v-model="form.password"
              type="password"
              placeholder="请输入教务密码"
              class="surface-input"
              autocomplete="current-password"
              required
            />
          </label>

          <div v-if="errorMsg" class="surface-error">
            {{ errorMsg }}
          </div>

          <button type="submit" :disabled="loading" class="surface-button-primary w-full">
            <AppIcon :name="loading ? 'spark' : 'shield'" class="h-4 w-4" />
            {{ loading ? '正在登录，请稍候…' : '登录并进入仪表盘' }}
          </button>
        </form>

        <div class="surface-warning mt-4">
          登录成功后会先进入入口仪表盘，再选择具体查询模块。
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { login } from '@/api/zhjw'
import AppIcon from '@/components/AppIcon.vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const router = useRouter()

const form = ref({
  username: '',
  password: ''
})

const loading = ref(false)
const errorMsg = ref('')

async function handleLogin() {
  loading.value = true
  errorMsg.value = ''

  try {
    const res = await login({
      username: form.value.username,
      password: form.value.password
    })

    if (res.code === 200 && res.data?.cookie) {
      userStore.setCookie(res.data.cookie)
      router.push('/home')
    } else {
      errorMsg.value = res.msg || '登录失败，请检查账号或密码。'
    }
  } catch (err) {
    console.error('登录失败:', err)
    errorMsg.value = err.response?.data?.msg || err.message || '登录失败，请检查网络连接。'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (userStore.isLoggedIn) {
    router.replace('/home')
  }
})
</script>
