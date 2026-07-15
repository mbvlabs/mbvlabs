<script setup lang="ts">
import { router } from "@inertiajs/vue3"
import type { HTMLAttributes } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps<{
  class?: HTMLAttributes["class"]
  href?: string
  ariaLabel?: string
}>()

function visit() {
  if (props.href) router.visit(props.href)
}

function handleKeydown(event: KeyboardEvent) {
  if (event.target !== event.currentTarget || (event.key !== "Enter" && event.key !== " ")) return
  event.preventDefault()
  visit()
}
</script>

<template>
  <tr
    data-slot="table-row"
    :role="href ? 'link' : undefined"
    :tabindex="href ? 0 : undefined"
    :aria-label="ariaLabel"
    :class="cn(
      'hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors has-aria-expanded:bg-muted/50',
      href && 'cursor-pointer focus-visible:bg-muted/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring',
      props.class,
    )"
    @click="visit"
    @keydown="handleKeydown"
  >
    <slot />
  </tr>
</template>
