<template>
  <AppLayout
    title="教务总览"
    description="把成绩、课表、考试安排和培养方案收拢到统一入口，在一套一致的中文软质感界面中快速完成查询。"
    :show-back="false"
  >
    <template #header-extra>
      <div class="grid grid-cols-2 gap-3 text-center">
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">登录状态</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ userStore.isLoggedIn ? '已连接' : '未连接' }}</p>
        </div>
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">功能模块</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ featureCards.length }}</p>
        </div>
      </div>
    </template>

    <section class="grid gap-6 xl:grid-cols-[minmax(0,1.2fr)_minmax(280px,0.8fr)]">
      <article class="surface-panel p-6 md:p-8">
        <span class="eyebrow">
          <AppIcon name="shield" class="h-4 w-4" />
          统一访问入口
        </span>

        <div class="mt-6 grid gap-6 lg:grid-cols-[minmax(0,1fr)_220px] lg:items-center">
          <div>
            <h2 class="font-display text-3xl font-extrabold tracking-tight text-ink md:text-5xl">
              让教务查询回到清晰、稳定、可快速到达的状态。
            </h2>
            <p class="mt-4 max-w-2xl text-sm leading-7 text-muted md:text-base">
              当前首页改造成统一控制台：顶部导航固定，模块卡片语义统一，表单与结果页都采用同一套 Neumorphism 设计语言，后续继续扩展新页面时不需要再重复造轮子。
            </p>

            <div class="mt-6 flex flex-wrap gap-3">
              <button type="button" class="surface-button-primary" @click="router.push('/grade')">
                <AppIcon name="grade" class="h-4 w-4" />
                先看成绩
              </button>
              <button type="button" class="surface-button" @click="router.push('/schedule')">
                <AppIcon name="schedule" class="h-4 w-4" />
                查看今日课表
              </button>
            </div>
          </div>

          <div class="surface-deep-well relative overflow-hidden p-6">
            <div class="hero-orbit absolute -right-4 -top-4 h-20 w-20 opacity-60"></div>
            <div class="hero-orbit absolute -bottom-5 left-2 h-14 w-14 opacity-40 [animation-delay:0.4s]"></div>
            <div class="relative grid gap-4">
              <div class="surface-well px-4 py-4">
                <p class="text-xs uppercase tracking-[0.2em] text-muted">模块覆盖</p>
                <p class="mt-2 font-display text-3xl font-bold text-ink">5 项</p>
                <p class="mt-2 text-sm leading-6 text-muted">成绩、课表、考试、选课结果、培养方案全部聚合在同一入口。</p>
              </div>
              <div class="surface-well px-4 py-4">
                <p class="text-xs uppercase tracking-[0.2em] text-muted">交互基线</p>
                <p class="mt-2 text-sm leading-6 text-muted">统一按钮深度、输入框凹陷态、卡片层级和中文信息密度。</p>
              </div>
            </div>
          </div>
        </div>
      </article>

      <aside class="surface-panel p-6">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="text-xs font-semibold uppercase tracking-[0.22em] text-muted">使用提示</p>
            <h3 class="mt-2 font-display text-2xl font-bold text-ink">当前状态</h3>
          </div>
          <button type="button" class="surface-icon-button" aria-label="退出登录" @click="handleLogout">
            <AppIcon name="logout" class="h-5 w-5" />
          </button>
        </div>

        <div class="mt-5 space-y-3">
          <div v-for="item in tips" :key="item.title" class="surface-well px-4 py-4">
            <div class="flex items-start gap-3">
              <div class="surface-pill flex h-10 w-10 items-center justify-center text-accent">
                <AppIcon :name="item.icon" class="h-4 w-4" />
              </div>
              <div>
                <p class="text-sm font-semibold text-ink">{{ item.title }}</p>
                <p class="mt-1 text-sm leading-6 text-muted">{{ item.desc }}</p>
              </div>
            </div>
          </div>
        </div>

        <div class="surface-warning mt-5">
          当前项目仍在持续完善中，若教务系统源站波动，页面可能出现空结果或接口错误提示。
        </div>
      </aside>
    </section>

    <section class="mt-6">
      <div class="mb-4 flex items-end justify-between gap-4">
        <div>
          <p class="text-sm font-semibold uppercase tracking-[0.18em] text-muted">Core Modules</p>
          <h2 class="font-display text-2xl font-bold tracking-tight text-ink md:text-3xl">快速进入教务模块</h2>
        </div>
        <p class="hidden max-w-xl text-sm leading-6 text-muted md:block">
          卡片、图标井、标题层级和说明文案都基于同一套 token，可继续沿用到后续新增页面。
        </p>
      </div>

      <div class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
        <button
          v-for="item in featureCards"
          :key="item.path"
          type="button"
          class="surface-card-interactive group p-6 text-left"
          @click="router.push(item.path)"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="surface-deep-well flex h-16 w-16 items-center justify-center text-accent transition-transform duration-300 group-hover:scale-105">
              <AppIcon :name="item.icon" class="h-7 w-7" />
            </div>
            <span class="surface-badge">{{ item.shortTitle }}</span>
          </div>

          <div class="mt-6">
            <h3 class="font-display text-2xl font-bold tracking-tight text-ink">{{ item.title }}</h3>
            <p class="mt-3 text-sm leading-7 text-muted">{{ item.description }}</p>
          </div>

          <div class="mt-6 flex items-center justify-between text-sm font-semibold text-accent">
            <span>进入模块</span>
            <span aria-hidden="true">→</span>
          </div>
        </button>
      </div>
    </section>
  </AppLayout>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import AppIcon from '@/components/AppIcon.vue'
import AppLayout from '@/components/AppLayout.vue'
import { primaryNavItems } from '@/constants/navigation'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const featureCards = computed(() => primaryNavItems.filter((item) => item.path !== '/home'))

const tips = computed(() => [
  {
    title: userStore.isLoggedIn ? '登录态已保存到本地浏览器' : '尚未检测到登录态',
    desc: userStore.isLoggedIn ? '当前浏览器中已缓存教务授权信息，可直接进入查询模块。' : '请先完成登录，成功后首页会显示已连接状态。',
    icon: 'user'
  },
  {
    title: '所有页面已统一为中文界面',
    desc: '首页、筛选区、结果卡片、空态和错误态都使用一致的中文信息架构。',
    icon: 'spark'
  },
  {
    title: '移动端和桌面端共用一套布局壳层',
    desc: '桌面端显示顶部导航，移动端切换为菜单抽屉，保证同一套视觉语言。',
    icon: 'schedule'
  }
])

function handleLogout() {
  userStore.clearCookie()
  router.replace('/login')
}
</script>
