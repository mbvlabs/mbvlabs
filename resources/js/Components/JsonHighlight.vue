<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  value: string
  class?: string
}>(), {
  class: '',
})

const highlighted = computed(() => highlightJson(props.value))

function formatJson(value: string): string {
  if (!value) {
    return '{}'
  }

  try {
    return JSON.stringify(JSON.parse(value), null, 2)
  } catch {
    return value
  }
}

function escapeHtml(value: string): string {
  return value
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;')
}

function highlightJson(value: string): string {
  return escapeHtml(formatJson(value)).replace(
    /(&quot;(?:\\u[a-fA-F0-9]{4}|\\[^u]|[^\\&])*&quot;(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g,
    (match) => {
      let className = 'text-foreground'
      if (match.endsWith(':')) {
        className = 'text-sky-700 dark:text-sky-300'
      } else if (match === 'true' || match === 'false') {
        className = 'text-violet-700 dark:text-violet-300'
      } else if (match === 'null') {
        className = 'text-muted-foreground'
      } else if (match.startsWith('&quot;')) {
        className = 'text-emerald-700 dark:text-emerald-300'
      } else {
        className = 'text-amber-700 dark:text-amber-300'
      }

      return `<span class="${className}">${match}</span>`
    },
  )
}
</script>

<template>
  <pre
    :class="['overflow-auto whitespace-pre-wrap font-mono text-xs leading-5', props.class]"
    v-html="highlighted"
  />
</template>
