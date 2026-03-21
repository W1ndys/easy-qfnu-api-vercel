<template>
  <AppLayout
    title="成绩查询"
    description="支持按学期、课程类型和课程名筛选成绩，保留全部成绩/最好成绩视图，并提供自定义绩点计算模式。"
  >
    <template #header-extra>
      <div class="grid grid-cols-3 gap-3 text-center">
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">总加权绩点</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ loading || error ? '--' : formatNumber(gradeData.total_stat.weighted_gpa) }}</p>
        </div>
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">总学分</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ loading || error ? '--' : formatNumber(gradeData.total_stat.total_credits) }}</p>
        </div>
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">课程数</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ loading || error ? '--' : gradeData.total_stat.course_count || 0 }}</p>
        </div>
      </div>
    </template>

    <section class="surface-panel p-4 sm:p-5 md:p-6">
      <div class="grid gap-4 lg:grid-cols-3">
        <label>
          <span class="surface-field-label">学期</span>
          <select v-model="filters.term" class="surface-select">
            <option value="">全部学期</option>
            <option v-for="term in semesterOptions" :key="term" :value="term">{{ term }}</option>
          </select>
        </label>

        <label>
          <span class="surface-field-label">课程类型</span>
          <select v-model="filters.course_type" class="surface-select">
            <option value="">全部类型</option>
            <option v-for="item in courseTypes" :key="item.id" :value="item.id">
              {{ item.name }}
            </option>
          </select>
        </label>

        <label>
          <span class="surface-field-label">课程名称</span>
          <input
            v-model.trim="filters.course_name"
            type="text"
            placeholder="输入课程名关键词"
            class="surface-input"
          />
        </label>
      </div>

      <div class="mt-4 grid gap-3 lg:grid-cols-[minmax(0,1fr)_auto_auto] lg:items-end">
        <div>
          <p class="surface-field-label">显示模式</p>
          <div class="surface-segment sm:grid-cols-2">
            <button
              type="button"
              class="surface-segment-option"
              :class="{ 'is-active': filters.display_type === 'all' }"
              @click="filters.display_type = 'all'"
            >
              全部成绩
            </button>
            <button
              type="button"
              class="surface-segment-option"
              :class="{ 'is-active': filters.display_type === 'max' }"
              @click="filters.display_type = 'max'"
            >
              最好成绩
            </button>
          </div>
        </div>

        <button type="button" class="surface-button-primary" :disabled="loading" @click="fetchGrades">
          <AppIcon name="search" class="h-4 w-4" />
          {{ loading ? '正在查询…' : '查询成绩' }}
        </button>

        <button
          v-if="!loading && !error && !empty && gradeData.grades.length > 0"
          type="button"
          :class="customCalcMode ? 'surface-button-quiet' : 'surface-button'"
          @click="toggleCustomCalcMode"
        >
          <AppIcon name="grade" class="h-4 w-4" />
          {{ customCalcMode ? '关闭自定义计算' : '开启自定义计算' }}
        </button>
      </div>
    </section>

    <section v-if="customCalcMode" class="surface-panel sticky top-24 z-20 mt-5 p-3.5 sm:p-4 md:top-28">
      <div class="grid gap-4 md:grid-cols-[minmax(0,1fr)_auto] md:items-center">
        <div class="grid grid-cols-3 gap-3">
          <div class="surface-well px-3 py-4 text-center">
            <p class="text-xs text-muted">已选绩点</p>
            <p class="mt-2 font-display text-2xl font-bold text-success">{{ customStat.weightedGpa }}</p>
          </div>
          <div class="surface-well px-3 py-4 text-center">
            <p class="text-xs text-muted">已选学分</p>
            <p class="mt-2 font-display text-2xl font-bold text-success">{{ customStat.totalCredits }}</p>
          </div>
          <div class="surface-well px-3 py-4 text-center">
            <p class="text-xs text-muted">已选课程</p>
            <p class="mt-2 font-display text-2xl font-bold text-success">{{ customStat.courseCount }}</p>
          </div>
        </div>

        <div class="grid gap-2 sm:grid-cols-2">
          <button type="button" class="surface-button-primary" @click="selectAll">全选当前结果</button>
          <button type="button" class="surface-button-quiet" @click="deselectAll">清空选择</button>
        </div>
      </div>
    </section>

    <section v-if="!loading && !error" class="mt-5 grid gap-3.5 md:grid-cols-3">
      <div class="surface-stat">
        <p class="text-xs uppercase tracking-[0.18em] text-muted">Semesters</p>
        <p class="mt-2 font-display text-3xl font-bold text-ink">{{ groupedGrades.length }}</p>
        <p class="mt-2 text-sm leading-6 text-muted">本次结果共覆盖 {{ groupedGrades.length }} 个学期分组。</p>
      </div>
      <div class="surface-stat">
        <p class="text-xs uppercase tracking-[0.18em] text-muted">Results</p>
        <p class="mt-2 font-display text-3xl font-bold text-ink">{{ gradeData.grades.length }}</p>
        <p class="mt-2 text-sm leading-6 text-muted">当前筛选条件下共返回 {{ gradeData.grades.length }} 门课程成绩。</p>
      </div>
      <div class="surface-stat">
        <p class="text-xs uppercase tracking-[0.18em] text-muted">Selected</p>
        <p class="mt-2 font-display text-3xl font-bold text-ink">{{ selectedKeys.size }}</p>
        <p class="mt-2 text-sm leading-6 text-muted">自定义模式下已勾选 {{ selectedKeys.size }} 门课程参与计算。</p>
      </div>
    </section>

    <section class="mt-5 space-y-3.5">
      <div v-if="loading" class="grid gap-4">
        <div class="surface-skeleton h-36"></div>
        <div class="surface-skeleton h-36"></div>
      </div>

      <div v-else-if="error" class="surface-error">
        {{ error }}
      </div>

      <div v-else-if="empty" class="surface-empty">
        暂未查询到符合条件的成绩数据，请调整筛选条件后重试。
      </div>

      <div v-else class="space-y-4">
        <article v-for="group in groupedGrades" :key="group.semester" class="surface-panel p-3.5 sm:p-4">
          <button
            type="button"
            class="flex w-full items-center justify-between gap-4 rounded-2xl px-3 py-2 text-left transition-colors"
            @click="toggleSemester(group.semester)"
          >
            <div>
              <p class="font-display text-2xl font-bold tracking-tight text-ink">{{ group.semester }}</p>
              <p class="mt-1 text-sm text-muted">共 {{ group.items.length }} 门课程，学分合计 {{ semesterCredits(group.items) }}</p>
            </div>
            <span class="surface-badge-neutral">{{ isSemesterOpen(group.semester) ? '收起' : '展开' }}</span>
          </button>

          <div v-show="isSemesterOpen(group.semester)" class="mt-4 grid gap-3">
            <article
              v-for="grade in group.items"
              :key="`${group.semester}-${grade.course_code}-${grade.course_name}`"
              class="surface-card p-3.5 sm:p-4"
              :class="[
                customCalcMode ? 'surface-card-tappable' : '',
                customCalcMode && isGradeSelected(grade) ? 'surface-card-selected' : ''
              ]"
              :role="customCalcMode ? 'button' : undefined"
              :tabindex="customCalcMode ? 0 : undefined"
              @click="customCalcMode && toggleGradeSelection(grade)"
              @keyup.enter.prevent="customCalcMode && toggleGradeSelection(grade)"
              @keyup.space.prevent="customCalcMode && toggleGradeSelection(grade)"
            >
              <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_120px] lg:items-start">
                <div>
                  <div class="flex flex-wrap items-center gap-2">
                    <h3 class="font-display text-2xl font-bold tracking-tight text-ink">{{ grade.course_name || '未命名课程' }}</h3>
                    <span class="surface-badge">{{ grade.course_prop || '性质未知' }}</span>
                    <span v-if="grade.exam_type" class="surface-badge-neutral">{{ grade.exam_type }}</span>
                  </div>

                  <div class="mt-3.5 grid gap-2.5 sm:grid-cols-2 xl:grid-cols-4">
                    <div class="surface-well px-4 py-3">
                      <p class="text-xs text-muted">课程编号</p>
                      <p class="mt-2 text-sm font-semibold text-ink">{{ grade.course_code || '-' }}</p>
                    </div>
                    <div class="surface-well px-4 py-3">
                      <p class="text-xs text-muted">学分</p>
                      <p class="mt-2 text-sm font-semibold text-ink">{{ grade.credit || '-' }}</p>
                    </div>
                    <div class="surface-well px-4 py-3">
                      <p class="text-xs text-muted">绩点</p>
                      <p class="mt-2 text-sm font-semibold text-ink">{{ grade.gpa || '-' }}</p>
                    </div>
                    <div class="surface-well px-4 py-3">
                      <p class="text-xs text-muted">考试方式</p>
                      <p class="mt-2 text-sm font-semibold text-ink">{{ grade.exam_type || '-' }}</p>
                    </div>
                  </div>
                </div>

                <div class="surface-deep-well flex min-h-[110px] flex-col items-center justify-center px-3.5 py-4 text-center">
                  <p class="text-xs uppercase tracking-[0.22em] text-muted">Score</p>
                  <p class="mt-3 font-display text-5xl font-extrabold tracking-tight text-accent">{{ grade.score || '-' }}</p>
                  <p v-if="customCalcMode" class="mt-2 text-xs text-muted">
                    {{ isGradeSelected(grade) ? '已加入计算' : '点击加入计算' }}
                  </p>
                </div>
              </div>
            </article>
          </div>
        </article>
      </div>
    </section>
  </AppLayout>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { getGrade } from '@/api/zhjw'
