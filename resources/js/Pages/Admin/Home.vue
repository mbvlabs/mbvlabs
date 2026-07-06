<script setup lang="ts">
import { Head, Link } from '@inertiajs/vue3'
import { ArrowUpRight, FileText, FolderKanban, Inbox, Newspaper } from '@lucide/vue'

import AdminLayout from '@/Layouts/AdminLayout.vue'

import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
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
  ProjectType: string | null
  Status: string
  CreatedAt: string
}

interface Work {
  ID: number
  Title: string
  ClientName: string | null
  Status: string
  PublishedAt: string | null
  CreatedAt: string
}

interface BlogPost {
  ID: number
  Title: string
  Status: string
  PublishedAt: string | null
  CreatedAt: string
}

interface Project {
  ID: number
  Name: string
  ProjectType: string | null
  PublishedAt: string | null
  IsFeatured: boolean
  CreatedAt: string
}

defineProps<{
  projectInquiries: ProjectInquiry[]
  works: Work[]
  blogPosts: BlogPost[]
  projects: Project[]
}>()

function formatDate(value: string | null): string {
  if (!value) {
    return '-'
  }

  return new Date(value).toLocaleDateString()
}
</script>

<template>
  <Head title="Admin" />

  <div class="mx-auto w-full min-w-0 max-w-6xl space-y-6">
    <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight">Admin</h1>
        <p class="text-muted-foreground">Latest content and incoming project inquiries.</p>
      </div>
    </div>

    <div class="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <Link
        href="/admin/project-inquiries"
        class="border bg-card p-5 text-card-foreground shadow-sm transition-colors hover:bg-muted/40"
      >
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-muted-foreground">Inquiries</span>
          <Inbox class="size-4 text-muted-foreground" />
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ projectInquiries.length }}</div>
        <div class="mt-1 flex items-center text-xs text-muted-foreground">
          <ArrowUpRight class="mr-1 size-3" />
          Latest inquiries
        </div>
      </Link>

      <Link
        href="/admin/works"
        class="border bg-card p-5 text-card-foreground shadow-sm transition-colors hover:bg-muted/40"
      >
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-muted-foreground">Work</span>
          <FileText class="size-4 text-muted-foreground" />
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ works.length }}</div>
        <div class="mt-1 flex items-center text-xs text-muted-foreground">
          <ArrowUpRight class="mr-1 size-3" />
          Latest entries
        </div>
      </Link>

      <Link
        href="/admin/blog-posts"
        class="border bg-card p-5 text-card-foreground shadow-sm transition-colors hover:bg-muted/40"
      >
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-muted-foreground">Blog Posts</span>
          <Newspaper class="size-4 text-muted-foreground" />
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ blogPosts.length }}</div>
        <div class="mt-1 flex items-center text-xs text-muted-foreground">
          <ArrowUpRight class="mr-1 size-3" />
          Latest posts
        </div>
      </Link>

      <Link
        href="/admin/projects"
        class="border bg-card p-5 text-card-foreground shadow-sm transition-colors hover:bg-muted/40"
      >
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-muted-foreground">Projects</span>
          <FolderKanban class="size-4 text-muted-foreground" />
        </div>
        <div class="mt-3 text-2xl font-semibold tracking-tight">{{ projects.length }}</div>
        <div class="mt-1 flex items-center text-xs text-muted-foreground">
          <ArrowUpRight class="mr-1 size-3" />
          Latest projects
        </div>
      </Link>
    </div>

    <div class="grid gap-4 xl:grid-cols-2">
      <section class="border bg-card">
        <div class="flex items-center justify-between border-b p-5">
          <div>
            <h2 class="font-semibold leading-none tracking-tight">Latest Inquiries</h2>
            <p class="mt-1 text-sm text-muted-foreground">Recent project requests.</p>
          </div>
          <Link href="/admin/project-inquiries">
            <Button variant="outline" size="sm">View all</Button>
          </Link>
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Project</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="item in projectInquiries" :key="item.ID">
              <TableCell>
                <Link :href="`/admin/project-inquiries/${item.ID}`" class="font-medium hover:underline">
                  {{ item.Name }}
                </Link>
                <p class="text-xs text-muted-foreground">{{ item.Email }}</p>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ item.ProjectType || '-' }}</TableCell>
              <TableCell><Badge variant="secondary">{{ item.Status }}</Badge></TableCell>
              <TableCell class="text-muted-foreground">{{ formatDate(item.CreatedAt) }}</TableCell>
            </TableRow>
            <TableRow v-if="projectInquiries.length === 0">
              <TableCell colspan="4" class="h-20 text-center text-muted-foreground">
                No inquiries found.
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </section>

      <section class="border bg-card">
        <div class="flex items-center justify-between border-b p-5">
          <div>
            <h2 class="font-semibold leading-none tracking-tight">Work</h2>
            <p class="mt-1 text-sm text-muted-foreground">Recent work entries.</p>
          </div>
          <Link href="/admin/works">
            <Button variant="outline" size="sm">View all</Button>
          </Link>
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Title</TableHead>
              <TableHead>Client</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Published</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="item in works" :key="item.ID">
              <TableCell>
                <Link :href="`/admin/works/${item.ID}`" class="font-medium hover:underline">
                  {{ item.Title }}
                </Link>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ item.ClientName || '-' }}</TableCell>
              <TableCell><Badge variant="secondary">{{ item.Status }}</Badge></TableCell>
              <TableCell class="text-muted-foreground">{{ formatDate(item.PublishedAt) }}</TableCell>
            </TableRow>
            <TableRow v-if="works.length === 0">
              <TableCell colspan="4" class="h-20 text-center text-muted-foreground">
                No work found.
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </section>

      <section class="border bg-card">
        <div class="flex items-center justify-between border-b p-5">
          <div>
            <h2 class="font-semibold leading-none tracking-tight">Blog Posts</h2>
            <p class="mt-1 text-sm text-muted-foreground">Recent writing.</p>
          </div>
          <Link href="/admin/blog-posts">
            <Button variant="outline" size="sm">View all</Button>
          </Link>
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Title</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Published</TableHead>
              <TableHead>Created</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="item in blogPosts" :key="item.ID">
              <TableCell>
                <Link :href="`/admin/blog-posts/${item.ID}`" class="font-medium hover:underline">
                  {{ item.Title }}
                </Link>
              </TableCell>
              <TableCell><Badge variant="secondary">{{ item.Status }}</Badge></TableCell>
              <TableCell class="text-muted-foreground">{{ formatDate(item.PublishedAt) }}</TableCell>
              <TableCell class="text-muted-foreground">{{ formatDate(item.CreatedAt) }}</TableCell>
            </TableRow>
            <TableRow v-if="blogPosts.length === 0">
              <TableCell colspan="4" class="h-20 text-center text-muted-foreground">
                No blog posts found.
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </section>

      <section class="border bg-card">
        <div class="flex items-center justify-between border-b p-5">
          <div>
            <h2 class="font-semibold leading-none tracking-tight">Projects</h2>
            <p class="mt-1 text-sm text-muted-foreground">Recent projects.</p>
          </div>
          <Link href="/admin/projects">
            <Button variant="outline" size="sm">View all</Button>
          </Link>
        </div>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Published</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="item in projects" :key="item.ID">
              <TableCell>
                <Link :href="`/admin/projects/${item.ID}`" class="font-medium hover:underline">
                  {{ item.Name }}
                </Link>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ item.ProjectType || '-' }}</TableCell>
              <TableCell>
                <Badge v-if="item.IsFeatured" variant="default">Featured</Badge>
                <Badge v-else variant="secondary">Standard</Badge>
              </TableCell>
              <TableCell class="text-muted-foreground">{{ formatDate(item.PublishedAt) }}</TableCell>
            </TableRow>
            <TableRow v-if="projects.length === 0">
              <TableCell colspan="4" class="h-20 text-center text-muted-foreground">
                No projects found.
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </section>
    </div>
  </div>
</template>
