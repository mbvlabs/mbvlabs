<script setup lang="ts">
import { Head, Link, router, useForm } from '@inertiajs/vue3'
import { CirclePause, CirclePlay, LoaderCircle } from '@lucide/vue'
import { useIntervalFn } from '@vueuse/core'
import { computed, ref } from 'vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import { routes } from '@/routes'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Toggle } from '@/components/ui/toggle'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

defineOptions({ layout: AdminLayout })

interface RiverQueue {
  Name: string
  CreatedAt: string | null
  UpdatedAt: string | null
  PausedAt: string | null
  IsPaused: boolean
  AvailableCount: number
  CancelledCount: number
  CompletedCount: number
  DiscardedCount: number
  PendingCount: number
  RetryableCount: number
  RunningCount: number
  ScheduledCount: number
  TotalCount: number
  ActiveClients: number
  MaxWorkers: number
  NumJobsRunning: number
  NumJobsCompleted: number
}

interface RiverJob {
  ID: number
  State: string
  Kind: string
  Queue: string
  Attempt: number
  MaxAttempts: number
  CreatedAt: string
  ScheduledAt: string
}

interface StateCount {
  State: string
  Count: number
}

const props = defineProps<{
  queues: RiverQueue[]
  stateCounts: StateCount[]
  recentJobs: RiverJob[]
}>()

const form = useForm({ name: '' })
const pollIntervalMs = 5000
const isPolling = ref(false)
const isPollingEnabled = ref(true)

const totalJobs = computed(() => props.stateCounts.reduce((sum, item) => sum + item.Count, 0))
const availableJobs = computed(() => countForState('available') + countForState('scheduled') + countForState('pending'))
const runningJobs = computed(() => countForState('running'))
const retryableJobs = computed(() => countForState('retryable') + countForState('discarded'))

useIntervalFn(() => {
  if (!isPollingEnabled.value) {
    return
  }

  if (form.processing) {
    return
  }

  router.reload({
    only: ['queues', 'stateCounts', 'recentJobs'],
    onStart: () => {
      isPolling.value = true
    },
    onFinish: () => {
      isPolling.value = false
    },
  })
}, pollIntervalMs)

function countForState(state: string): number {
  return props.stateCounts.find((item) => item.State === state)?.Count || 0
}

function pauseQueue(name: string) {
  form.name = name
  form.post(routes.adminQueuePause(), { preserveScroll: true })
}

function resumeQueue(name: string) {
  form.name = name
  form.post(routes.adminQueueResume(), { preserveScroll: true })
}

function queueJobsHref(name: string): string {
  return `${routes.adminQueueJobs()}?queue=${encodeURIComponent(name)}`
}

function goToRecentJob(id: number) {
  router.visit(routes.adminQueueJobShow(id))
}

function isJobInProgress(job: RiverJob): boolean {
  return job.State === 'running'
}

function formatDate(value: string | null): string {
  if (!value) {
    return '-'
  }
  return new Date(value).toLocaleString()
}
</script>

