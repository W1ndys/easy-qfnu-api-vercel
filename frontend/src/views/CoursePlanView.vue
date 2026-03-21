<template>
  <AppLayout
    title="培养方案"
    description="按课程组查看毕业要求、已修学分与课程状态，并保留培养目标和详细说明的展开阅读。"
  >
    <template #header-extra>
      <div class="grid grid-cols-2 gap-3 text-center">
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">已修学分</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ summaryStats.earned }}</p>
        </div>
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">课程组</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ summaryStats.totalGroups }}</p>
        </div>
      </div>
    </template>

    <section class="space-y-4">
      <div v-if="loading" class="grid gap-4">
        <div class="surface-skeleton h-36"></div>
        <div class="surface-skeleton h-36"></div>
        <div class="surface-skeleton h-44"></div>
      </div>

      <div v-else-if="error" class="surface-error">
        {{ error }}
      </div>

      <div v-else-if="empty" class="surface-empty">
        当前暂无培养方案数据。
      </div>

      <template v-else>
        <section class="grid gap-3.5 md:grid-cols-3">
          <div class="surface-stat">
            <p class="text-xs uppercase tracking-[0.18em] text-muted">Earned</p>
            <p class="mt-2 font-display text-3xl font-bold text-ink">{{ summaryStats.earned }}</p>
            <p class="mt-2 text-sm leading-6 text-muted">课程组累计已修学分。</p>
          </div>
          <div class="surface-stat">
            <p class="text-xs uppercase tracking-[0.18em] text-muted">Required</p>
            <p class="mt-2 font-display text-3xl font-bold text-ink">{{ summaryStats.required }}</p>
            <p class="mt-2 text-sm leading-6 text-muted">所有课程组要求学分总和。</p>
          </div>
          <div class="surface-stat">
            <p class="text-xs uppercase tracking-[0.18em] text-muted">Completed</p>
            <p class="mt-2 font-display text-3xl font-bold text-ink">{{ summaryStats.completedGroups }}</p>
            <p class="mt-2 text-sm leading-6 text-muted">已达成要求的课程组数量。</p>
          </div>
        </section>

        <article class="surface-panel p-4 sm:p-5">
          <button type="button" class="flex w-full items-center justify-between gap-4 text-left" @click="showObjectives = !showObjectives">
            <div>
              <p class="text-xs font-semibold uppercase tracking-[0.22em] text-muted">Objectives</p>
              <h2 class="mt-2 font-display text-2xl font-bold tracking-tight text-ink">培养目标</h2>
            </div>
            <span class="surface-badge-neutral">{{ showObjectives ? '收起' : '展开' }}</span>
          </button>
          <div v-show="showObjectives" class="surface-text-block mt-4">
            {{ plan.objectives || '暂无培养目标' }}
          </div>
        </article>

        <article class="surface-panel p-4 sm:p-5">
          <button type="button" class="flex w-full items-center justify-between gap-4 text-left" @click="showDescription = !showDescription">
            <div>
              <p class="text-xs font-semibold uppercase tracking-[0.22em] text-muted">Description</p>
              <h2 class="mt-2 font-display text-2xl font-bold tracking-tight text-ink">详细说明</h2>
            </div>
            <span class="surface-badge-neutral">{{ showDescription ? '收起' : '展开' }}</span>
          </button>
          <div v-show="showDescription" class="surface-text-block mt-4">
            {{ plan.description || '暂无详细说明' }}
          </div>
        </article>

        <article
          v-for="(group, index) in groups"
          :key="`${group.group_name}-${index}`"
          class="surface-panel p-4 sm:p-5"
        >
          <button type="button" class="w-full text-left" @click="toggleGroup(index)">
            <div class="flex flex-wrap items-center justify-between gap-3">
              <div>
                <p class="font-display text-2xl font-bold tracking-tight text-ink">{{ group.group_name || '未命名课程组' }}</p>
                <p class="mt-2 text-sm leading-6 text-muted">
                  已修 {{ formatNumber(group.earned_credits) }} / 要求 {{ formatNumber(group.required_credits) }} 学分
                </p>
              </div>
              <span class="surface-badge-neutral">{{ isGroupOpen(index) ? '收起课程' : '展开课程' }}</span>
            </div>

            <div class="mt-4 surface-progress-track">
              <div class="surface-progress-bar" :style="{ width: `${progressPercent(group)}%` }"></div>
            </div>
          </button>

          <div v-show="isGroupOpen(index)" class="mt-4 grid gap-3">
            <article
              v-for="course in group.courses || []"
              :key="`${course.course_code}-${course.course_name}`"
              class="surface-card p-3.5 sm:p-4"
            >
              <div class="flex flex-wrap items-center justify-between gap-3">
                <h3 class="font-semibold text-ink">{{ course.course_name || '未命名课程' }}</h3>
                <span :class="statusClass(course.status)">
                  {{ course.status || '状态未知' }}
                </span>
              </div>

              <div class="mt-4 grid gap-3 md:grid-cols-2 xl:grid-cols-4">
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">学分</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ formatNumber(course.credits) }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">学时</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ course.hours || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">开课学期</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ course.term || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">课程属性</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ course.course_attr || '-' }}</p>
                </div>
              </div>
            </article>
          </div>
        </article>
      </template>
    </section>
  </AppLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { getCoursePlan } from '@/api/zhjw'
