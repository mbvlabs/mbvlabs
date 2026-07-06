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

const form = useForm({})

function destroy(id: number) {
  if (confirm('Are you sure you want to delete this work?')) {
    form.delete(`/admin/works/${id}`)
  }
}
</script>

<template>
  <Head title="Work" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
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
            <TableHead>Client</TableHead>
            <TableHead>Industry</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Featured</TableHead>
            <TableHead>Published</TableHead>
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="item in items" :key="item.ID">
            <TableCell class="font-medium">{{ item.Title }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.ClientName || '—' }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.ClientIndustry || '—' }}</TableCell>
            <TableCell>
              <Badge v-if="item.Status === 'published'" variant="default">Published</Badge>
              <Badge v-else variant="secondary">Draft</Badge>
            </TableCell>
            <TableCell>
              <Badge v-if="item.IsFeatured" variant="default">Featured</Badge>
              <Badge v-else variant="secondary">Not featured</Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">
              {{ item.PublishedAt ? new Date(item.PublishedAt).toLocaleDateString() : '—' }}
            </TableCell>
            <TableCell class="text-right">
              <div class="flex items-center justify-end gap-1">
                <Link :href="`/admin/works/${item.ID}`">
                  <Button variant="ghost" size="icon">
                    <Eye class="size-4" />
                    <span class="sr-only">View</span>
                  </Button>
                </Link>
                <Link :href="`/admin/works/${item.ID}/edit`">
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
            <TableCell colspan="7" class="h-24 text-center text-muted-foreground">
              No work found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
