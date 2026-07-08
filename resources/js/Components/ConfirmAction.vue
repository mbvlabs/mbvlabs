<script setup lang="ts">
import {
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogOverlay,
  AlertDialogPortal,
  AlertDialogRoot,
  AlertDialogTitle,
  AlertDialogTrigger,
} from 'reka-ui'

import { Button } from '@/components/ui/button'

type ButtonVariant = 'default' | 'outline' | 'secondary' | 'ghost' | 'destructive' | 'link'

withDefaults(defineProps<{
  title: string
  description: string
  actionLabel?: string
  actionVariant?: ButtonVariant
  cancelLabel?: string
  disabled?: boolean
}>(), {
  actionLabel: 'Confirm',
  actionVariant: 'default',
  cancelLabel: 'Cancel',
  disabled: false,
})

const emit = defineEmits<{
  confirm: []
}>()
</script>

<template>
  <AlertDialogRoot>
    <AlertDialogTrigger as-child :disabled="disabled">
      <slot />
    </AlertDialogTrigger>

    <AlertDialogPortal>
      <AlertDialogOverlay
        class="fixed inset-0 z-50 bg-black/20 duration-100 supports-backdrop-filter:backdrop-blur-sm data-open:animate-in data-open:fade-in-0 data-closed:animate-out data-closed:fade-out-0"
      />
      <AlertDialogContent
        class="fixed left-1/2 top-1/2 z-50 grid w-[calc(100%-2rem)] max-w-md -translate-x-1/2 -translate-y-1/2 gap-4 border bg-popover p-6 text-popover-foreground shadow-md duration-200 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95"
      >
        <div class="space-y-2">
          <AlertDialogTitle class="text-base font-semibold tracking-tight">
            {{ title }}
          </AlertDialogTitle>
          <AlertDialogDescription class="text-sm text-muted-foreground">
            {{ description }}
          </AlertDialogDescription>
        </div>

        <div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <AlertDialogCancel as-child>
            <Button type="button" variant="outline">
              {{ cancelLabel }}
            </Button>
          </AlertDialogCancel>
          <AlertDialogAction as-child>
            <Button type="button" :variant="actionVariant" @click="emit('confirm')">
              {{ actionLabel }}
            </Button>
          </AlertDialogAction>
        </div>
      </AlertDialogContent>
    </AlertDialogPortal>
  </AlertDialogRoot>
</template>