import AppIcon from '@/components/AppIcon.vue'
import AppLayout from '@/components/AppLayout.vue'
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

const customCalcMode = ref(false)
const selectedKeys = ref(new Set())

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

function gradeKey(grade) {
  return `${grade.semester}-${grade.course_code}-${grade.course_name}`
}

function semesterCredits(items) {
  const total = items.reduce((sum, item) => sum + (Number(item.credit) || 0), 0)
  return total.toFixed(2)
}

const customStat = computed(() => {
  if (selectedKeys.value.size === 0) {
    return { weightedGpa: '0.00', totalCredits: '0.00', courseCount: 0 }
  }

  let weightedSum = 0
  let totalCredits = 0
  let courseCount = 0

  for (const g of gradeData.value.grades) {
    if (!selectedKeys.value.has(gradeKey(g))) continue
    const credit = parseFloat(g.credit) || 0
    const gpa = parseFloat(g.gpa) || 0
    if (credit > 0) {
      weightedSum += gpa * credit
      totalCredits += credit
      courseCount++
    }
  }

  return {
    weightedGpa: totalCredits > 0 ? (weightedSum / totalCredits).toFixed(2) : '0.00',
    totalCredits: totalCredits.toFixed(2),
    courseCount
  }
})

function toggleGradeSelection(grade) {
  const key = gradeKey(grade)
  const next = new Set(selectedKeys.value)
  if (next.has(key)) {
    next.delete(key)
  } else {
    next.add(key)
  }
  selectedKeys.value = next
}

function isGradeSelected(grade) {
  return selectedKeys.value.has(gradeKey(grade))
}

function selectAll() {
  const next = new Set()
  for (const g of gradeData.value.grades) {
    next.add(gradeKey(g))
  }
  selectedKeys.value = next
}

function deselectAll() {
  selectedKeys.value = new Set()
}

function toggleCustomCalcMode() {
  customCalcMode.value = !customCalcMode.value
  if (!customCalcMode.value) {
    selectedKeys.value = new Set()
  }
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
    selectedKeys.value = new Set()
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
