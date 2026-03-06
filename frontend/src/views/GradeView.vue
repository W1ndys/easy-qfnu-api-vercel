<template>
  <AppLayout title="成绩查询">
    <section class="rounded-xl bg-white p-4 shadow-sm">
      <div class="grid gap-3">
        <label class="text-sm text-gray-600">
          学期
          <select
            v-model="filters.term"
            class="mt-1 min-h-11 w-full rounded-lg border border-gray-200 px-3 text-base focus:border-blue-500 focus:outline-none"
          >
            <option value="">全部学期</option>
            <option v-for="term in semesterOptions" :key="term" :value="term">{{ term }}</option>
          </select>
        </label>

        <label class="text-sm text-gray-600">
          课程类型
          <select
            v-model="filters.course_type"
            class="mt-1 min-h-11 w-full rounded-lg border border-gray-200 px-3 text-base focus:border-blue-500 focus:outline-none"
          >
            <option value="">全部类型</option>
            <option v-for="item in courseTypes" :key="item.id" :value="item.id">
              {{ item.name }}
            </option>
          </select>
        </label>

        <label class="text-sm text-gray-600">
          课程名称
          <input
            v-model.trim="filters.course_name"
            type="text"
            placeholder="输入课程名关键词"
            class="mt-1 min-h-11 w-full rounded-lg border border-gray-200 px-3 text-base focus:border-blue-500 focus:outline-none"
          />
        </label>

        <div>
          <p class="mb-2 text-sm text-gray-600">显示模式</p>
          <div class="grid grid-cols-2 gap-2 rounded-lg bg-gray-100 p-1">
            <button
              type="button"
              class="min-h-11 rounded-md text-sm font-medium transition-colors"
              :class="filters.display_type === 'all' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-600'"
              @click="filters.display_type = 'all'"
            >
              全部成绩
            </button>
            <button
              type="button"
              class="min-h-11 rounded-md text-sm font-medium transition-colors"
              :class="filters.display_type === 'max' ? 'bg-white text-blue-600 shadow-sm' : 'text-gray-600'"
              @click="filters.display_type = 'max'"
            >
              最好成绩
            </button>
          </div>
        </div>

        <button
          type="button"
          class="min-h-11 rounded-lg bg-blue-600 px-4 text-sm font-semibold text-white transition-colors hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-blue-300"
          :disabled="loading"
          @click="fetchGrades"
        >
          {{ loading ? '查询中...' : '查询' }}
        </button>
      </div>
    </section>

    <section class="mt-4 grid grid-cols-3 gap-3" v-if="!loading && !error">
      <div class="rounded-xl bg-white p-3 shadow-sm">
        <p class="text-xs text-gray-500">总加权绩点</p>
        <p class="mt-1 text-xl font-bold text-gray-900">{{ formatNumber(gradeData.total_stat.weighted_gpa) }}</p>
      </div>
      <div class="rounded-xl bg-white p-3 shadow-sm">
        <p class="text-xs text-gray-500">总学分</p>
        <p class="mt-1 text-xl font-bold text-gray-900">{{ formatNumber(gradeData.total_stat.total_credits) }}</p>
      </div>
      <div class="rounded-xl bg-white p-3 shadow-sm">
        <p class="text-xs text-gray-500">课程总数</p>
        <p class="mt-1 text-xl font-bold text-gray-900">{{ gradeData.total_stat.course_count || 0 }}</p>
      </div>
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
        未查询到成绩数据
      </div>

      <div v-else class="space-y-4">
        <article
          v-for="group in groupedGrades"
          :key="group.semester"
          class="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm"
        >
          <button
            type="button"
            class="flex min-h-11 w-full items-center justify-between bg-gray-50 px-4 text-left"
            @click="toggleSemester(group.semester)"
          >
            <span class="font-semibold text-gray-800">{{ group.semester }}</span>
            <span class="text-xs text-gray-500">{{ group.items.length }} 门</span>
          </button>

          <div v-show="isSemesterOpen(group.semester)" class="space-y-3 p-3">
            <div
              v-for="grade in group.items"
              :key="`${group.semester}-${grade.course_code}-${grade.course_name}`"
              class="rounded-lg border border-gray-100 p-3"
            >
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0">
                  <p class="truncate text-sm font-semibold text-gray-900">{{ grade.course_name }}</p>
                  <span class="mt-1 inline-flex rounded-full bg-blue-50 px-2 py-0.5 text-xs text-blue-600">
                    {{ grade.course_prop || '未知性质' }}
                  </span>
                </div>
                <p class="text-2xl font-bold text-blue-600">{{ grade.score || '-' }}</p>
              </div>

              <div class="mt-3 grid grid-cols-3 gap-2 text-xs text-gray-600">
                <p>学分：{{ grade.credit || '-' }}</p>
                <p>绩点：{{ grade.gpa || '-' }}</p>
                <p class="truncate">考试：{{ grade.exam_type || '-' }}</p>
              </div>
            </div>
          </div>
        </article>
      </div>
    </section>
  </AppLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import AppLayout from '@/components/AppLayout.vue'
