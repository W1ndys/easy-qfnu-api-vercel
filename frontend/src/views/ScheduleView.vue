<template>
  <AppLayout
    title="课表查询"
    description="按日期查看当天课程安排、节次、地点和班级信息，页面会同步展示当前周次。"
  >
    <template #header-extra>
      <div class="grid grid-cols-2 gap-3 text-center">
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">查询日期</p>
          <p class="mt-2 text-sm font-semibold text-ink">{{ selectedDate || '--' }}</p>
        </div>
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">当前周次</p>
          <p class="mt-2 text-sm font-semibold text-ink">{{ currentWeekRaw || '未获取' }}</p>
        </div>
      </div>
    </template>

    <section class="surface-panel p-6 md:p-8">
      <div class="grid gap-4 md:grid-cols-[minmax(0,1fr)_auto] md:items-end">
        <label>
          <span class="surface-field-label">日期</span>
          <input v-model="selectedDate" type="date" class="surface-input" />
        </label>

        <button type="button" class="surface-button-primary" :disabled="loading" @click="fetchSchedule">
          <AppIcon name="search" class="h-4 w-4" />
          {{ loading ? '查询中…' : '刷新课表' }}
        </button>
      </div>
    </section>

    <section v-if="!loading && !error" class="mt-6 grid gap-4 md:grid-cols-3">
      <div class="surface-stat">
        <p class="text-xs uppercase tracking-[0.18em] text-muted">Today</p>
        <p class="mt-2 font-display text-3xl font-bold text-ink">{{ sortedCourses.length }}</p>
        <p class="mt-2 text-sm leading-6 text-muted">当前日期共安排 {{ sortedCourses.length }} 门课程。</p>
      </div>
      <div class="surface-stat">
        <p class="text-xs uppercase tracking-[0.18em] text-muted">Week</p>
        <p class="mt-2 font-display text-3xl font-bold text-ink">{{ currentWeekRaw || '--' }}</p>
        <p class="mt-2 text-sm leading-6 text-muted">接口返回的原始周次信息会直接在此展示。</p>
      </div>
      <div class="surface-stat">
        <p class="text-xs uppercase tracking-[0.18em] text-muted">Date</p>
        <p class="mt-2 font-display text-3xl font-bold text-ink">{{ selectedDate.slice(5) }}</p>
        <p class="mt-2 text-sm leading-6 text-muted">切换日期后可快速查看对应课程安排。</p>
      </div>
    </section>

    <section class="mt-6 space-y-4">
      <div v-if="loading" class="grid gap-4">
        <div class="surface-skeleton h-36"></div>
        <div class="surface-skeleton h-36"></div>
      </div>

      <div v-else-if="error" class="surface-error">
        {{ error }}
      </div>

      <div v-else-if="empty" class="surface-empty">
        当前日期暂无课程安排，可以切换其他日期继续查看。
      </div>

      <div v-else class="grid gap-4">
        <article
          v-for="course in sortedCourses"
          :key="`${course.index}-${course.name}-${course.rawTimeString}`"
          class="surface-card p-5 md:p-6"
        >
          <div class="grid gap-4 lg:grid-cols-[84px_minmax(0,1fr)] lg:items-start">
            <div class="surface-deep-well flex min-h-[84px] items-center justify-center text-accent">
              <AppIcon name="schedule" class="h-8 w-8" />
            </div>

            <div>
              <div class="flex flex-wrap items-center gap-2">
                <h2 class="font-display text-2xl font-bold tracking-tight text-ink">{{ course.name || '未命名课程' }}</h2>
                <span class="surface-badge">{{ course.category || '未分类' }}</span>
              </div>

              <div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">上课时间</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ formatCourseTime(course) }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">上课地点</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ course.location || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">班级</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ course.classes || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">学分</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ course.credit || '-' }}</p>
                </div>
              </div>
            </div>
          </div>
        </article>
      </div>
    </section>
  </AppLayout>
</template>

<script setup>
import { onMounted, ref, watch, computed } from 'vue'
import { getClassSchedules } from '@/api/zhjw'
import AppIcon from '@/components/AppIcon.vue'
import AppLayout from '@/components/AppLayout.vue'
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
