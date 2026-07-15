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
  PublishedAt: string | null
  CreatedAt: string
  UpdatedAt: string
  PublicationSchedule: PublicationSchedule | null
}

defineProps<{
  items: BlogPost[]
}>()
</script>

<template>
  <Head title="Blog Posts" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Blog Posts</h1>
        <p class="text-muted-foreground">Manage your blog posts.</p>
      </div>
      <Link href="/admin/blog-posts/new">
        <Button>
          <Plus class="mr-2 size-4" />
          New Blog Post
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
            :href="routes.adminBlogPostShow(item.ID)"
            :aria-label="`View blog post ${item.Title}`"
          >
            <TableCell class="max-w-80 truncate font-medium">{{ item.Title }}</TableCell>
            <TableCell class="max-w-64 truncate text-muted-foreground">{{ item.Slug }}</TableCell>
            <TableCell>
              <Badge variant="secondary">
                {{ item.Status === 'draft' && item.PublicationSchedule ? 'Scheduled' : item.Status }}
              </Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">{{ new Date(item.CreatedAt).toLocaleDateString() }}</TableCell>
            <TableCell class="text-muted-foreground">{{ new Date(item.UpdatedAt).toLocaleDateString() }}</TableCell>
            <TableCell class="text-muted-foreground">
              {{ item.PublishedAt ? new Date(item.PublishedAt).toLocaleDateString() : '—' }}
            </TableCell>
          </TableRow>
          <TableRow v-if="items.length === 0">
            <TableCell colspan="6" class="h-24 text-center text-muted-foreground">
              No blog posts found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