import { getGrade } from '@/api/zhjw'
import { useUserStore } from '@/stores/user'
import { resolveRequestError } from '@/utils/requestError'

const userStore = useUserStore()

const courseTypes = [
  { id: '01', name: '公共课' },
  { id: '02', name: '公共基础课' },
  { id: '03', name: '专业基础课' },
  { id: '04', name: '专业课' },
  { id: '05', name: '专业选修课' },
  { id: '06', name: '公共选修课' },
  { id: '07', name: '专业任选课' },
  { id: '08', name: '实践教学环节' },
  { id: '09', name: '公共任选课' },
  { id: '10', name: '教师教育基础课程（必修）' },
  { id: '11', name: '专业必修课' },
  { id: '12', name: '学科基础必修课' },
  { id: '13', name: '专业方向限选课' },
  { id: '14', name: '考试报名虚拟课程' },
  { id: '15', name: '教师教育选修课程' },
  { id: '16', name: '公共必修课' }
]

const filters = reactive({
  term: '',
  course_type: '',
  course_name: '',
  display_type: 'all'
})

const loading = ref(false)
const error = ref('')
const empty = ref(false)
const semesterOpenState = ref({})
const gradeData = ref({
  grades: [],
  total_stat: { weighted_gpa: 0, total_credits: 0, course_count: 0 }
})

const semesterOptions = computed(() => {
  const set = new Set(
    (gradeData.value.grades || [])
      .map((item) => item.semester)
      .filter(Boolean)
  )
  return Array.from(set).sort().reverse()
})

const groupedGrades = computed(() => {
  const groupedMap = new Map()
  for (const item of gradeData.value.grades || []) {
    const semester = item.semester || '未知学期'
    if (!groupedMap.has(semester)) {
      groupedMap.set(semester, [])
    }
    groupedMap.get(semester).push(item)
  }
  return Array.from(groupedMap.entries())
    .sort((a, b) => b[0].localeCompare(a[0]))
    .map(([semester, items]) => ({ semester, items }))
})

function getCookie() {
  return userStore.cookie || localStorage.getItem('zhjw_cookie') || ''
}

function formatNumber(value) {
  const num = Number(value)
  return Number.isFinite(num) ? num.toFixed(2) : '0.00'
}

function toggleSemester(semester) {
  semesterOpenState.value[semester] = !semesterOpenState.value[semester]
}

function isSemesterOpen(semester) {
  return semesterOpenState.value[semester] !== false
}

async function fetchGrades() {
  loading.value = true
  error.value = ''
  empty.value = false

  const params = {}
  if (filters.term) params.term = filters.term
  if (filters.course_type) params.course_type = filters.course_type
  if (filters.course_name) params.course_name = filters.course_name
  if (filters.display_type) params.display_type = filters.display_type

  try {
    const res = await getGrade(getCookie(), params)
    if (res.code !== 200) {
      error.value = res.msg || '获取成绩失败'
      gradeData.value = { grades: [], total_stat: { weighted_gpa: 0, total_credits: 0, course_count: 0 } }
      return
    }

    gradeData.value = {
      grades: res.data?.grades || [],
      total_stat: res.data?.total_stat || { weighted_gpa: 0, total_credits: 0, course_count: 0 }
    }

    empty.value = gradeData.value.grades.length === 0
    semesterOpenState.value = {}
    for (const term of semesterOptions.value) {
      semesterOpenState.value[term] = true
    }
  } catch (err) {
    const parsed = resolveRequestError(err, '未查询到成绩数据')
    if (parsed.message) {
      error.value = parsed.message
    }
    empty.value = parsed.isEmpty
    gradeData.value = { grades: [], total_stat: { weighted_gpa: 0, total_credits: 0, course_count: 0 } }
  } finally {
    loading.value = false
  }
}

onMounted(fetchGrades)
</script>
