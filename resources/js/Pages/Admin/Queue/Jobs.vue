<script setup lang="ts">
import { Head, Link, router, useForm } from '@inertiajs/vue3'
import { Eye, RotateCcw, Search, Trash2, XCircle } from '@lucide/vue'
import { useIntervalFn } from '@vueuse/core'
import { reactive } from 'vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import ConfirmAction from '@/Components/ConfirmAction.vue'
import { routes } from '@/routes'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

defineOptions({ layout: AdminLayout })

interface RiverJob {
  ID: number
  State: string
  Attempt: number
  MaxAttempts: number
  CreatedAt: string
  ScheduledAt: string
  Kind: string
  Queue: string
}

interface RiverQueue {
  Name: string
}

interface Filters {
  State: string
  Queue: string
  Kind: string
}

interface Pagination {
  Page: number
  PageSize: number
  TotalCount: number
  TotalPages: number
}

const props = defineProps<{
  items: RiverJob[]
  queues: RiverQueue[]
  filters: Filters
  pagination: Pagination
}>()

const stateOptions = [
  'available',
  'cancelled',
  'completed',
  'discarded',
  'pending',
  'retryable',
  'running',
  'scheduled',
]

const filterForm = reactive({
  state: props.filters.State,
  queue: props.filters.Queue,
  kind: props.filters.Kind,
})
const actionForm = useForm({})
const pollIntervalMs = 5000

useIntervalFn(() => {
  if (actionForm.processing) {
    return
  }

  router.reload({
    only: ['items', 'queues', 'pagination'],
  })
}, pollIntervalMs)

function applyFilters() {
  router.get(routes.adminQueueJobs(), filterForm, {
    preserveScroll: true,
    preserveState: true,
  })
}

function clearFilters() {
  filterForm.state = ''
  filterForm.queue = ''
  filterForm.kind = ''
  applyFilters()
}

function pageHref(page: number): string {
  const params = new URLSearchParams()
  if (filterForm.state) {
    params.set('state', filterForm.state)
  }
  if (filterForm.queue) {
    params.set('queue', filterForm.queue)
  }
  if (filterForm.kind) {
    params.set('kind', filterForm.kind)
  }
  params.set('page', String(page))
  params.set('per_page', String(props.pagination.PageSize))
  return `${routes.adminQueueJobs()}?${params.toString()}`
}

function cancelJob(id: number) {
  actionForm.post(routes.adminQueueJobCancel(id), { preserveScroll: true })
}

function retryJob(id: number) {
  actionForm.post(routes.adminQueueJobRetry(id), { preserveScroll: true })
}

function discardJob(id: number) {
  actionForm.post(routes.adminQueueJobDiscard(id), { preserveScroll: true })
}

function formatDate(value: string): string {
  return new Date(value).toLocaleString()
}
</script>

<template>
  <Head title="Queue Jobs" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Queue Jobs</h1>
        <p class="text-muted-foreground">Inspect and manage River jobs.</p>
      </div>
    </div>

    <form class="grid gap-3 border bg-card p-4 md:grid-cols-[1fr_1fr_2fr_auto_auto]" @submit.prevent="applyFilters">
      <select v-model="filterForm.state" class="h-10 border bg-background px-3 text-sm">
        <option value="">All states</option>
        <option v-for="state in stateOptions" :key="state" :value="state">{{ state }}</option>
      </select>
      <select v-model="filterForm.queue" class="h-10 border bg-background px-3 text-sm">
        <option value="">All queues</option>
        <option v-for="queue in queues" :key="queue.Name" :value="queue.Name">{{ queue.Name }}</option>
      </select>
      <Input v-model="filterForm.kind" placeholder="Filter by kind" />
      <Button type="submit" variant="outline">
        <Search class="mr-2 size-4" />
        Filter
      </Button>
      <Button type="button" variant="ghost" @click="clearFilters">
        Clear
      </Button>
    </form>

    <div class="rounded-none border bg-card">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Kind</TableHead>
            <TableHead>Queue</TableHead>
            <TableHead>State</TableHead>
            <TableHead>Attempts</TableHead>
            <TableHead>Created</TableHead>
            <TableHead>Scheduled</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="item in items" :key="item.ID">
            <TableCell class="font-medium">#{{ item.ID }}</TableCell>
            <TableCell class="max-w-xs truncate">{{ item.Kind }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.Queue }}</TableCell>
            <TableCell><Badge variant="secondary">{{ item.State }}</Badge></TableCell>
            <TableCell class="text-muted-foreground">{{ item.Attempt }} / {{ item.MaxAttempts }}</TableCell>
            <TableCell class="text-muted-foreground">{{ formatDate(item.CreatedAt) }}</TableCell>
            <TableCell class="text-muted-foreground">{{ formatDate(item.ScheduledAt) }}</TableCell>
            <TableCell class="text-right">
              <div class="flex items-center justify-end gap-1">
                <Link :href="routes.adminQueueJobShow(item.ID)">
                  <Button variant="ghost" size="icon">
                    <Eye class="size-4" />
                    <span class="sr-only">View</span>
                  </Button>
                </Link>
                <ConfirmAction
                  title="Retry job"
                  :description="`Move job ${item.ID} back to available so River can process it again.`"
                  action-label="Retry"
                  action-variant="default"
                  :disabled="actionForm.processing"
                  @confirm="retryJob(item.ID)"
                >
                  <Button variant="ghost" size="icon" :disabled="actionForm.processing">
                    <RotateCcw class="size-4" />
                    <span class="sr-only">Retry</span>
                  </Button>
                </ConfirmAction>
                <ConfirmAction
                  title="Cancel job"
                  :description="`Cancel job ${item.ID}. River will treat it as finalized and it will not run.`"
                  action-label="Cancel job"
                  action-variant="destructive"
                  :disabled="actionForm.processing"
                  @confirm="cancelJob(item.ID)"
                >
                  <Button variant="ghost" size="icon" :disabled="actionForm.processing">
                    <XCircle class="size-4" />
                    <span class="sr-only">Cancel</span>
                  </Button>
                </ConfirmAction>
                <ConfirmAction
                  title="Discard job"
                  :description="`Discard job ${item.ID}. This finalizes the job without running it.`"
                  action-label="Discard"
                  action-variant="destructive"
                  :disabled="actionForm.processing"
                  @confirm="discardJob(item.ID)"
                >
                  <Button variant="ghost" size="icon" :disabled="actionForm.processing">
                    <Trash2 class="size-4 text-destructive" />
                    <span class="sr-only">Discard</span>
                  </Button>
                </ConfirmAction>
              </div>
            </TableCell>
          </TableRow>
          <TableRow v-if="items.length === 0">
            <TableCell colspan="8" class="h-24 text-center text-muted-foreground">
              No jobs found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <p class="text-sm text-muted-foreground">
        Showing page {{ pagination.Page }} of {{ Math.max(pagination.TotalPages, 1) }}
        for {{ pagination.TotalCount }} jobs.
      </p>
      <div class="flex gap-2">
        <Link v-if="pagination.Page > 1" :href="pageHref(pagination.Page - 1)">
          <Button variant="outline" size="sm">Previous</Button>
        </Link>
        <Button v-else variant="outline" size="sm" disabled>Previous</Button>

        <Link v-if="pagination.Page < pagination.TotalPages" :href="pageHref(pagination.Page + 1)">
          <Button variant="outline" size="sm">Next</Button>
        </Link>
        <Button v-else variant="outline" size="sm" disabled>Next</Button>
      </div>
    </div>
  </div>
</template>
