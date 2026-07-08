<script setup lang="ts">
import { Head, Link, router, useForm } from '@inertiajs/vue3'
import { ArrowLeft, LoaderCircle, RotateCcw, Trash2, XCircle } from '@lucide/vue'
import { useIntervalFn } from '@vueuse/core'
import { computed, ref } from 'vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import ConfirmAction from '@/Components/ConfirmAction.vue'
import JsonHighlight from '@/Components/JsonHighlight.vue'
import { routes } from '@/routes'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Toggle } from '@/components/ui/toggle'

defineOptions({ layout: AdminLayout })

interface RiverJob {
  ID: number
  State: string
  Attempt: number
  MaxAttempts: number
  AttemptedAt: string | null
  CreatedAt: string
  FinalizedAt: string | null
  ScheduledAt: string
  Priority: number
  Args: string
  AttemptedBy: string[] | null
  Errors: string[] | null
  Kind: string
  Metadata: string
  Queue: string
  Tags: string[] | null
}

const props = defineProps<{
  item: RiverJob
}>()

const form = useForm({})
const pollIntervalMs = 5000
const isPollingEnabled = ref(true)

interface ParsedRiverError {
  raw: string
  displayNumber: number
  attempt: number | null
  at: string | null
}

const jobErrors = computed(() => props.item.Errors ?? [])
const attemptedBy = computed(() => props.item.AttemptedBy ?? [])
const tags = computed(() => props.item.Tags ?? [])

const sortedErrors = computed<ParsedRiverError[]>(() => {
  return jobErrors.value
    .map((raw, index) => {
      const parsed = parseRiverError(raw)
      return {
        raw,
        displayNumber: index + 1,
        attempt: parsed.attempt,
        at: parsed.at,
      }
    })
    .sort((a, b) => {
      if (a.at && b.at) {
        return new Date(b.at).getTime() - new Date(a.at).getTime()
      }
      return b.displayNumber - a.displayNumber
    })
})

useIntervalFn(() => {
  if (!isPollingEnabled.value) {
    return
  }

  if (form.processing) {
    return
  }

  router.reload({
    only: ['item'],
  })
}, pollIntervalMs)

function cancelJob() {
  form.post(routes.adminQueueJobCancel(props.item.ID), { preserveScroll: true })
}

function retryJob() {
  form.post(routes.adminQueueJobRetry(props.item.ID), { preserveScroll: true })
}

function discardJob() {
  form.post(routes.adminQueueJobDiscard(props.item.ID), { preserveScroll: true })
}

function formatDate(value: string | null): string {
  if (!value) {
    return '-'
  }
  return new Date(value).toLocaleString()
}

function parseRiverError(value: string): { attempt: number | null, at: string | null } {
  try {
    const parsed = JSON.parse(value) as { attempt?: unknown, at?: unknown }
    return {
      attempt: typeof parsed.attempt === 'number' ? parsed.attempt : null,
      at: typeof parsed.at === 'string' ? parsed.at : null,
    }
  } catch {
    return { attempt: null, at: null }
  }
}
</script>

