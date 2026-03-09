<template>
  <div class="min-h-screen bg-gray-50">
    <div class="mx-auto flex min-h-screen max-w-lg flex-col px-4 py-6 pb-[max(env(safe-area-inset-bottom),1rem)]">
      <header class="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">曲奇助手</h1>
          <p class="mt-2 text-sm text-gray-500">快速访问常用教务查询功能</p>
        </div>
        <button
          type="button"
          class="min-h-11 rounded-lg border border-gray-200 bg-white px-4 text-sm font-medium text-gray-700 shadow-sm transition-colors hover:bg-gray-100"
          @click="handleLogout"
        >
          退出登录
        </button>
      </header>

      <div class="mb-6 rounded-xl border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
        该网站目前处于开发初期，可能不稳定。
      </div>

      <main class="flex-1">
        <div class="grid grid-cols-2 gap-4 sm:grid-cols-3">
          <button
            v-for="item in featureCards"
            :key="item.path"
            type="button"
            class="aspect-square rounded-xl bg-white p-4 text-left shadow-sm transition-transform duration-150 hover:shadow active:scale-95"
            @click="router.push(item.path)"
          >
            <div class="flex h-full flex-col items-center justify-center gap-2">
              <span class="text-3xl" aria-hidden="true">{{ item.icon }}</span>
              <span class="text-sm font-semibold text-gray-800">{{ item.title }}</span>
              <span class="text-center text-xs text-gray-500">{{ item.desc }}</span>
            </div>
          </button>
        </div>
      </main>

      <footer class="pt-8 text-center text-xs text-gray-400">easy-qfnu-api-lite · v1</footer>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const featureCards = [
  { title: '成绩查询', icon: '📊', path: '/grade', desc: '历史成绩与绩点' },
  { title: '课表查询', icon: '📅', path: '/schedule', desc: '按日期查看课程' },
  { title: '考试安排', icon: '📝', path: '/exam', desc: '考试时间与地点' },
  { title: '选课结果', icon: '📋', path: '/selection', desc: '已选课程记录' },
  { title: '培养方案', icon: '📚', path: '/course-plan', desc: '培养计划与学分' }
]

function handleLogout() {
  userStore.clearCookie()
  router.replace('/login')
}
</script>
