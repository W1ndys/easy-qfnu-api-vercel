<template>
  <AppLayout title="考试安排">
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
          @click="fetchExams"
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
        暂无考试安排
      </div>

      <article
        v-else
        v-for="item in exams"
        :key="`${item.course_id}-${item.exam_time}-${item.index}`"
        class="rounded-xl bg-white p-4 shadow-sm"
      >
        <div class="flex items-start justify-between gap-3">
          <h2 class="text-base font-semibold text-gray-900">{{ item.course_name || '未命名课程' }}</h2>
          <span class="rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-600">{{ item.session || '场次未定' }}</span>
        </div>

        <p class="mt-2 text-sm font-medium text-blue-600">{{ item.exam_time || '-' }}</p>

        <div class="mt-3 grid grid-cols-2 gap-2 text-sm text-gray-600">
          <p>地点：{{ item.exam_room || '-' }}</p>
          <p>座位号：{{ item.seat_number || '-' }}</p>
          <p>校区：{{ item.campus || '-' }}</p>
          <p>准考证：{{ item.admission_no || '-' }}</p>
          <p>任课教师：{{ item.instructor || '-' }}</p>
          <p>课程编号：{{ item.course_id || '-' }}</p>
        </div>

        <p v-if="item.remarks" class="mt-2 text-sm text-gray-500">备注：{{ item.remarks }}</p>
      </article>
    </section>
  </AppLayout>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getExamSchedules } from '@/api/zhjw'
import { useUserStore } from '@/stores/user'
import { resolveRequestError } from '@/utils/requestError'

const userStore = useUserStore()

const term = ref('')
const loading = ref(false)
const error = ref('')
const empty = ref(false)
const exams = ref([])

function getCookie() {
  return userStore.cookie || localStorage.getItem('zhjw_cookie') || ''
}

function parseExamTimestamp(raw) {
  if (!raw) return Number.MAX_SAFE_INTEGER

  const matcher = raw.match(/(\d{4})[^\d]?(\d{1,2})[^\d]?(\d{1,2}).*?(\d{1,2}):(\d{1,2})/)
  if (matcher) {
    const year = Number(matcher[1])
    const month = Number(matcher[2]) - 1
    const day = Number(matcher[3])
    const hour = Number(matcher[4])
    const minute = Number(matcher[5])
    return new Date(year, month, day, hour, minute).getTime()
  }

  const fallback = Date.parse(raw)
  return Number.isNaN(fallback) ? Number.MAX_SAFE_INTEGER : fallback
}

async function fetchExams() {
  loading.value = true
  error.value = ''
  empty.value = false

  try {
    const params = term.value ? { term: term.value } : {}
    const res = await getExamSchedules(getCookie(), params)
    if (res.code !== 200) {
      error.value = res.msg || '获取考试安排失败'
      exams.value = []
      return
    }

    exams.value = (res.data || []).slice().sort((a, b) => parseExamTimestamp(a.exam_time) - parseExamTimestamp(b.exam_time))
    empty.value = exams.value.length === 0
  } catch (err) {
    const parsed = resolveRequestError(err, '暂无考试安排')
    if (parsed.message) {
      error.value = parsed.message
    }
    empty.value = parsed.isEmpty
    exams.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchExams)
</script>
