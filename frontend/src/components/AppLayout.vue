<template>
  <div class="relative min-h-screen overflow-hidden bg-surface text-ink">
    <div class="pointer-events-none absolute inset-x-0 top-0 h-[34rem] overflow-hidden">
      <div class="hero-orbit absolute -left-16 top-24 h-40 w-40 opacity-60"></div>
      <div class="hero-orbit absolute right-[-3rem] top-10 h-56 w-56 opacity-70"></div>
      <div class="hero-orbit absolute bottom-10 left-1/3 h-24 w-24 opacity-50 [animation-delay:0.6s]"></div>
    </div>

    <header class="sticky top-0 z-40 px-4 pt-4 sm:px-6 lg:px-8">
      <div class="page-shell px-0 sm:px-0 lg:px-0">
        <div class="surface-panel px-4 py-3 supports-[backdrop-filter]:bg-surface/80 supports-[backdrop-filter]:backdrop-blur md:px-6">
          <div class="flex items-center gap-3">
            <button
              v-if="showBack"
              type="button"
              class="surface-icon-button md:hidden"
              aria-label="返回上一页"
              @click="router.back()"
            >
              <AppIcon name="back" class="h-5 w-5" />
            </button>

            <button
              v-else
              type="button"
              class="surface-icon-button md:hidden"
              aria-label="回到首页"
              @click="router.push('/home')"
            >
              <AppIcon name="home" class="h-5 w-5" />
            </button>

            <button type="button" class="min-w-0 flex-1 text-left" @click="router.push('/home')">
              <p class="text-[11px] font-semibold uppercase tracking-[0.24em] text-muted">QFNU Academic Hub</p>
              <p class="font-display text-lg font-extrabold tracking-tight text-ink md:text-xl">曲园教务控制台</p>
            </button>

            <nav class="hidden flex-1 items-center justify-center gap-2 xl:flex">
              <RouterLink
                v-for="item in primaryNavItems"
                :key="item.path"
                :to="item.path"
                class="surface-nav-link"
                :class="{ 'surface-nav-link-active': isActive(item.path) }"
              >
                <AppIcon :name="item.icon" class="h-4 w-4" />
                <span>{{ item.shortTitle }}</span>
              </RouterLink>
            </nav>

            <div class="hidden items-center gap-2 md:flex xl:hidden">
              <button
                v-if="showBack"
                type="button"
                class="surface-icon-button"
                aria-label="返回上一页"
                @click="router.back()"
              >
                <AppIcon name="back" class="h-5 w-5" />
              </button>
              <RouterLink to="/home" class="surface-button-quiet">返回总览</RouterLink>
            </div>

            <button
              type="button"
              class="surface-icon-button xl:hidden"
              :aria-label="mobileMenuOpen ? '关闭导航菜单' : '打开导航菜单'"
              @click="mobileMenuOpen = !mobileMenuOpen"
            >
              <AppIcon :name="mobileMenuOpen ? 'close' : 'menu'" class="h-5 w-5" />
            </button>
          </div>

          <Transition name="menu-slide">
            <nav v-if="mobileMenuOpen" class="mt-4 grid gap-2 xl:hidden">
              <RouterLink
                v-for="item in primaryNavItems"
                :key="`mobile-${item.path}`"
                :to="item.path"
                class="surface-nav-link justify-start px-5"
                :class="{ 'surface-nav-link-active': isActive(item.path) }"
                @click="mobileMenuOpen = false"
              >
                <AppIcon :name="item.icon" class="h-4 w-4" />
                <span class="flex-1 text-left">{{ item.title }}</span>
              </RouterLink>
            </nav>
          </Transition>
        </div>
      </div>
    </header>

    <main class="relative z-10 pb-12 pt-6">
      <div class="page-shell">
        <section v-if="title || description || slots['header-extra']" class="surface-panel mb-6 p-6 md:p-8">
          <div class="grid gap-6 lg:grid-cols-[minmax(0,1fr)_280px] lg:items-center">
            <div class="space-y-4">
              <span class="eyebrow">
                <AppIcon name="spark" class="h-4 w-4" />
                {{ currentLabel }}
              </span>

              <div>
                <h1 class="section-title">{{ title }}</h1>
                <p v-if="description" class="section-copy mt-3">{{ description }}</p>
              </div>
            </div>

            <div v-if="slots['header-extra']" class="surface-deep-well p-4 md:p-5">
              <slot name="header-extra" />
            </div>

            <div v-else class="surface-deep-well hidden p-5 lg:block">
              <p class="text-xs font-semibold uppercase tracking-[0.22em] text-muted">当前模块</p>
              <p class="mt-3 font-display text-2xl font-bold text-ink">{{ title }}</p>
              <p class="mt-3 text-sm leading-6 text-muted">{{ currentSummary }}</p>
            </div>
          </div>
        </section>

        <slot />
      </div>
    </main>

    <footer class="relative z-10 px-4 pb-8 text-center text-xs text-muted sm:px-6 lg:px-8">
      <div class="page-shell px-0 sm:px-0 lg:px-0">
        <p>Easy QFNU API Lite · 中文软质感前端界面</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed, ref, useSlots, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import AppIcon from '@/components/AppIcon.vue'
import { primaryNavItems } from '@/constants/navigation'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  description: {
    type: String,
    default: ''
  },
  showBack: {
    type: Boolean,
    default: true
  }
})

const slots = useSlots()
const route = useRoute()
const router = useRouter()
const mobileMenuOpen = ref(false)

const currentNavItem = computed(() => primaryNavItems.find((item) => route.path.startsWith(item.path)))
const currentLabel = computed(() => currentNavItem.value?.shortTitle || props.title || '模块')
const currentSummary = computed(() => currentNavItem.value?.description || props.description || '在统一壳层中查看当前模块的数据与筛选条件。')

watch(
  () => route.path,
  () => {
    mobileMenuOpen.value = false
  }
)

function isActive(path) {
  return route.path === path
}
</script>

<style scoped>
.menu-slide-enter-active,
.menu-slide-leave-active {
  transition: opacity 0.28s ease, transform 0.28s ease;
}

.menu-slide-enter-from,
.menu-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
