<template>
  <AppLayout title="培养方案">
    <section class="space-y-3">
      <div v-if="loading" class="space-y-3">
        <div class="h-24 animate-pulse rounded-xl bg-white shadow-sm"></div>
        <div class="h-24 animate-pulse rounded-xl bg-white shadow-sm"></div>
        <div class="h-28 animate-pulse rounded-xl bg-white shadow-sm"></div>
      </div>

      <div v-else-if="error" class="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-600">
        {{ error }}
      </div>

      <div v-else-if="empty" class="rounded-xl bg-white p-6 text-center text-sm text-gray-500 shadow-sm">
        暂无培养方案数据
      </div>

      <template v-else>
        <article class="overflow-hidden rounded-xl bg-white shadow-sm">
          <button
            type="button"
            class="flex min-h-11 w-full items-center justify-between bg-gray-50 px-4 py-3 text-left"
            @click="showObjectives = !showObjectives"
          >
            <span class="font-semibold text-gray-800">培养目标</span>
            <span class="text-xs text-gray-500">{{ showObjectives ? '收起' : '展开' }}</span>
          </button>
          <div v-show="showObjectives" class="p-4 text-sm leading-6 text-gray-600">
            {{ plan.objectives || '暂无培养目标' }}
          </div>
        </article>

        <article class="overflow-hidden rounded-xl bg-white shadow-sm">
          <button
            type="button"
            class="flex min-h-11 w-full items-center justify-between bg-gray-50 px-4 py-3 text-left"
            @click="showDescription = !showDescription"
          >
            <span class="font-semibold text-gray-800">详细说明</span>
            <span class="text-xs text-gray-500">{{ showDescription ? '收起' : '展开' }}</span>
          </button>
          <div v-show="showDescription" class="p-4 text-sm leading-6 text-gray-600">
            {{ plan.description || '暂无详细说明' }}
          </div>
        </article>

        <article
          v-for="(group, index) in groups"
          :key="`${group.group_name}-${index}`"
          class="overflow-hidden rounded-xl bg-white shadow-sm"
        >
          <button
            type="button"
            class="w-full px-4 py-3 text-left"
            @click="toggleGroup(index)"
          >
            <div class="flex items-center justify-between gap-3">
              <p class="text-sm font-semibold text-gray-800">{{ group.group_name || '未命名课程组' }}</p>
              <p class="whitespace-nowrap text-xs text-gray-500">
                {{ formatNumber(group.earned_credits) }} / {{ formatNumber(group.required_credits) }}
              </p>
            </div>
            <div class="mt-2 h-2 overflow-hidden rounded-full bg-gray-100">
              <div
                class="h-full rounded-full bg-blue-500 transition-all"
                :style="{ width: `${progressPercent(group)}%` }"
              ></div>
            </div>
          </button>

          <div v-show="isGroupOpen(index)" class="space-y-2 border-t border-gray-100 px-3 py-3">
            <div
              v-for="course in group.courses || []"
              :key="`${course.course_code}-${course.course_name}`"
              class="rounded-lg border border-gray-100 p-3"
            >
              <div class="flex items-start justify-between gap-3">
                <p class="text-sm font-semibold text-gray-900">{{ course.course_name || '未命名课程' }}</p>
                <span
                  class="rounded-full px-2 py-0.5 text-xs"
                  :class="statusClass(course.status)"
                >
                  {{ course.status || '状态未知' }}
                </span>
              </div>

              <div class="mt-2 grid grid-cols-2 gap-2 text-xs text-gray-600">
                <p>学分：{{ formatNumber(course.credits) }}</p>
                <p>学时：{{ course.hours || '-' }}</p>
                <p>开课学期：{{ course.term || '-' }}</p>
                <p>课程属性：{{ course.course_attr || '-' }}</p>
              </div>
            </div>
          </div>
        </article>
      </template>
    </section>
  </AppLayout>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getCoursePlan } from '@/api/zhjw'
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
  if (!status) return 'bg-gray-100 text-gray-600'
  if (status.includes('已修')) return 'bg-green-100 text-green-700'
  if (status.includes('未修')) return 'bg-gray-100 text-gray-600'
  return 'bg-yellow-100 text-yellow-700'
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
