<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { useAttrs } from "vue"
import { Check, ChevronDown } from "@lucide/vue"
import {
  SelectContent,
  SelectIcon,
  SelectItem,
  SelectItemIndicator,
  SelectItemText,
  SelectPortal,
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectViewport,
} from "reka-ui"
import { cn } from "@/lib/utils"

defineOptions({
  inheritAttrs: false,
})

type SelectOption = {
  label: string
  value: string
}

const props = withDefaults(defineProps<{
  modelValue?: string
  options: SelectOption[]
  placeholder?: string
  class?: HTMLAttributes["class"]
}>(), {
  placeholder: "Select an option",
})

const emits = defineEmits<{
  (e: "update:modelValue", payload: string): void
}>()

const attrs = useAttrs()
</script>

<template>
  <SelectRoot
    :model-value="modelValue"
    @update:model-value="emits('update:modelValue', String($event ?? ''))"
  >
    <SelectTrigger
      v-bind="attrs"
      :class="cn(
        'border-transparent border-b-input bg-transparent focus-visible:border-b-ring aria-invalid:border-b-destructive dark:aria-invalid:border-b-destructive/50 flex h-10 w-full min-w-0 items-center justify-between border px-0 py-1 text-left text-base transition-[color,border-color] md:text-sm outline-none placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 data-[placeholder]:text-muted-foreground',
        props.class,
      )"
    >
      <SelectValue :placeholder="placeholder" />
      <SelectIcon as-child>
        <ChevronDown class="size-4 opacity-60" />
      </SelectIcon>
    </SelectTrigger>
    <SelectPortal>
      <SelectContent
        position="popper"
        :side-offset="4"
        class="z-50 min-w-[var(--reka-select-trigger-width)] overflow-hidden border bg-popover text-popover-foreground shadow-md"
      >
        <SelectViewport class="p-1">
          <SelectItem
            v-for="option in options"
            :key="option.value"
            :value="option.value"
            class="relative flex min-h-9 cursor-default select-none items-center py-1.5 pl-8 pr-2 text-sm outline-none data-[highlighted]:bg-muted data-[highlighted]:text-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50"
          >
            <SelectItemIndicator class="absolute left-2 flex size-4 items-center justify-center">
              <Check class="size-4" />
            </SelectItemIndicator>
            <SelectItemText>{{ option.label }}</SelectItemText>
          </SelectItem>
        </SelectViewport>
      </SelectContent>
    </SelectPortal>
  </SelectRoot>
</template>
