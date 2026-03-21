<template>
  <AppLayout
    title="选课结果"
    description="查看指定学期的选课记录、课程属性、教师和选课时间，便于快速确认已选课程。"
  >
    <template #header-extra>
      <div class="grid grid-cols-2 gap-3 text-center">
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">结果数量</p>
          <p class="mt-2 font-display text-xl font-bold text-ink">{{ results.length }}</p>
        </div>
        <div class="surface-well px-3 py-4">
          <p class="text-xs text-muted">当前学期</p>
          <p class="mt-2 text-sm font-semibold text-ink">{{ term || '全部学期' }}</p>
        </div>
      </div>
    </template>

    <section class="surface-panel p-4 sm:p-5 md:p-6">
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

        <button type="button" class="surface-button-primary" :disabled="loading" @click="fetchSelectionResults">
          <AppIcon name="search" class="h-4 w-4" />
          {{ loading ? '查询中…' : '查询选课结果' }}
        </button>
      </div>
    </section>

    <section class="mt-5 space-y-3.5">
      <div v-if="loading" class="grid gap-4">
        <div class="surface-skeleton h-40"></div>
        <div class="surface-skeleton h-40"></div>
      </div>

      <div v-else-if="error" class="surface-error">
        {{ error }}
      </div>

      <div v-else-if="empty" class="surface-empty">
        暂未查询到选课记录，可以切换学期后再次尝试。
      </div>

      <div v-else class="grid gap-4">
        <article
          v-for="item in results"
          :key="`${item.course_id}-${item.select_time}-${item.index}`"
          class="surface-card p-4 sm:p-5"
        >
          <div class="grid gap-4 lg:grid-cols-[minmax(0,1fr)_220px] lg:items-start">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <h2 class="font-display text-2xl font-bold tracking-tight text-ink">{{ item.course_name || '未命名课程' }}</h2>
                <span class="surface-badge">{{ item.credit || '-' }} 学分</span>
                <span class="surface-badge-success">{{ item.course_prop || '性质未知' }}</span>
              </div>

              <div class="mt-3">
                <span class="surface-badge-neutral">{{ item.course_attr || '属性未知' }}</span>
              </div>

              <div class="mt-3.5 grid gap-2.5 md:grid-cols-2 xl:grid-cols-4">
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">授课教师</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.teacher || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">学时</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.hours || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">课程编号</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.course_id || '-' }}</p>
                </div>
                <div class="surface-well px-4 py-3">
                  <p class="text-xs text-muted">操作人</p>
                  <p class="mt-2 text-sm font-semibold text-ink">{{ item.operator || '-' }}</p>
                </div>
              </div>
            </div>

            <div class="surface-deep-well flex min-h-[148px] flex-col justify-center px-4 py-4">
              <p class="text-xs uppercase tracking-[0.22em] text-muted">Selected At</p>
              <p class="mt-3 text-base font-semibold leading-8 text-accent">{{ item.select_time || '-' }}</p>
            </div>
          </div>
        </article>
      </div>
    </section>
  </AppLayout>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { getSelectionResults } from '@/api/zhjw'
import AppIcon from '@/components/AppIcon.vue'
import AppLayout from '@/components/AppLayout.vue'
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
