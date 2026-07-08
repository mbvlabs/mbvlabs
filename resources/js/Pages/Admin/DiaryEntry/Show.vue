<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { ArrowLeft, Pencil, Trash2 } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import { routes } from '@/routes'
import { Button } from '@/components/ui/button'

defineOptions({ layout: AdminLayout })

interface DiaryEntry {
  ID: number
  EntryDate: string
  MorningThoughts: string
  EveningThoughts: string
  CreatedAt: string
  UpdatedAt: string
}

const props = defineProps<{
  item: DiaryEntry
}>()

const form = useForm({})

function destroy() {
  if (confirm('Are you sure you want to delete this diary entry?')) {
    form.delete(routes.adminDiaryEntryDestroy(props.item.ID))
  }
}
</script>

<template>
  <Head :title="`Diary Entry ${item.EntryDate}`" />

  <div class="mx-auto w-full min-w-0 max-w-5xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ item.EntryDate }}</h1>
        <p class="text-muted-foreground">Diary entry</p>
      </div>
      <div class="flex gap-2">
        <Link :href="routes.adminDiaryEntryIndex()">
          <Button variant="outline">
            <ArrowLeft class="mr-2 size-4" />
            Back
          </Button>
        </Link>
        <Link :href="routes.adminDiaryEntryEdit(item.ID)">
          <Button variant="outline">
            <Pencil class="mr-2 size-4" />
            Edit
          </Button>
        </Link>
        <Button variant="outline" @click="destroy">
          <Trash2 class="mr-2 size-4 text-destructive" />
          Delete
        </Button>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <section class="border bg-card p-6">
        <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Morning</h2>
        <p class="mt-4 whitespace-pre-wrap text-sm leading-6">{{ item.MorningThoughts || '-' }}</p>
      </section>

      <section class="border bg-card p-6">
        <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Evening</h2>
        <p class="mt-4 whitespace-pre-wrap text-sm leading-6">{{ item.EveningThoughts || '-' }}</p>
      </section>
    </div>

    <div class="border bg-card">
      <div class="divide-y">
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Created</dt>
          <dd class="col-span-2 text-sm">{{ new Date(item.CreatedAt).toLocaleString() }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Updated</dt>
          <dd class="col-span-2 text-sm">{{ new Date(item.UpdatedAt).toLocaleString() }}</dd>
        </div>
      </div>
    </div>
  </div>
</template>
