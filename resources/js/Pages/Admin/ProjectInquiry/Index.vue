<script setup lang="ts">
import { Head, Link, useForm } from '@inertiajs/vue3'
import { Pencil, Trash2, Eye } from '@lucide/vue'

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

const form = useForm({})

function destroy(id: number) {
  if (confirm('Are you sure you want to delete this project inquiry?')) {
    form.delete(`/admin/project-inquiries/${id}`)
  }
}
</script>

<template>
  <Head title="Project Inquiries" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
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
            <TableHead class="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="item in items" :key="item.ID">
            <TableCell class="font-medium">{{ item.Name }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.Email }}</TableCell>
            <TableCell class="text-muted-foreground">{{ item.ProjectType || '—' }}</TableCell>
            <TableCell>
              <Badge variant="secondary">{{ item.Status }}</Badge>
            </TableCell>
            <TableCell class="text-muted-foreground">
              {{ new Date(item.CreatedAt).toLocaleDateString() }}
            </TableCell>
            <TableCell class="text-right">
              <div class="flex items-center justify-end gap-1">
                <Link :href="`/admin/project-inquiries/${item.ID}`">
                  <Button variant="ghost" size="icon">
                    <Eye class="size-4" />
                    <span class="sr-only">View</span>
                  </Button>
                </Link>
                <Link :href="`/admin/project-inquiries/${item.ID}/edit`">
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
              No project inquiries found.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>
