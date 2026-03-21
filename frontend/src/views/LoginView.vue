<template>
  <div class="relative min-h-screen overflow-hidden bg-surface px-4 py-8 sm:px-6 lg:px-8">
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="hero-orbit absolute left-[-3rem] top-16 h-40 w-40 opacity-60"></div>
      <div class="hero-orbit absolute right-[-4rem] top-24 h-64 w-64 opacity-70"></div>
      <div class="hero-orbit absolute bottom-12 left-1/3 h-24 w-24 opacity-50 [animation-delay:0.5s]"></div>
    </div>

    <div class="page-shell relative flex min-h-[calc(100vh-4rem)] items-center px-0 sm:px-0 lg:px-0">
      <div class="grid w-full gap-6 lg:grid-cols-[minmax(0,1.1fr)_minmax(360px,460px)] lg:items-center">
        <section class="surface-panel hidden p-8 lg:block xl:p-10">
          <span class="eyebrow">
            <AppIcon name="shield" class="h-4 w-4" />
            曲园教务入口
          </span>

          <h1 class="mt-6 font-display text-5xl font-extrabold tracking-tight text-ink">
            中文、统一、柔和层次的教务查询体验。
          </h1>
          <p class="mt-5 max-w-2xl text-base leading-8 text-muted">
            登录后即可进入新版控制台。所有功能模块都采用同一套软质感设计 token，表单、卡片、统计块和导航逻辑已经统一，后续扩展不会再出现风格割裂。
          </p>

          <div class="mt-8 grid gap-4 sm:grid-cols-2">
            <div v-for="item in highlights" :key="item.title" class="surface-deep-well p-5">
              <div class="surface-pill flex h-12 w-12 items-center justify-center text-accent">
                <AppIcon :name="item.icon" class="h-5 w-5" />
              </div>
              <h2 class="mt-5 text-lg font-semibold text-ink">{{ item.title }}</h2>
              <p class="mt-2 text-sm leading-7 text-muted">{{ item.desc }}</p>
            </div>
          </div>
        </section>

        <section class="surface-panel p-6 sm:p-8">
          <div class="mx-auto max-w-md">
            <div class="surface-deep-well p-6 text-center">
              <div class="surface-pill mx-auto flex h-14 w-14 items-center justify-center text-accent">
                <AppIcon name="lock" class="h-6 w-6" />
              </div>
              <h1 class="mt-4 font-display text-3xl font-extrabold tracking-tight text-ink">登录教务系统</h1>
              <p class="mt-3 text-sm leading-7 text-muted">
                输入学号与密码后进入新版首页。授权信息仅保存在当前浏览器中。
              </p>
            </div>

            <form class="mt-6 space-y-4" @submit.prevent="handleLogin">
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
                {{ loading ? '正在登录，请稍候…' : '进入教务控制台' }}
              </button>
            </form>

            <div class="surface-warning mt-5">
              若长时间登录失败，请先确认教务系统源站状态、账号密码是否正确，以及本地网络是否能够访问学校服务。
            </div>
          </div>
        </section>
      </div>
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

const highlights = [
  {
    title: '统一视觉层级',
    desc: '登录页、首页和业务页都使用同一套圆角、阴影、字体与交互状态。',
    icon: 'spark'
  },
  {
    title: '查询入口清晰',
    desc: '成绩、课表、考试和培养方案会收拢在顶部导航与首页卡片中。',
    icon: 'grade'
  },
  {
    title: '移动端可直接使用',
    desc: '移动端提供菜单抽屉和 44px 以上点击区域，避免误触。',
    icon: 'schedule'
  },
  {
    title: '信息展示更集中',
    desc: '统计信息、空态和错误态都统一表达，不再散落在不同样式里。',
    icon: 'book'
  }
]

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
