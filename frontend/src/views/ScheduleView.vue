<template>
  <AppLayout title="课表查询">
    <section class="rounded-xl bg-white p-4 shadow-sm">
      <label class="text-sm text-gray-600">
        日期
        <input
          v-model="selectedDate"
          type="date"
          class="mt-1 min-h-11 w-full rounded-lg border border-gray-200 px-3 text-base focus:border-blue-500 focus:outline-none"
        />
      </label>
    </section>

    <section class="mt-4 rounded-xl bg-white p-4 shadow-sm" v-if="!loading && !error">
      <p class="text-sm text-gray-500">当前周次</p>
      <p class="mt-1 text-lg font-semibold text-gray-900">{{ currentWeekRaw || '未知' }}</p>
    </section>

    <section class="mt-4 space-y-3">
      <div v-if="loading" class="space-y-3">
        <div class="h-24 animate-pulse rounded-xl bg-white shadow-sm"></div>
        <div class="h-24 animate-pulse rounded-xl bg-white shadow-sm"></div>
      </div>

      <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-600">
        {{ error }}
      </div>

      <div v-else-if="empty" class="rounded-xl bg-white p-6 text-center text-sm text-gray-500 shadow-sm">
        今日无课
      </div>

      <article
        v-else
        v-for="course in sortedCourses"
        :key="`${course.index}-${course.name}-${course.rawTimeString}`"
        class="rounded-xl bg-white p-4 shadow-sm"
      >
        <div class="flex items-start justify-between gap-3">
          <h2 class="text-base font-semibold text-gray-900">{{ course.name || '未命名课程' }}</h2>
          <span class="rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-600">
            {{ course.category || '未分类' }}
          </span>
        </div>

        <div class="mt-3 space-y-2 text-sm text-gray-600">
          <p>时间：{{ formatCourseTime(course) }}</p>
          <p>地点：{{ course.location || '-' }}</p>
          <p>班级：{{ course.classes || '-' }}</p>
          <p>学分：{{ course.credit || '-' }}</p>
        </div>
      </article>
    </section>
  </AppLayout>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getClassSchedules } from '@/api/zhjw'
import { useUserStore } from '@/stores/user'
import { resolveRequestError } from '@/utils/requestError'

const userStore = useUserStore()

const selectedDate = ref(getToday())
const loading = ref(false)
const error = ref('')
const empty = ref(false)
const currentWeekRaw = ref('')
const courses = ref([])

const dayMap = {
  1: '一',
  2: '二',
  3: '三',
  4: '四',
  5: '五',
  6: '六',
  7: '日'
}

const sortedCourses = computed(() => {
  return [...courses.value].sort((a, b) => {
    const dayDiff = (a.timeParsed?.dayOfWeek || 0) - (b.timeParsed?.dayOfWeek || 0)
    if (dayDiff !== 0) return dayDiff
    return (a.timeParsed?.periodArray?.[0] || 0) - (b.timeParsed?.periodArray?.[0] || 0)
  })
})

function getToday() {
  const now = new Date()
  const y = now.getFullYear()
  const m = String(now.getMonth() + 1).padStart(2, '0')
  const d = String(now.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function getCookie() {
  return userStore.cookie || localStorage.getItem('zhjw_cookie') || ''
}

function formatCourseTime(course) {
  const day = course.timeParsed?.dayOfWeek
  const periods = course.timeParsed?.periodArray || []
  if (day && periods.length) {
    const start = periods[0]
    const end = periods[periods.length - 1]
    return `星期${dayMap[day]} 第${start}-${end}节`
  }
  return course.rawTimeString || '-'
}

async function fetchSchedule() {
  loading.value = true
  error.value = ''
  empty.value = false

  try {
    const res = await getClassSchedules(getCookie(), { date: selectedDate.value })
    if (res.code !== 200) {
      error.value = res.msg || '获取课表失败'
      courses.value = []
      currentWeekRaw.value = ''
      return
    }

    currentWeekRaw.value = res.data?.currentWeekRaw || ''
    courses.value = res.data?.courses || []
    empty.value = courses.value.length === 0
  } catch (err) {
    const parsed = resolveRequestError(err, '今日无课')
    if (parsed.message) {
      error.value = parsed.message
    }
    empty.value = parsed.isEmpty
    courses.value = []
    currentWeekRaw.value = ''
  } finally {
    loading.value = false
  }
}

watch(selectedDate, fetchSchedule)
onMounted(fetchSchedule)
</script>
