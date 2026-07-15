<script setup lang="ts">
import { Head, Link } from '@inertiajs/vue3'
import { Plus } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import { routes } from '@/routes'

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
</script>

<template>
  <Head title="Projects" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
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
            <TableHead>Slug</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Created</TableHead>
            <TableHead>Updated</TableHead>
            <TableHead>Published</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="item in items"
            :key="item.ID"
            :href="routes.adminProjectShow(item.ID)"
            :aria-label="`View project ${item.Name}`"
          >
            <TableCell class="max-w-80 truncate font-medium">{{ item.Name }}</TableCell>
            <TableCell class="max-w-64 truncate text-muted-foreground">{{ item.Slug }}</TableCell>
            <TableCell>
              <Badge v-if="item.PublishedAt" variant="default">Published</Badge>
              <Badge v-else variant="secondary">Draft</Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">{{ new Date(item.CreatedAt).toLocaleDateString() }}</TableCell>
            <TableCell class="text-muted-foreground">{{ new Date(item.UpdatedAt).toLocaleDateString() }}</TableCell>
            <TableCell class="text-muted-foreground">
              {{ item.PublishedAt ? new Date(item.PublishedAt).toLocaleDateString() : '—' }}
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
