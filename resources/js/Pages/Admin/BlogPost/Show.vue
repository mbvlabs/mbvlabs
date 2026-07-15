<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { Pencil, Trash2, ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import ConfirmAction from '@/Components/ConfirmAction.vue'
import MarkdownPreview from '@/Components/MarkdownPreview.vue'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

defineOptions({ layout: AdminLayout })

interface PublicationSchedule {
  job_id: number
  scheduled_at: string
}

interface BlogPost {
  ID: number
  Title: string
  Slug: string
  Excerpt: string
  Body: string
  Status: string
  CoverImageUrl: string | null
  Tags: string | null
  PublishedAt: string | null
  CreatedAt: string
  UpdatedAt: string
  PublicationSchedule: PublicationSchedule | null
}

const props = defineProps<{
  item: BlogPost
}>()

const form = useForm({})

function destroy() {
  form.delete(`/admin/blog-posts/${props.item.ID}`)
}
</script>

<template>
  <Head :title="item.Title" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ item.Title }}</h1>
        <p class="text-muted-foreground">{{ item.Excerpt }}</p>
      </div>
      <div class="flex gap-2">
        <Link href="/admin/blog-posts">
          <Button variant="outline">
            <ArrowLeft class="mr-2 size-4" />
            Back
          </Button>
        </Link>
        <Link :href="`/admin/blog-posts/${item.ID}/edit`">
          <Button variant="outline">
            <Pencil class="mr-2 size-4" />
            Edit
          </Button>
        </Link>
        <ConfirmAction
          title="Are you sure?"
          :description="`This permanently deletes ${item.Title}. This action cannot be undone.`"
          action-label="Delete blog post"
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

    <div v-if="item.CoverImageUrl" class="overflow-hidden rounded-none border bg-card">
      <img :src="item.CoverImageUrl" :alt="item.Title" class="w-full object-cover" />
    </div>

    <div class="rounded-none border bg-card">
      <div class="divide-y">
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Slug</dt>
          <dd class="col-span-2 text-sm">{{ item.Slug }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Status</dt>
          <dd class="col-span-2 text-sm">
            <Badge variant="secondary">
              {{ item.Status === 'draft' && item.PublicationSchedule ? 'Scheduled' : item.Status }}
            </Badge>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Tags</dt>
          <dd class="col-span-2 text-sm">{{ item.Tags || '—' }}</dd>
        </div>
        <div v-if="item.PublicationSchedule" class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Scheduled for</dt>
          <dd class="col-span-2 text-sm">
            {{ new Date(item.PublicationSchedule.scheduled_at).toLocaleString() }}
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Published</dt>
          <dd class="col-span-2 text-sm">
            {{ item.PublishedAt ? new Date(item.PublishedAt).toLocaleDateString() : 'Not published' }}
          </dd>
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

    <div v-if="item.Body" class="rounded-none border bg-card">
      <MarkdownPreview :source="item.Body" class="p-6" />
    </div>
  </div>
</template>
