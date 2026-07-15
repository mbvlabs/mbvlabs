<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { ArrowLeft } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select } from '@/components/ui/select'
import { Textarea } from '@/components/ui/textarea'

defineOptions({ layout: AdminLayout })

const projectTypeOptions = [
  { label: 'Web app', value: 'web-app' },
  { label: 'Website', value: 'website' },
  { label: 'Automation', value: 'automation' },
  { label: 'Integration', value: 'integration' },
  { label: 'Consulting', value: 'consulting' },
  { label: 'Open source', value: 'open-source' },
  { label: 'Modernization', value: 'modernization' },
]

const statusOptions = [
  { label: 'New', value: 'new' },
  { label: 'Contacted', value: 'contacted' },
  { label: 'Qualified', value: 'qualified' },
  { label: 'Proposal', value: 'proposal' },
  { label: 'Won', value: 'won' },
  { label: 'Lost', value: 'lost' },
]

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

const form = useForm({
  name: props.item.Name,
  email: props.item.Email,
  company: props.item.Company ?? '',
  role: props.item.Role ?? '',
  projectType: props.item.ProjectType ?? 'web-app',
  timeline: props.item.Timeline ?? '',
  message: props.item.Message,
  source: props.item.Source ?? '',
  status: props.item.Status,
  metadata: props.item.Metadata ?? '{}',
})

function submit() {
  form.put(`/admin/project-inquiries/${props.item.ID}`)
}
</script>

<template>
  <Head :title="`Edit ${item.Name}`" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Edit Project Inquiry</h1>
        <p class="text-muted-foreground">{{ item.Name }}</p>
      </div>
      <Link :href="`/admin/project-inquiries/${item.ID}`">
        <Button variant="outline">
          <ArrowLeft class="mr-2 size-4" />
          Cancel
        </Button>
      </Link>
    </div>

    <form @submit.prevent="submit" class="border bg-card">
      <div class="space-y-8 p-6">
        <section class="space-y-4">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Contact</h2>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-1">
              <Label for="name">Name</Label>
              <Input id="name" v-model="form.name" />
            </div>

            <div class="space-y-1">
              <Label for="email">Email</Label>
              <Input id="email" type="email" v-model="form.email" />
            </div>

            <div class="space-y-1">
              <Label for="company">Company</Label>
              <Input id="company" v-model="form.company" />
            </div>

            <div class="space-y-1">
              <Label for="role">Role</Label>
              <Input id="role" v-model="form.role" />
            </div>
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Project</h2>
          <div class="grid grid-cols-1 gap-6 md:grid-cols-2">
            <div class="space-y-1">
              <Label for="projectType">Project Type</Label>
              <Select id="projectType" v-model="form.projectType" :options="projectTypeOptions" />
            </div>

            <div class="space-y-1">
              <Label for="status">Status</Label>
              <Select id="status" v-model="form.status" :options="statusOptions" />
            </div>

            <div class="space-y-1">
              <Label for="timeline">Timeline</Label>
              <Input id="timeline" v-model="form.timeline" />
            </div>

            <div class="space-y-1 md:col-span-2">
              <Label for="source">Source</Label>
              <Input id="source" v-model="form.source" />
            </div>
          </div>

          <div class="space-y-1">
            <Label for="message">Message</Label>
            <Textarea id="message" v-model="form.message" class="min-h-32" />
          </div>
        </section>

        <section class="space-y-2">
          <Label for="metadata">Metadata</Label>
          <Textarea id="metadata" v-model="form.metadata" class="min-h-40 font-mono text-sm" />
        </section>
      </div>

      <div class="flex justify-end gap-2 border-t px-6 py-4">
        <Link :href="`/admin/project-inquiries/${item.ID}`">
          <Button type="button" variant="outline">Cancel</Button>
        </Link>
        <Button type="submit" :disabled="form.processing">
          {{ form.processing ? 'Updating...' : 'Update Project Inquiry' }}
        </Button>
      </div>
    </form>
  </div>
</template>
