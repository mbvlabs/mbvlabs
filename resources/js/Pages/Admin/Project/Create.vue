<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { watch } from 'vue'
import { ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import RichEditor from '@/Components/RichEditor.vue'

import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'

defineOptions({ layout: AdminLayout })

const projectTypeOptions = [
  { label: 'Web app', value: 'web-app' },
  { label: 'Website', value: 'website' },
  { label: 'Automation', value: 'automation' },
  { label: 'Integration', value: 'integration' },
  { label: 'Consulting', value: 'consulting' },
  { label: 'Framework', value: 'framework' },
]

const form = useForm({
  name: '',
  slug: '',
  tagline: '',
  description: '',
  projectType: 'web-app',
  repositoryUrl: '',
  liveUrl: '',
  imageUrl: '',
  technologies: '',
  startedAt: '',
  launchedAt: '',
  publishedAt: '',
  isFeatured: false,
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

watch(() => form.name, (val) => {
  form.slug = slugify(val)
})

function setFeatured(checked: boolean | 'indeterminate') {
  form.isFeatured = checked === true
}

function submit() {
  form.post('/admin/projects')
}
</script>

<template>
  <Head title="New Project" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">New Project</h1>
        <p class="text-muted-foreground">Add a new project to your portfolio.</p>
      </div>
      <Link href="/admin/projects">
        <Button variant="outline">
          <ArrowLeft class="mr-2 size-4" />
          Back
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
              <Label for="name">Name</Label>
              <Input id="name" v-model="form.name" />
              <p v-if="form.errors.name" class="text-sm text-destructive">{{ form.errors.name }}</p>
            </div>

            <div class="space-y-1">
              <Label for="slug">Slug</Label>
              <Input id="slug" v-model="form.slug" />
              <p v-if="form.errors.slug" class="text-sm text-destructive">{{ form.errors.slug }}</p>
            </div>

            <div class="space-y-1">
              <Label for="tagline">Tagline</Label>
              <Input id="tagline" v-model="form.tagline" />
              <p v-if="form.errors.tagline" class="text-sm text-destructive">{{ form.errors.tagline }}</p>
            </div>

            <div class="space-y-1">
              <Label for="projectType">Project Type</Label>
              <Select id="projectType" v-model="form.projectType" :options="projectTypeOptions" />
              <p v-if="form.errors.projectType" class="text-sm text-destructive">{{ form.errors.projectType }}</p>
            </div>
          </div>

          <div class="space-y-1">
            <Label for="description">Description</Label>
            <RichEditor id="description" v-model="form.description" placeholder="Write the project description..." class="min-h-[24rem]" />
            <p v-if="form.errors.description" class="text-sm text-destructive">{{ form.errors.description }}</p>
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Links and media</h2>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-1">
              <Label for="repositoryUrl">Repository URL</Label>
              <Input id="repositoryUrl" v-model="form.repositoryUrl" />
              <p v-if="form.errors.repositoryUrl" class="text-sm text-destructive">{{ form.errors.repositoryUrl }}</p>
            </div>

            <div class="space-y-1">
              <Label for="liveUrl">Live URL</Label>
              <Input id="liveUrl" v-model="form.liveUrl" />
              <p v-if="form.errors.liveUrl" class="text-sm text-destructive">{{ form.errors.liveUrl }}</p>
            </div>

            <div class="space-y-1 md:col-span-2">
              <Label for="imageUrl">Image URL</Label>
              <Input id="imageUrl" v-model="form.imageUrl" />
              <p v-if="form.errors.imageUrl" class="text-sm text-destructive">{{ form.errors.imageUrl }}</p>
            </div>

            <div class="space-y-1 md:col-span-2">
              <Label for="technologies">Technologies</Label>
              <Input id="technologies" v-model="form.technologies" placeholder="Go, Vue, PostgreSQL" />
              <p v-if="form.errors.technologies" class="text-sm text-destructive">{{ form.errors.technologies }}</p>
            </div>
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Publishing</h2>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
            <div class="space-y-1">
              <Label for="startedAt">Started Date</Label>
              <Input id="startedAt" type="date" v-model="form.startedAt" />
              <p v-if="form.errors.startedAt" class="text-sm text-destructive">{{ form.errors.startedAt }}</p>
            </div>

            <div class="space-y-1">
              <Label for="launchedAt">Launched Date</Label>
              <Input id="launchedAt" type="date" v-model="form.launchedAt" />
              <p v-if="form.errors.launchedAt" class="text-sm text-destructive">{{ form.errors.launchedAt }}</p>
            </div>

            <div class="space-y-1">
              <Label for="publishedAt">Published Date</Label>
              <Input id="publishedAt" type="date" v-model="form.publishedAt" />
              <p v-if="form.errors.publishedAt" class="text-sm text-destructive">{{ form.errors.publishedAt }}</p>
            </div>
          </div>

          <div class="space-y-2">
            <div class="flex items-center gap-2">
              <Checkbox
                id="isFeatured"
                :model-value="form.isFeatured"
                @update:model-value="setFeatured"
              />
              <Label for="isFeatured">Featured project</Label>
            </div>
            <p v-if="form.errors.isFeatured" class="text-sm text-destructive">{{ form.errors.isFeatured }}</p>
          </div>
        </section>
      </div>

      <div class="flex justify-end gap-2 border-t px-6 py-4">
        <Link href="/admin/projects">
          <Button type="button" variant="outline">Cancel</Button>
        </Link>
        <Button type="submit" :disabled="form.processing">
          {{ form.processing ? 'Creating...' : 'Create Project' }}
        </Button>
      </div>
    </form>
  </div>
</template>
