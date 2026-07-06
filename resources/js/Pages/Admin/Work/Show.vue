<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { Pencil, Trash2, ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import MarkdownPreview from '@/Components/MarkdownPreview.vue'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

defineOptions({ layout: AdminLayout })

interface Work {
  ID: number
  Title: string
  Slug: string
  ClientName: string | null
  ClientIndustry: string | null
  Summary: string
  Content: string
  CoverImageUrl: string | null
  Specialisms: string[] | null
  Platforms: string[] | null
  Technologies: string[] | null
  StartedAt: string | null
  CompletedAt: string | null
  Status: string
  PublishedAt: string | null
  IsFeatured: boolean
  CreatedAt: string
  UpdatedAt: string
}

const props = defineProps<{
  item: Work
}>()

const form = useForm({})

function destroy() {
  if (confirm('Are you sure you want to delete this work?')) {
    form.delete(`/admin/works/${props.item.ID}`)
  }
}
</script>

<template>
  <Head :title="item.Title" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ item.Title }}</h1>
        <p class="text-muted-foreground">{{ item.Summary }}</p>
      </div>
      <div class="flex gap-2">
        <Link href="/admin/works">
          <Button variant="outline">
            <ArrowLeft class="mr-2 size-4" />
            Back
          </Button>
        </Link>
        <Link :href="`/admin/works/${item.ID}/edit`">
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

    <div v-if="item.CoverImageUrl" class="overflow-hidden rounded-none border bg-card">
      <img :src="item.CoverImageUrl" :alt="item.Title" class="w-full object-cover" />
    </div>

    <div class="rounded-none border bg-card">
      <div class="divide-y">
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Client</dt>
          <dd class="col-span-2 text-sm">{{ item.ClientName || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Industry</dt>
          <dd class="col-span-2 text-sm">{{ item.ClientIndustry || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Status</dt>
          <dd class="col-span-2 text-sm">
            <Badge v-if="item.Status === 'published'" variant="default">Published</Badge>
            <Badge v-else variant="secondary">Draft</Badge>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Featured</dt>
          <dd class="col-span-2 text-sm">
            <Badge v-if="item.IsFeatured" variant="default">Featured</Badge>
            <Badge v-else variant="secondary">Not featured</Badge>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Specialisms</dt>
          <dd class="col-span-2 text-sm">{{ (item.Specialisms ?? []).join(', ') || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Platforms</dt>
          <dd class="col-span-2 text-sm">{{ (item.Platforms ?? []).join(', ') || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Technologies</dt>
          <dd class="col-span-2 text-sm">{{ (item.Technologies ?? []).join(', ') || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Started</dt>
          <dd class="col-span-2 text-sm">{{ item.StartedAt ? new Date(item.StartedAt).toLocaleDateString() : '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Completed</dt>
          <dd class="col-span-2 text-sm">{{ item.CompletedAt ? new Date(item.CompletedAt).toLocaleDateString() : '—' }}</dd>
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

    <div v-if="item.Content" class="rounded-none border bg-card">
      <MarkdownPreview :source="item.Content" class="p-6" />
    </div>
  </div>
</template>
