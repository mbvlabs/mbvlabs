<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { Plus, Pencil, Trash2, Eye } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'

import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

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

defineProps<{
  items: Project[]
}>()

const form = useForm({})

function destroy(id: number) {
  if (confirm('Are you sure you want to delete this project?')) {
    form.delete(`/admin/projects/${id}`)
  }
}
</script>

<template>
  <Head title="Projects" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Projects</h1>
        <p class="text-muted-foreground">Manage your projects.</p>
      </div>
      <Link href="/admin/projects/new">
        <Button>
          <Plus class="mr-2 size-4" />
          New Project
        </Button>
      </Link>
    </div>

    <div class="rounded-none border bg-card">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Tagline</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Published</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="item in items" :key="item.ID">
            <TableCell class="font-medium">{{ item.Name }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.Tagline || '—' }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.ProjectType || '—' }}</TableCell>
            <TableCell>
              <Badge v-if="item.IsFeatured" variant="default">Featured</Badge>
              <Badge v-else variant="secondary">Standard</Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">
              {{ item.PublishedAt ? new Date(item.PublishedAt).toLocaleDateString() : '—' }}
            </TableCell>
            <TableCell class="text-right">
              <div class="flex items-center justify-end gap-1">
                <Link :href="`/admin/projects/${item.ID}`">
                  <Button variant="ghost" size="icon">
                    <Eye class="size-4" />
                    <span class="sr-only">View</span>
                  </Button>
                </Link>
                <Link :href="`/admin/projects/${item.ID}/edit`">
                  <Button variant="ghost" size="icon">
                    <Pencil class="size-4" />
                    <span class="sr-only">Edit</span>
                  </Button>
                </Link>
                <Button variant="ghost" size="icon" @click="destroy(item.ID)">
                  <Trash2 class="size-4 text-destructive" />
                  <span class="sr-only">Delete</span>
                </Button>
              </div>
            </TableCell>
          </TableRow>
          <TableRow v-if="items.length === 0">
            <TableCell colspan="6" class="h-24 text-center text-muted-foreground">
              No projects found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
