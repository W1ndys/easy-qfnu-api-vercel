<template>
  <AppLayout
    title="考试安排"
    description="按学期查看考试时间、地点、座位号和备注信息，结果会自动按考试时间排序。"
  >
    <template #header-extra>
      <div class="grid grid-cols-2 gap-3 text-center">
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">结果数量</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ exams.length }}</p>
        </div>
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">当前学期</p>
          <p class="mt-2 text-sm font-semibold text-ink">{{ term || '全部学期' }}</p>
        </div>
      </div>
    </template>

    <section class="surface-panel p-6 md:p-8">
      <div class="grid gap-4 md:grid-cols-[minmax(0,1fr)_auto] md:items-end">
        <label>
          <span class="surface-field-label">学期筛选</span>
          <input
            v-model.trim="term"
            type="text"
            placeholder="例如：2024-2025-1"
            class="surface-input"
          />
        </label>

        <button type="button" class="surface-button-primary" :disabled="loading" @click="fetchExams">
          <AppIcon name="search" class="h-4 w-4" />
          {{ loading ? '查询中…' : '查询考试' }}
        </button>
      </div>
    </section>

    <section class="mt-6 space-y-4">
      <div v-if="loading" class="grid gap-4">
        <div class="surface-skeleton h-40"></div>
        <div class="surface-skeleton h-40"></div>
      </div>

      <div v-else-if="error" class="surface-error">
        {{ error }}
      </div>

      <div v-else-if="empty" class="surface-empty">
        当前条件下暂无考试安排，可以尝试更换学期后再次查询。
      </div>

      <div v-else class="grid gap-4">
        <article
          v-for="item in exams"
          :key="`${item.course_id}-${item.exam_time}-${item.index}`"
          class="surface-card p-5 md:p-6"
        >
          <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_220px] lg:items-start">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <h2 class="font-display text-2xl font-bold tracking-tight text-ink">{{ item.course_name || '未命名课程' }}</h2>
                <span class="surface-badge">{{ item.session || '场次未定' }}</span>
              </div>

              <div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">考试地点</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.exam_room || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">座位号</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.seat_number || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">校区</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.campus || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">准考证号</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.admission_no || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">任课教师</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.instructor || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">课程编号</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.course_id || '-' }}</p>
                </div>
              </div>

              <div v-if="item.remarks" class="surface-warning mt-4">
                备注：{{ item.remarks }}
              </div>
            </div>

            <div class="surface-deep-well flex min-h-[172px] flex-col justify-center px-5 py-6">
              <p class="text-xs uppercase tracking-[0.22em] text-muted">Exam Time</p>
              <p class="mt-3 font-display text-3xl font-bold tracking-tight text-accent">{{ item.exam_time || '-' }}</p>
            </div>
          </div>
        </article>
      </div>
    </section>
  </AppLayout>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getExamSchedules } from '@/api/zhjw'
import AppIcon from '@/components/AppIcon.vue'
import AppLayout from '@/components/AppLayout.vue'
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
