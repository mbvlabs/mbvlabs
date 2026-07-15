<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { Pencil, Trash2, ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import ConfirmAction from '@/Components/ConfirmAction.vue'
import MarkdownPreview from '@/Components/MarkdownPreview.vue'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

defineOptions({ layout: AdminLayout })

interface ProjectInquiry {
  ID: number
  Name: string
  Email: string
  Company: string | null
  Role: string | null
  ProjectType: string | null
  Timeline: string | null
  Message: string
  Source: string | null
  Status: string
  Metadata: string | null
  CreatedAt: string
  UpdatedAt: string
}

const props = defineProps<{
  item: ProjectInquiry
}>()

const form = useForm({})

function destroy() {
  form.delete(`/admin/project-inquiries/${props.item.ID}`)
}
</script>

<template>
  <Head :title="item.Name" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ item.Name }}</h1>
        <p class="text-muted-foreground">{{ item.Email }}</p>
      </div>
      <div class="flex gap-2">
        <Link href="/admin/project-inquiries">
          <Button variant="outline">
            <ArrowLeft class="mr-2 size-4" />
            Back
          </Button>
        </Link>
        <Link :href="`/admin/project-inquiries/${item.ID}/edit`">
          <Button variant="outline">
            <Pencil class="mr-2 size-4" />
            Edit
          </Button>
        </Link>
        <ConfirmAction
          title="Are you sure?"
          :description="`This permanently deletes the inquiry from ${item.Name}. This action cannot be undone.`"
          action-label="Delete inquiry"
          action-variant="destructive"
          confirmation="DELETE"
          :disabled="form.processing"
          @confirm="destroy"
        >
          <Button variant="destructive" :disabled="form.processing">
            <Trash2 class="mr-2 size-4" />
            Delete
          </Button>
        </ConfirmAction>
      </div>
    </div>

    <div class="rounded-none border bg-card">
      <div class="divide-y">
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Name</dt>
          <dd class="col-span-2 text-sm">{{ item.Name }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Email</dt>
          <dd class="col-span-2 text-sm">{{ item.Email }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Company</dt>
          <dd class="col-span-2 text-sm">{{ item.Company || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Role</dt>
          <dd class="col-span-2 text-sm">{{ item.Role || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Project Type</dt>
          <dd class="col-span-2 text-sm">{{ item.ProjectType || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Timeline</dt>
          <dd class="col-span-2 text-sm">{{ item.Timeline || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Message</dt>
          <dd class="col-span-2 text-sm">
            <MarkdownPreview :source="item.Message" />
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Source</dt>
          <dd class="col-span-2 text-sm">{{ item.Source || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Status</dt>
          <dd class="col-span-2 text-sm">
            <Badge variant="secondary">{{ item.Status }}</Badge>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Metadata</dt>
          <dd class="col-span-2 text-sm">{{ item.Metadata || '—' }}</dd>
        </div>
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
