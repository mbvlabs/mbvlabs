<script setup lang="ts">
import { Moon, Sun } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { onMounted, ref } from 'vue'

const isDark = ref(true)

function setTheme(dark: boolean) {
  isDark.value = dark
  const root = document.documentElement
  if (dark) {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
  localStorage.setItem('theme', dark ? 'dark' : 'light')
}

function toggleTheme() {
  setTheme(!isDark.value)
}

onMounted(() => {
  const saved = localStorage.getItem('theme')
  if (saved === 'light') {
    setTheme(false)
  } else {
    setTheme(true)
  }
})
</script>

<template>
  <Button variant="ghost" size="icon" @click="toggleTheme">
    <Sun v-if="!isDark" class="size-4" />
    <Moon v-else class="size-4" />
    <span class="sr-only">Toggle theme</span>
  </Button>
</template>
