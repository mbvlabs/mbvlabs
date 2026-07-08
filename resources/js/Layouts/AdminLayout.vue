<script setup lang="ts">
import { Link, usePage } from '@inertiajs/vue3'
import { computed } from 'vue'
import {
  Home,
  LayoutDashboard,
  FileText,
  Newspaper,
  FolderKanban,
  Inbox,
  BookOpenText,
  ListChecks,
  ServerCog,
} from '@lucide/vue'

import ThemeToggle from '@/Components/ThemeToggle.vue'
import { Separator } from '@/components/ui/separator'
import { routes } from '@/routes'
import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarRail,
  SidebarSeparator,
  SidebarTrigger,
} from '@/components/ui/sidebar'

interface NavItem {
  title: string
  href: string
  icon: typeof Home
  exact?: boolean
}

const mainNav: NavItem[] = [
  { title: 'Dashboard', href: routes.adminHomePage(), icon: LayoutDashboard },
  { title: 'Work', href: routes.adminWorkIndex(), icon: FileText },
  { title: 'Diary', href: routes.adminDiaryEntryIndex(), icon: BookOpenText },
  { title: 'Blog Posts', href: routes.adminBlogPostIndex(), icon: Newspaper },
  { title: 'Projects', href: routes.adminProjectIndex(), icon: FolderKanban },
  { title: 'Project Inquiries', href: routes.adminProjectInquiryIndex(), icon: Inbox },
]

const queueNav: NavItem[] = [
  { title: 'Overview', href: routes.adminQueueIndex(), icon: ServerCog, exact: true },
  { title: 'Jobs', href: routes.adminQueueJobs(), icon: ListChecks },
]

const page = usePage()

const currentPath = computed(() => {
  const path = page.url.split(/[?#]/)[0]
  return path.replace(/\/$/, '') || '/'
})

function isActiveNav(item: NavItem): boolean {
  if (item.href === '#') {
    return false
  }

  if (item.exact || item.href === '/admin') {
    return currentPath.value === item.href
  }

  return currentPath.value === item.href || currentPath.value.startsWith(`${item.href}/`)
}
</script>

<template>
  <SidebarProvider>
    <Sidebar variant="inset" collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              size="lg"
              class="group-data-[collapsible=icon]:p-0!"
              as-child
            >
              <Link href="/">
                <div
                  class="flex aspect-square size-8 items-center justify-center rounded-none bg-neutral-900 text-white"
                >
                  <Home class="size-4" />
                </div>
                <div class="grid flex-1 text-left text-sm leading-tight">
                  <span class="truncate font-semibold">mbvlabs. inc</span>
                  <span class="truncate text-xs">Admin</span>
                </div>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Platform</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in mainNav" :key="item.title">
                <SidebarMenuButton as-child :is-active="isActiveNav(item)" :tooltip="item.title">
                  <Link :href="item.href">
                    <component :is="item.icon" />
                    <span>{{ item.title }}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        <SidebarSeparator />

        <SidebarGroup>
          <SidebarGroupLabel>Queue</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem v-for="item in queueNav" :key="item.title">
                <SidebarMenuButton as-child :is-active="isActiveNav(item)" :tooltip="item.title">
                  <Link :href="item.href">
                    <component :is="item.icon" />
                    <span>{{ item.title }}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

      </SidebarContent>
      <SidebarRail />
    </Sidebar>

    <SidebarInset>
      <header
        class="flex h-16 shrink-0 items-center gap-2 border-b px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12"
      >
        <div class="flex items-center gap-2 px-2">
          <SidebarTrigger class="-ml-1" />
          <Separator orientation="vertical" class="mr-2 h-4" />
          <div class="ml-auto">
            <ThemeToggle />
          </div>
        </div>
      </header>
      <div class="flex flex-1 flex-col gap-4 p-4">
        <slot />
      </div>
    </SidebarInset>
  </SidebarProvider>
</template>
