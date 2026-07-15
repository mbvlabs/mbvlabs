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

defineProps<{
  items: Work[]
}>()
</script>

<template>
  <Head title="Work" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Work</h1>
        <p class="text-muted-foreground">Manage your work.</p>
      </div>
      <Link href="/admin/works/new">
        <Button>
          <Plus class="mr-2 size-4" />
          New Work
        </Button>
      </Link>
    </div>

    <div class="rounded-none border bg-card">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Title</TableHead>
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
            :href="routes.adminWorkShow(item.ID)"
            :aria-label="`View work ${item.Title}`"
          >
            <TableCell class="max-w-80 truncate font-medium">{{ item.Title }}</TableCell>
            <TableCell class="max-w-64 truncate text-muted-foreground">{{ item.Slug }}</TableCell>
            <TableCell>
              <Badge v-if="item.Status === 'published'" variant="default">Published</Badge>
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
              No work found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
