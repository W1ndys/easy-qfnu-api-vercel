<template>
  <AppLayout
    title="入口仪表盘"
    description="登录成功后可从这里快速进入各个教务功能页面。"
    :show-back="false"
  >
    <template #header-extra>
      <div class="grid grid-cols-2 gap-2.5 text-center">
        <div class="surface-well px-3 py-3.5">
          <p class="text-xs text-muted">登录状态</p>
          <p class="mt-2 text-sm font-semibold text-ink">{{ userStore.isLoggedIn ? '已登录' : '未登录' }}</p>
        </div>
        <div class="surface-well px-3 py-3.5">
          <p class="text-xs text-muted">功能入口</p>
          <p class="mt-2 text-sm font-semibold text-ink">{{ featureCards.length }} 个</p>
        </div>
      </div>
    </template>

    <section class="grid gap-3.5">
      <article class="surface-panel p-4 sm:p-5">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="text-xs font-semibold uppercase tracking-[0.2em] text-muted">Dashboard</p>
            <h2 class="mt-2 font-display text-xl font-bold tracking-tight text-ink sm:text-2xl">快捷入口</h2>
          </div>
          <button type="button" class="surface-icon-button" aria-label="退出登录" @click="handleLogout">
            <AppIcon name="logout" class="h-5 w-5" />
          </button>
        </div>

        <p class="mt-3 text-sm leading-6 text-muted">
          选择下方模块即可进入对应查询页面。顶部导航和这里的卡片入口保持一致。
        </p>
      </article>

      <section class="grid gap-3.5 md:grid-cols-2 xl:grid-cols-3">
        <button
          v-for="item in featureCards"
          :key="item.path"
          type="button"
          class="surface-card-interactive group p-4 text-left sm:p-5"
          @click="router.push(item.path)"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="surface-deep-well flex h-12 w-12 items-center justify-center text-accent transition-transform duration-300 group-hover:scale-105 sm:h-14 sm:w-14">
              <AppIcon :name="item.icon" class="h-5 w-5 sm:h-6 sm:w-6" />
            </div>
            <span class="surface-badge">{{ item.shortTitle }}</span>
          </div>

          <div class="mt-4">
            <h3 class="font-display text-lg font-bold tracking-tight text-ink sm:text-xl">{{ item.title }}</h3>
            <p class="mt-2.5 text-sm leading-6 text-muted">{{ item.description }}</p>
          </div>

          <div class="mt-4 flex items-center justify-between text-sm font-semibold text-accent">
            <span>立即进入</span>
            <span aria-hidden="true">→</span>
          </div>
        </button>
      </section>
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

function handleLogout() {
  userStore.clearCookie()
  router.replace('/login')
}
</script>
