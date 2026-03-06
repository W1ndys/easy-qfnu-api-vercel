<template>
  <AppLayout title="选课结果">
    <section class="rounded-xl bg-white p-4 shadow-sm">
      <div class="flex gap-2">
        <input
          v-model.trim="term"
          type="text"
          placeholder="学期，如 2024-2025-1"
          class="min-h-11 flex-1 rounded-lg border border-gray-200 px-3 text-base focus:border-blue-500 focus:outline-none"
        />
        <button
          type="button"
          class="min-h-11 rounded-lg bg-blue-600 px-4 text-sm font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-blue-300"
          :disabled="loading"
          @click="fetchSelectionResults"
        >
          查询
        </button>
      </div>
    </section>

    <section class="mt-4 space-y-3">
      <div v-if="loading" class="space-y-3">
        <div class="h-28 animate-pulse rounded-xl bg-white shadow-sm"></div>
        <div class="h-28 animate-pulse rounded-xl bg-white shadow-sm"></div>
      </div>

      <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-600">
        {{ error }}
      </div>

      <div v-else-if="empty" class="rounded-xl bg-white p-6 text-center text-sm text-gray-500 shadow-sm">
        暂无选课记录
      </div>

      <article
        v-else
        v-for="item in results"
        :key="`${item.course_id}-${item.select_time}-${item.index}`"
        class="rounded-xl bg-white p-4 shadow-sm"
      >
        <div class="flex items-start justify-between gap-3">
          <h2 class="text-base font-semibold text-gray-900">{{ item.course_name || '未命名课程' }}</h2>
          <span class="rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-600">{{ item.credit || '-' }} 学分</span>
        </div>

        <div class="mt-3 flex flex-wrap gap-2">
          <span class="rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600">{{ item.course_attr || '属性未知' }}</span>
          <span class="rounded-full bg-green-100 px-2 py-0.5 text-xs text-green-700">{{ item.course_prop || '性质未知' }}</span>
        </div>

        <div class="mt-3 grid grid-cols-2 gap-2 text-sm text-gray-600">
          <p>教师：{{ item.teacher || '-' }}</p>
          <p>学时：{{ item.hours || '-' }}</p>
          <p>课程编号：{{ item.course_id || '-' }}</p>
          <p>操作人：{{ item.operator || '-' }}</p>
        </div>

        <p class="mt-2 text-sm text-gray-500">选课时间：{{ item.select_time || '-' }}</p>
      </article>
    </section>
  </AppLayout>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getSelectionResults } from '@/api/zhjw'
import { useUserStore } from '@/stores/user'
import { resolveRequestError } from '@/utils/requestError'

const userStore = useUserStore()

const term = ref('')
const loading = ref(false)
const error = ref('')
const empty = ref(false)
const results = ref([])

function getCookie() {
  return userStore.cookie || localStorage.getItem('zhjw_cookie') || ''
}

async function fetchSelectionResults() {
  loading.value = true
  error.value = ''
  empty.value = false

  try {
    const params = term.value ? { term: term.value } : {}
    const res = await getSelectionResults(getCookie(), params)
    if (res.code !== 200) {
      error.value = res.msg || '获取选课结果失败'
      results.value = []
      return
    }

    const list = Array.isArray(res.data) ? res.data : res.data?.results || []
    results.value = list
    empty.value = list.length === 0
  } catch (err) {
    const parsed = resolveRequestError(err, '暂无选课记录')
    if (parsed.message) {
      error.value = parsed.message
    }
    empty.value = parsed.isEmpty
    results.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchSelectionResults)
</script>
