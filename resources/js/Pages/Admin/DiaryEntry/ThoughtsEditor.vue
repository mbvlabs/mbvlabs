<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import RichEditor from '@/Components/RichEditor.vue'
import { Label } from '@/components/ui/label'
import { Toggle } from '@/components/ui/toggle'

type ThoughtPeriod = 'morning' | 'evening'

const props = withDefaults(defineProps<{
  morningThoughts: string
  eveningThoughts: string
  morningError?: string
  eveningError?: string
  initialPeriod?: ThoughtPeriod | ''
}>(), {
  morningError: '',
  eveningError: '',
  initialPeriod: 'morning',
})

const emit = defineEmits<{
  'update:morningThoughts': [value: string]
  'update:eveningThoughts': [value: string]
}>()

const activePeriod = ref<ThoughtPeriod>(props.initialPeriod === 'evening' ? 'evening' : 'morning')

const morningModel = computed({
  get: () => props.morningThoughts,
  set: (value: string) => emit('update:morningThoughts', value),
})

const eveningModel = computed({
  get: () => props.eveningThoughts,
  set: (value: string) => emit('update:eveningThoughts', value),
})

const activeModel = computed({
  get: () => activePeriod.value === 'morning' ? morningModel.value : eveningModel.value,
  set: (value: string) => {
    if (activePeriod.value === 'morning') {
      morningModel.value = value
      return
    }
    eveningModel.value = value
  },
})

const activeLabel = computed(() => activePeriod.value === 'morning' ? 'Morning Thoughts' : 'Evening Thoughts')
const activeError = computed(() => activePeriod.value === 'morning' ? props.morningError : props.eveningError)
const activePlaceholder = computed(() => `Write the ${activePeriod.value} thoughts...`)

watch(() => props.initialPeriod, (period) => {
  if (period === 'morning' || period === 'evening') {
    activePeriod.value = period
  }
})

function selectPeriod(period: ThoughtPeriod) {
  activePeriod.value = period
}

function periodToggleClass(period: ThoughtPeriod) {
  return [
    'min-w-28 border transition-colors',
    activePeriod.value === period
      ? 'bg-primary text-primary-foreground hover:bg-primary/90 aria-pressed:bg-primary aria-pressed:text-primary-foreground'
      : 'bg-background text-muted-foreground hover:bg-muted hover:text-foreground',
  ]
}
</script>

<template>
  <section class="space-y-3">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
      <div class="space-y-1">
        <Label :for="`${activePeriod}Thoughts`">{{ activeLabel }}</Label>
      </div>

      <div class="inline-flex w-fit border bg-background" role="group" aria-label="Thought period">
        <Toggle
          type="button"
          size="sm"
          :model-value="activePeriod === 'morning'"
          :aria-label="'Show morning thoughts editor'"
          :class="periodToggleClass('morning')"
          @update:model-value="selectPeriod('morning')"
        >
          Morning
        </Toggle>
        <Toggle
          type="button"
          size="sm"
          :model-value="activePeriod === 'evening'"
          :aria-label="'Show evening thoughts editor'"
          :class="periodToggleClass('evening')"
          @update:model-value="selectPeriod('evening')"
        >
          Evening
        </Toggle>
      </div>
    </div>

    <RichEditor
      :id="`${activePeriod}Thoughts`"
      :key="activePeriod"
      v-model="activeModel"
      :placeholder="activePlaceholder"
      class="h-[30rem]"
    />

    <p v-if="activeError" class="text-sm text-destructive">
      {{ activeError }}
    </p>
  </section>
</template>
