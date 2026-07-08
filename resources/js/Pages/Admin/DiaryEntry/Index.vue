<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { CalendarDays, Eye, Pencil, Plus, Trash2 } from '@lucide/vue'

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

const form = useForm({})

function destroy(id: number) {
  if (confirm('Are you sure you want to delete this diary entry?')) {
    form.delete(routes.adminDiaryEntryDestroy(id))
  }
}
</script>

<template>
  <Head title="Diary" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
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
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="item in items" :key="item.ID">
            <TableCell class="font-medium">{{ item.EntryDate }}</TableCell>
            <TableCell class="max-w-xs truncate text-muted-foreground">
              {{ item.MorningThoughts || '-' }}
            </TableCell>
            <TableCell class="max-w-xs truncate text-muted-foreground">
              {{ item.EveningThoughts || '-' }}
            </TableCell>
            <TableCell class="text-muted-foreground">{{ new Date(item.UpdatedAt).toLocaleString() }}</TableCell>
            <TableCell class="text-right">
              <div class="flex items-center justify-end gap-1">
                <Link :href="routes.adminDiaryEntryShow(item.ID)">
                  <Button variant="ghost" size="icon">
                    <Eye class="size-4" />
                    <span class="sr-only">View</span>
                  </Button>
                </Link>
                <Link :href="routes.adminDiaryEntryEdit(item.ID)">
                  <Button variant="ghost" size="icon">
                    <Pencil class="size-4" />
                    <span class="sr-only">Edit</span>
                  </Button>
                </Link>
                <Button variant="ghost" size="icon" @click="destroy(item.ID)">
                  <Trash2 class="size-4 text-destructive" />
                  <span class="sr-only">Delete</span>
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableRow v-if="items.length === 0">
            <TableCell colspan="5" class="h-24 text-center text-muted-foreground">
              No diary entries found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
