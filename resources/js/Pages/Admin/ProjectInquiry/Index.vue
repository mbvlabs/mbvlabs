<script setup lang="ts">
import { Head } from '@inertiajs/vue3'

import AdminLayout from '@/Layouts/AdminLayout.vue'
import { routes } from '@/routes'
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

defineProps<{
  items: ProjectInquiry[]
}>()
</script>

<template>
  <Head title="Project Inquiries" />

  <div class="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div>
      <h1 class="text-2xl font-semibold tracking-tight">Project Inquiries</h1>
      <p class="text-muted-foreground">Manage incoming project inquiries.</p>
    </div>

    <div class="rounded-none border bg-card">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>Email</TableHead>
            <TableHead>Project Type</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Created</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="item in items"
            :key="item.ID"
            :href="routes.adminProjectInquiryShow(item.ID)"
            :aria-label="`View project inquiry from ${item.Name}`"
          >
            <TableCell class="font-medium">{{ item.Name }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.Email }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.ProjectType || '—' }}</TableCell>
            <TableCell>
              <Badge variant="secondary">{{ item.Status }}</Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">
              {{ new Date(item.CreatedAt).toLocaleDateString() }}
            </TableCell>
          </TableRow>
          <TableRow v-if="items.length === 0">
            <TableCell colspan="5" class="h-24 text-center text-muted-foreground">
              No project inquiries found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
