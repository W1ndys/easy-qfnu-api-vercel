<template>
  <div class="min-h-screen bg-gray-50">
    <div class="mx-auto flex min-h-screen max-w-lg flex-col">
      <header class="sticky top-0 z-20 border-b border-gray-200 bg-white/95 backdrop-blur">
        <div class="grid h-12 grid-cols-[44px_1fr_44px] items-center px-2">
          <button
            v-if="showBack"
            type="button"
            class="flex h-11 w-11 items-center justify-center rounded-md text-xl text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900"
            @click="router.back()"
            aria-label="返回"
          >
            ←
          </button>
          <div v-else class="h-11 w-11"></div>

          <h1 class="truncate text-center text-base font-semibold text-gray-800">
            {{ title }}
          </h1>

          <div class="h-11 w-11"></div>
        </div>
      </header>

      <main class="flex-1 px-4 py-4 pb-[max(env(safe-area-inset-bottom),1rem)]">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

defineProps({
  title: {
    type: String,
    default: ''
  },
  showBack: {
    type: Boolean,
    default: true
  }
})

const router = useRouter()
</script>
