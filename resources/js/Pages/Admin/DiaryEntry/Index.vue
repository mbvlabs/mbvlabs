<script setup lang="ts">
import { Head, Link } from '@inertiajs/vue3'
import { CalendarDays, Plus } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import { routes } from '@/routes'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

defineOptions({ layout: AdminLayout })

interface DiaryEntry {
  ID: number
  EntryDate: string
  MorningThoughts: string
  EveningThoughts: string
  CreatedAt: string
  UpdatedAt: string
}

defineProps<{
  items: DiaryEntry[]
}>()
</script>

<template>
  <Head title="Diary" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Diary</h1>
        <p class="text-muted-foreground">Daily morning and evening idea capture.</p>
      </div>
      <div class="flex gap-2">
        <Link :href="routes.adminDiaryEntryToday()">
          <Button variant="outline">
            <CalendarDays class="mr-2 size-4" />
            Today
          </Button>
        </Link>
        <Link :href="routes.adminDiaryEntryNew()">
          <Button>
            <Plus class="mr-2 size-4" />
            New Entry
          </Button>
        </Link>
      </div>
    </div>

    <div class="rounded-none border bg-card">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Date</TableHead>
            <TableHead>Morning</TableHead>
            <TableHead>Evening</TableHead>
            <TableHead>Updated</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="item in items"
            :key="item.ID"
            :href="routes.adminDiaryEntryShow(item.ID)"
            :aria-label="`View diary entry ${item.EntryDate}`"
          >
            <TableCell class="font-medium">{{ item.EntryDate }}</TableCell>
            <TableCell class="max-w-xs truncate text-muted-foreground">
              {{ item.MorningThoughts || '-' }}
            </TableCell>
            <TableCell class="max-w-xs truncate text-muted-foreground">
              {{ item.EveningThoughts || '-' }}
            </TableCell>
            <TableCell class="text-muted-foreground">{{ new Date(item.UpdatedAt).toLocaleString() }}</TableCell>
          </TableRow>
          <TableRow v-if="items.length === 0">
            <TableCell colspan="4" class="h-24 text-center text-muted-foreground">
              No diary entries found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
