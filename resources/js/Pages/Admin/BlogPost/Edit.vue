<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { ref, watch } from 'vue'
import { ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'
import RichEditor from '@/Components/RichEditor.vue'

defineOptions({ layout: AdminLayout })

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Published', value: 'published' },
  { label: 'Archived', value: 'archived' },
]

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
  CreatedAt: string
  UpdatedAt: string
  PublicationSchedule: PublicationSchedule | null
}

const props = defineProps<{
  item: BlogPost
}>()

const form = useForm({
  title: props.item.Title,
  slug: props.item.Slug,
  excerpt: props.item.Excerpt,
  body: props.item.Body,
  status: props.item.Status,
  coverImageUrl: props.item.CoverImageUrl ?? '',
  tags: props.item.Tags ?? '',
})

const slugManuallyEdited = ref(false)

function slugify(text: string): string {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

watch(() => form.title, (val) => {
  if (!slugManuallyEdited.value) {
    form.slug = slugify(val)
  }
})

function onSlugInput() {
  slugManuallyEdited.value = true
}

function submit() {
  form.put(`/admin/blog-posts/${props.item.ID}`)
}
</script>

<template>
  <Head :title="`Edit ${item.Title}`" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Edit Blog Post</h1>
        <p class="text-muted-foreground">{{ item.Title }}</p>
      </div>
      <Link :href="`/admin/blog-posts/${item.ID}`">
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
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Primary details</h2>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-1">
              <Label for="title">Title</Label>
              <Input id="title" v-model="form.title" />
              <p v-if="form.errors.title" class="text-sm text-destructive">{{ form.errors.title }}</p>
            </div>

            <div class="space-y-1">
              <Label for="slug">Slug</Label>
              <Input id="slug" v-model="form.slug" @input="onSlugInput" />
              <p v-if="form.errors.slug" class="text-sm text-destructive">{{ form.errors.slug }}</p>
            </div>
          </div>

          <div class="space-y-1">
            <Label for="excerpt">Excerpt</Label>
            <Textarea id="excerpt" v-model="form.excerpt" />
            <p v-if="form.errors.excerpt" class="text-sm text-destructive">{{ form.errors.excerpt }}</p>
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Publishing</h2>
          <div v-if="item.PublicationSchedule" class="border bg-muted/40 px-4 py-3 text-sm">
            Scheduled for {{ new Date(item.PublicationSchedule.scheduled_at).toLocaleString() }}.
            Changing the status to Published or Archived clears this schedule.
          </div>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-1">
              <Label for="status">Status</Label>
              <Select id="status" v-model="form.status" :options="statusOptions" />
              <p v-if="form.errors.status" class="text-sm text-destructive">{{ form.errors.status }}</p>
            </div>

            <div class="space-y-1 md:col-span-2">
              <Label for="coverImageUrl">Cover Image URL</Label>
              <Input id="coverImageUrl" v-model="form.coverImageUrl" />
              <p v-if="form.errors.coverImageUrl" class="text-sm text-destructive">{{ form.errors.coverImageUrl }}</p>
            </div>

            <div class="space-y-1 md:col-span-2">
              <Label for="tags">Tags</Label>
              <Input id="tags" v-model="form.tags" />
              <p v-if="form.errors.tags" class="text-sm text-destructive">{{ form.errors.tags }}</p>
            </div>
          </div>
        </section>

        <section class="space-y-2">
          <Label for="body">Body</Label>
          <RichEditor id="body" v-model="form.body" placeholder="Write your blog post..." class="min-h-[32rem]" />
          <p v-if="form.errors.body" class="text-sm text-destructive">{{ form.errors.body }}</p>
        </section>
      </div>

      <div class="flex justify-end gap-2 border-t px-6 py-4">
        <Link :href="`/admin/blog-posts/${item.ID}`">
          <Button type="button" variant="outline">Cancel</Button>
        </Link>
        <Button type="submit" :disabled="form.processing">
          {{ form.processing ? 'Updating...' : 'Update Blog Post' }}
        </Button>
      </div>
    </form>
  </div>
</template>