<template>
  <Head :title="`Queue Job ${item.ID}`" />

  <div class="mx-auto w-full min-w-0 max-w-5xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Job #{{ item.ID }}</h1>
        <p class="text-muted-foreground">{{ item.Kind }}</p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <div class="flex items-center gap-1.5 px-2 text-xs font-medium uppercase tracking-widest text-muted-foreground">
          <LoaderCircle
            :class="['size-3', isPollingEnabled ? 'animate-spin' : '']"
          />
          <span>{{ isPollingEnabled ? 'Polling' : 'Paused' }}</span>
        </div>
        <Toggle
          v-model="isPollingEnabled"
          variant="outline"
          size="sm"
          :aria-label="isPollingEnabled ? 'Pause job polling' : 'Resume job polling'"
        >
          {{ isPollingEnabled ? 'Pause polling' : 'Resume polling' }}
        </Toggle>
      </div>
    </div>

    <div class="flex flex-wrap gap-2 border-b pb-4">
      <Link :href="routes.adminQueueJobs()">
        <Button variant="outline">
          <ArrowLeft class="mr-2 size-4" />
          Back
        </Button>
      </Link>
      <ConfirmAction
        title="Retry job"
        :description="`Move job ${item.ID} back to available so River can process it again.`"
        action-label="Retry"
        :disabled="form.processing"
        @confirm="retryJob"
      >
        <Button variant="outline" :disabled="form.processing">
          <RotateCcw class="mr-2 size-4" />
          Retry
        </Button>
      </ConfirmAction>
      <ConfirmAction
        title="Cancel job"
        :description="`Cancel job ${item.ID}. River will treat it as finalized and it will not run.`"
        action-label="Cancel job"
        action-variant="destructive"
        :disabled="form.processing"
        @confirm="cancelJob"
      >
        <Button variant="outline" :disabled="form.processing">
          <XCircle class="mr-2 size-4" />
          Cancel
        </Button>
      </ConfirmAction>
      <ConfirmAction
        title="Discard job"
        :description="`Discard job ${item.ID}. This finalizes the job without running it.`"
        action-label="Discard"
        action-variant="destructive"
        :disabled="form.processing"
        @confirm="discardJob"
      >
        <Button variant="outline" :disabled="form.processing">
          <Trash2 class="mr-2 size-4 text-destructive" />
          Discard
        </Button>
      </ConfirmAction>
    </div>

    <section class="border bg-card">
      <div class="grid gap-0 divide-y md:grid-cols-2 md:divide-x md:divide-y-0">
        <div class="space-y-4 p-6">
          <div>
            <dt class="text-sm font-medium text-muted-foreground">State</dt>
            <dd class="mt-1"><Badge variant="secondary">{{ item.State }}</Badge></dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-muted-foreground">Queue</dt>
            <dd class="mt-1 text-sm">{{ item.Queue }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-muted-foreground">Attempts</dt>
            <dd class="mt-1 text-sm">{{ item.Attempt }} / {{ item.MaxAttempts }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-muted-foreground">Priority</dt>
            <dd class="mt-1 text-sm">{{ item.Priority }}</dd>
          </div>
        </div>
        <div class="space-y-4 p-6">
          <div>
            <dt class="text-sm font-medium text-muted-foreground">Created</dt>
            <dd class="mt-1 text-sm">{{ formatDate(item.CreatedAt) }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-muted-foreground">Scheduled</dt>
            <dd class="mt-1 text-sm">{{ formatDate(item.ScheduledAt) }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-muted-foreground">Attempted</dt>
            <dd class="mt-1 text-sm">{{ formatDate(item.AttemptedAt) }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-muted-foreground">Finalized</dt>
            <dd class="mt-1 text-sm">{{ formatDate(item.FinalizedAt) }}</dd>
          </div>
        </div>
      </div>
    </section>

    <section class="space-y-6">
      <div class="border bg-card p-6">
        <h2 class="font-semibold leading-none tracking-tight">Args</h2>
        <JsonHighlight :value="item.Args" class="mt-4 max-h-96" />
      </div>
      <div class="border bg-card p-6">
        <h2 class="font-semibold leading-none tracking-tight">Metadata</h2>
        <JsonHighlight :value="item.Metadata" class="mt-4 max-h-96" />
      </div>
    </section>

    <section class="border bg-card p-6">
      <h2 class="font-semibold leading-none tracking-tight">Errors</h2>
      <div v-if="sortedErrors.length > 0" class="mt-4 space-y-3">
        <div
          v-for="error in sortedErrors"
          :key="error.displayNumber"
          class="border bg-muted/30"
        >
          <div class="flex flex-wrap items-center gap-2 border-b px-3 py-2 text-xs text-muted-foreground">
            <span class="font-semibold text-foreground">Error {{ error.displayNumber }}</span>
            <span v-if="error.at">{{ formatDate(error.at) }}</span>
            <span v-if="error.attempt !== null">Attempt {{ error.attempt }}</span>
          </div>
          <JsonHighlight :value="error.raw" class="max-h-80 p-3" />
        </div>
      </div>
      <p v-else class="mt-4 text-sm text-muted-foreground">No recorded errors.</p>
    </section>

    <section class="border bg-card p-6">
      <h2 class="font-semibold leading-none tracking-tight">Worker Data</h2>
      <dl class="mt-4 grid gap-4 md:grid-cols-2">
        <div>
          <dt class="text-sm font-medium text-muted-foreground">Attempted By</dt>
          <dd class="mt-1 text-sm">{{ attemptedBy.length ? attemptedBy.join(', ') : '-' }}</dd>
        </div>
        <div>
          <dt class="text-sm font-medium text-muted-foreground">Tags</dt>
          <dd class="mt-1 text-sm">{{ tags.length ? tags.join(', ') : '-' }}</dd>
        </div>
      </dl>
    </section>
  </div>
</template>
