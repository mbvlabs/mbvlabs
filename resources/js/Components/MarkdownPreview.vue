<script setup lang="ts">
import type { HTMLAttributes } from 'vue'
import { computed } from 'vue'
import { marked } from 'marked'

import { cn } from '@/lib/utils'

const props = defineProps<{
  source?: string | null
  class?: HTMLAttributes['class']
}>()

const html = computed(() => {
  return marked.parse(props.source ?? '', { async: false }) as string
})
</script>

<template>
  <article
    :class="cn('prose prose-sm max-w-none dark:prose-invert', props.class)"
    v-html="html"
  />
</template>