<template>
  <Head title="Queue" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Queue</h1>
        <p class="text-muted-foreground">River queue health, workers, and recent job activity.</p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <div class="flex items-center gap-1.5 text-xs font-medium uppercase tracking-widest text-muted-foreground">
          <LoaderCircle
            :class="['size-3', isPollingEnabled ? 'animate-spin' : '']"
          />
          <span>{{ isPollingEnabled ? 'Polling' : 'Paused' }}</span>
        </div>
        <Toggle
          v-model="isPollingEnabled"
          variant="outline"
          size="sm"
          :aria-label="isPollingEnabled ? 'Pause queue polling' : 'Resume queue polling'"
        >
          {{ isPollingEnabled ? 'Pause polling' : 'Resume polling' }}
        </Toggle>
      </div>
    </div>

    <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <section class="border bg-card p-5">
        <div>
          <span class="text-sm font-medium text-muted-foreground">Total Jobs</span>
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ totalJobs }}</div>
        <p class="mt-1 text-xs text-muted-foreground">All recorded River jobs</p>
      </section>

      <section class="border bg-card p-5">
        <div>
          <span class="text-sm font-medium text-muted-foreground">Queued</span>
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ availableJobs }}</div>
        <p class="mt-1 text-xs text-muted-foreground">Available, pending, or scheduled</p>
      </section>

      <section class="border bg-card p-5">
        <div>
          <span class="text-sm font-medium text-muted-foreground">Running</span>
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ runningJobs }}</div>
        <p class="mt-1 text-xs text-muted-foreground">Currently claimed by workers</p>
      </section>

      <section class="border bg-card p-5">
        <div>
          <span class="text-sm font-medium text-muted-foreground">Needs Review</span>
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ retryableJobs }}</div>
        <p class="mt-1 text-xs text-muted-foreground">Retryable or discarded jobs</p>
      </section>
    </div>

    <section class="border bg-card">
      <div class="flex items-center justify-between border-b p-5">
        <div>
          <h2 class="font-semibold leading-none tracking-tight">Queues</h2>
          <p class="mt-1 text-sm text-muted-foreground">Pause, resume, and inspect queue load.</p>
        </div>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Available</TableHead>
            <TableHead>Running</TableHead>
            <TableHead>Retryable</TableHead>
            <TableHead>Workers</TableHead>
            <TableHead>Updated</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="queue in queues" :key="queue.Name">
            <TableCell>
              <Link :href="queueJobsHref(queue.Name)" class="font-medium hover:underline">
                {{ queue.Name }}
              </Link>
              <p class="text-xs text-muted-foreground">{{ queue.TotalCount }} total jobs</p>
            </TableCell>
            <TableCell>
              <Badge :variant="queue.IsPaused ? 'destructive' : 'secondary'">
                {{ queue.IsPaused ? 'paused' : 'active' }}
              </Badge>
            </TableCell>
            <TableCell>{{ queue.AvailableCount + queue.PendingCount + queue.ScheduledCount }}</TableCell>
            <TableCell>{{ queue.RunningCount }}</TableCell>
            <TableCell>{{ queue.RetryableCount }}</TableCell>
            <TableCell class="text-muted-foreground">
              {{ queue.ActiveClients }} clients / {{ queue.MaxWorkers }} max
            </TableCell>
            <TableCell class="text-muted-foreground">{{ formatDate(queue.UpdatedAt) }}</TableCell>
            <TableCell class="text-right">
              <Button
                v-if="queue.IsPaused"
                variant="ghost"
                size="icon"
                @click="resumeQueue(queue.Name)"
              >
                <CirclePlay class="size-4" />
                <span class="sr-only">Resume queue</span>
              </Button>
              <Button
                v-else
                variant="ghost"
                size="icon"
                @click="pauseQueue(queue.Name)"
              >
                <CirclePause class="size-4" />
                <span class="sr-only">Pause queue</span>
              </Button>
            </TableCell>
          </TableRow>
          <TableRow v-if="queues.length === 0">
            <TableCell colspan="8" class="h-24 text-center text-muted-foreground">
              No queues found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </section>

    <section class="border bg-card">
      <div class="border-b p-5">
        <div>
          <div class="flex flex-wrap items-center gap-3">
            <h2 class="font-semibold leading-none tracking-tight">Recent Jobs</h2>
          </div>
          <p class="mt-1 text-sm text-muted-foreground">Latest jobs recorded by River.</p>
        </div>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Kind</TableHead>
            <TableHead>Queue</TableHead>
            <TableHead>State</TableHead>
            <TableHead>Attempts</TableHead>
            <TableHead>Scheduled</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="job in recentJobs"
            :key="job.ID"
            class="cursor-pointer"
            role="link"
            tabindex="0"
            @click="goToRecentJob(job.ID)"
            @keydown.enter.prevent="goToRecentJob(job.ID)"
            @keydown.space.prevent="goToRecentJob(job.ID)"
          >
            <TableCell>
              <span class="font-medium">#{{ job.ID }}</span>
            </TableCell>
            <TableCell class="max-w-sm truncate">{{ job.Kind }}</TableCell>
            <TableCell class="text-muted-foreground">{{ job.Queue }}</TableCell>
            <TableCell>
              <div class="flex items-center gap-2">
                <LoaderCircle
                  v-if="isJobInProgress(job)"
                  class="size-3 animate-spin text-muted-foreground"
                />
                <Badge variant="secondary">{{ job.State }}</Badge>
              </div>
            </TableCell>
            <TableCell class="text-muted-foreground">{{ job.Attempt }} / {{ job.MaxAttempts }}</TableCell>
            <TableCell class="text-muted-foreground">{{ formatDate(job.ScheduledAt) }}</TableCell>
          </TableRow>
          <TableRow v-if="recentJobs.length === 0">
            <TableCell colspan="6" class="h-20 text-center text-muted-foreground">
              No jobs found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </section>
  </div>
</template>
