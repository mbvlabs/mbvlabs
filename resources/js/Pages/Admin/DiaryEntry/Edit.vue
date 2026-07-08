<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import { routes } from '@/routes'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import ThoughtsEditor from './ThoughtsEditor.vue'

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
  focus?: 'morning' | 'evening' | ''
}>()

const form = useForm({
  entryDate: props.item.EntryDate,
  morningThoughts: props.item.MorningThoughts,
  eveningThoughts: props.item.EveningThoughts,
})

function submit() {
  form.put(routes.adminDiaryEntryUpdate(props.item.ID))
}
</script>

<template>
  <Head :title="`Edit Diary Entry ${item.EntryDate}`" />

  <div class="mx-auto w-full min-w-0 max-w-5xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Edit Diary Entry</h1>
        <p class="text-muted-foreground">{{ item.EntryDate }}</p>
      </div>
      <Link :href="routes.adminDiaryEntryShow(item.ID)">
        <Button variant="outline">
          <ArrowLeft class="mr-2 size-4" />
          Cancel
        </Button>
      </Link>
    </div>

    <form @submit.prevent="submit" class="border bg-card">
      <div class="space-y-8 p-6">
        <div
          v-if="Object.keys(form.errors).length > 0"
          class="border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive"
        >
          <p class="font-medium">Please fix the highlighted fields.</p>
        </div>

        <section class="space-y-4">
          <div class="max-w-sm space-y-1">
            <Label for="entryDate">Entry Date</Label>
            <Input id="entryDate" v-model="form.entryDate" type="date" />
            <p v-if="form.errors.entryDate" class="text-sm text-destructive">{{ form.errors.entryDate }}</p>
          </div>
        </section>

        <ThoughtsEditor
          v-model:morning-thoughts="form.morningThoughts"
          v-model:evening-thoughts="form.eveningThoughts"
          :morning-error="form.errors.morningThoughts"
          :evening-error="form.errors.eveningThoughts"
          :initial-period="focus"
        />
      </div>

      <div class="flex justify-end gap-2 border-t px-6 py-4">
        <Link :href="routes.adminDiaryEntryShow(item.ID)">
          <Button type="button" variant="outline">Cancel</Button>
        </Link>
        <Button type="submit" :disabled="form.processing">
          {{ form.processing ? 'Saving...' : 'Save Diary Entry' }}
        </Button>
      </div>
    </form>
  </div>
</template>