import AppLayout from '@/components/AppLayout.vue'
import { useUserStore } from '@/stores/user'
import { resolveRequestError } from '@/utils/requestError'

const userStore = useUserStore()

const loading = ref(false)
const error = ref('')
const empty = ref(false)
const showObjectives = ref(false)
const showDescription = ref(false)
const openGroups = ref({})
const plan = ref({
  objectives: '',
  description: '',
  groups: []
})

const groups = computed(() => plan.value.groups || [])

const summaryStats = computed(() => {
  let required = 0
  let earned = 0
  let completedGroups = 0

  for (const group of groups.value) {
    const groupRequired = Number(group.required_credits) || 0
    const groupEarned = Number(group.earned_credits) || 0
    required += groupRequired
    earned += groupEarned
    if (groupRequired > 0 && groupEarned >= groupRequired) {
      completedGroups++
    }
  }

  return {
    required: required.toFixed(2),
    earned: earned.toFixed(2),
    completedGroups,
    totalGroups: groups.value.length
  }
})

function getCookie() {
  return userStore.cookie || localStorage.getItem('zhjw_cookie') || ''
}

function formatNumber(value) {
  const num = Number(value)
  return Number.isFinite(num) ? num.toFixed(2) : '0.00'
}

function progressPercent(group) {
  const required = Number(group.required_credits) || 0
  const earned = Number(group.earned_credits) || 0
  if (required <= 0) return 0
  return Math.min(100, Math.round((earned / required) * 100))
}

function statusClass(status) {
  if (!status) return 'surface-badge-neutral'
  if (status.includes('已修')) return 'surface-badge-success'
  if (status.includes('未修')) return 'surface-badge-neutral'
  return 'surface-badge-warm'
}

function toggleGroup(index) {
  openGroups.value[index] = !openGroups.value[index]
}

function isGroupOpen(index) {
  return openGroups.value[index] === true
}

async function fetchPlan() {
  loading.value = true
  error.value = ''
  empty.value = false

  try {
    const res = await getCoursePlan(getCookie())
    if (res.code !== 200) {
      error.value = res.msg || '获取培养方案失败'
      plan.value = { objectives: '', description: '', groups: [] }
      return
    }

    plan.value = {
      objectives: res.data?.objectives || '',
      description: res.data?.description || '',
      groups: res.data?.groups || []
    }
    empty.value = plan.value.groups.length === 0
    openGroups.value = {}
  } catch (err) {
    const parsed = resolveRequestError(err, '暂无培养方案数据')
    if (parsed.message) {
      error.value = parsed.message
    }
    empty.value = parsed.isEmpty
    plan.value = { objectives: '', description: '', groups: [] }
  } finally {
    loading.value = false
  }
}

onMounted(fetchPlan)
</script>
