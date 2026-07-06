<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { Pencil, Trash2, ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import MarkdownPreview from '@/Components/MarkdownPreview.vue'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

defineOptions({ layout: AdminLayout })

interface Project {
  ID: number
  Name: string
  Slug: string
  Tagline: string | null
  Description: string | null
  ProjectType: string | null
  RepositoryUrl: string | null
  LiveUrl: string | null
  ImageUrl: string | null
  Technologies: string | null
  StartedAt: string | null
  LaunchedAt: string | null
  PublishedAt: string | null
  IsFeatured: boolean
  CreatedAt: string
  UpdatedAt: string
}

const props = defineProps<{
  item: Project
}>()

const form = useForm({})

function destroy() {
  if (confirm('Are you sure you want to delete this project?')) {
    form.delete(`/admin/projects/${props.item.ID}`)
  }
}
</script>

<template>
  <Head :title="item.Name" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">{{ item.Name }}</h1>
        <p class="text-muted-foreground">{{ item.Tagline }}</p>
      </div>
      <div class="flex gap-2">
        <Link href="/admin/projects">
          <Button variant="outline">
            <ArrowLeft class="mr-2 size-4" />
            Back
          </Button>
        </Link>
        <Link :href="`/admin/projects/${item.ID}/edit`">
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

    <div v-if="item.ImageUrl" class="overflow-hidden rounded-none border bg-card">
      <img :src="item.ImageUrl" :alt="item.Name" class="w-full object-cover" />
    </div>

    <div class="rounded-none border bg-card">
      <div class="divide-y">
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">ID</dt>
          <dd class="col-span-2 text-sm">{{ item.ID }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Name</dt>
          <dd class="col-span-2 text-sm">{{ item.Name }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Slug</dt>
          <dd class="col-span-2 text-sm">{{ item.Slug }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Tagline</dt>
          <dd class="col-span-2 text-sm">{{ item.Tagline || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Description</dt>
          <dd class="col-span-2 text-sm">
            <MarkdownPreview v-if="item.Description" :source="item.Description" />
            <span v-else>—</span>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Project Type</dt>
          <dd class="col-span-2 text-sm">{{ item.ProjectType || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Repository URL</dt>
          <dd class="col-span-2 text-sm">
            <a v-if="item.RepositoryUrl" :href="item.RepositoryUrl" target="_blank" rel="noreferrer" class="text-primary underline">{{ item.RepositoryUrl }}</a>
            <span v-else>—</span>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Live URL</dt>
          <dd class="col-span-2 text-sm">
            <a v-if="item.LiveUrl" :href="item.LiveUrl" target="_blank" rel="noreferrer" class="text-primary underline">{{ item.LiveUrl }}</a>
            <span v-else>—</span>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Image URL</dt>
          <dd class="col-span-2 text-sm">
            <a v-if="item.ImageUrl" :href="item.ImageUrl" target="_blank" rel="noreferrer" class="text-primary underline">{{ item.ImageUrl }}</a>
            <span v-else>—</span>
          </dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Technologies</dt>
          <dd class="col-span-2 text-sm">{{ item.Technologies || '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Started</dt>
          <dd class="col-span-2 text-sm">{{ item.StartedAt ? new Date(item.StartedAt).toLocaleDateString() : '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Launched</dt>
          <dd class="col-span-2 text-sm">{{ item.LaunchedAt ? new Date(item.LaunchedAt).toLocaleDateString() : '—' }}</dd>
        </div>
        <div class="grid grid-cols-3 gap-4 px-6 py-4">
          <dt class="text-sm font-medium text-muted-foreground">Featured</dt>
          <dd class="col-span-2 text-sm">
            <Badge v-if="item.IsFeatured" variant="default">Featured</Badge>
            <Badge v-else variant="secondary">Not featured</Badge>
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
  </div>
</template>
