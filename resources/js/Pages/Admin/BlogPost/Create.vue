<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { watch } from 'vue'
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

const form = useForm({
  title: '',
  slug: '',
  excerpt: '',
  body: '',
  status: 'draft',
  coverImageUrl: '',
  tags: '',
})

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
  form.slug = slugify(val)
})

function submit() {
  form.post('/admin/blog-posts')
}
</script>

<template>
  <Head title="New Blog Post" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">New Blog Post</h1>
        <p class="text-muted-foreground">Write a new blog post.</p>
      </div>
      <Link href="/admin/blog-posts">
        <Button variant="outline">
          <ArrowLeft class="mr-2 size-4" />
          Back
        </Button>
      </Link>
    </div>

    <form @submit.prevent="submit" class="border bg-card">
      <div class="space-y-8 p-6">
        <section class="space-y-4">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Primary details</h2>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-1">
              <Label for="title">Title</Label>
              <Input id="title" v-model="form.title" />
            </div>

            <div class="space-y-1">
              <Label for="slug">Slug</Label>
              <Input id="slug" v-model="form.slug" />
            </div>
          </div>

          <div class="space-y-1">
            <Label for="excerpt">Excerpt</Label>
            <Textarea id="excerpt" v-model="form.excerpt" />
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Publishing</h2>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-1">
              <Label for="status">Status</Label>
              <Select id="status" v-model="form.status" :options="statusOptions" />
            </div>

            <div class="space-y-1 md:col-span-2">
              <Label for="coverImageUrl">Cover Image URL</Label>
              <Input id="coverImageUrl" v-model="form.coverImageUrl" />
            </div>

            <div class="space-y-1 md:col-span-2">
              <Label for="tags">Tags</Label>
              <Input id="tags" v-model="form.tags" placeholder="strategy, systems, launch" />
            </div>
          </div>
        </section>

        <section class="space-y-2">
          <Label for="body">Body</Label>
          <RichEditor id="body" v-model="form.body" placeholder="Write your blog post..." class="min-h-[32rem]" />
        </section>
      </div>

      <div class="flex justify-end gap-2 border-t px-6 py-4">
        <Link href="/admin/blog-posts">
          <Button type="button" variant="outline">Cancel</Button>
        </Link>
        <Button type="submit" :disabled="form.processing">
          {{ form.processing ? 'Creating...' : 'Create Blog Post' }}
        </Button>
      </div>
    </form>
  </div>
</template>
