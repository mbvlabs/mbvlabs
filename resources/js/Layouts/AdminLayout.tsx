import { Link, usePage } from '@inertiajs/react'
import { BookOpenText, FileText, FolderKanban, Home, Inbox, LayoutDashboard, ListChecks, Moon, Newspaper, ServerCog, Sun } from 'lucide-react'
import { useEffect, useState, type ReactNode } from 'react'
import FlashMessages from '@/Components/FlashMessages'
import { Button } from '@/Components/ui/button'
import { Separator } from '@/Components/ui/separator'
import {
  Sidebar, SidebarContent, SidebarGroup, SidebarGroupContent, SidebarGroupLabel,
  SidebarHeader, SidebarInset, SidebarMenu, SidebarMenuButton, SidebarMenuItem,
  SidebarProvider, SidebarRail, SidebarSeparator, SidebarTrigger, useSidebar,
} from '@/Components/ui/sidebar'
import { TooltipProvider } from '@/Components/ui/tooltip'

const platform = [
  ['Dashboard', '/admin', LayoutDashboard], ['Work', '/admin/works', FileText],
  ['Diary', '/admin/diary', BookOpenText], ['Blog Posts', '/admin/blog-posts', Newspaper],
  ['Projects', '/admin/projects', FolderKanban], ['Project Inquiries', '/admin/project-inquiries', Inbox],
] as const
const queue = [['Overview', '/admin/queue', ServerCog], ['Jobs', '/admin/queue/jobs', ListChecks]] as const
type NavItem = (typeof platform)[number] | (typeof queue)[number]

function AdminNavigation({ label, items }: { label: string; items: readonly NavItem[] }) {
  const { url } = usePage()
  const { setOpenMobile } = useSidebar()
  const path = url.split(/[?#]/)[0].replace(/\/$/, '') || '/'
  return (
    <SidebarGroup>
      <SidebarGroupLabel>{label}</SidebarGroupLabel>
      <SidebarGroupContent>
        <nav aria-label={label}>
          <SidebarMenu>
            {items.map(([title, href, Icon]) => {
              const active = path === href || (href !== '/admin' && href !== '/admin/queue' && path.startsWith(`${href}/`))
              return <SidebarMenuItem key={href}>
                <SidebarMenuButton asChild isActive={active} tooltip={title}>
                  <Link href={href} aria-current={active ? 'page' : undefined} onClick={() => setOpenMobile(false)}>
                    <Icon aria-hidden="true" />
                    <span>{title}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            })}
          </SidebarMenu>
        </nav>
      </SidebarGroupContent>
    </SidebarGroup>
  )
}

export default function AdminLayout({ children }: { children: ReactNode }) {
  const [dark, setDark] = useState(true)
  useEffect(() => {
    setDark(localStorage.getItem('theme') !== 'light')
  }, [])

  function toggleTheme() {
    localStorage.setItem('theme', dark ? 'light' : 'dark')
    setDark(!dark)
  }

  return (
    <TooltipProvider>
      <SidebarProvider className="admin-shell text-foreground" data-theme={dark ? 'dark' : 'light'}>
        <Sidebar variant="inset" collapsible="icon">
          <SidebarHeader>
            <SidebarMenu><SidebarMenuItem>
              <SidebarMenuButton asChild size="lg" tooltip="MBV Labs home">
                <a href="/" aria-label="MBV Labs home">
                  <span className="grid size-8 shrink-0 place-items-center bg-neutral-900 text-white"><Home className="size-4" /></span>
                  <span className="grid truncate text-sm leading-tight group-data-[collapsible=icon]:hidden">
                    <strong>mbvlabs. inc</strong><span className="text-xs">Admin</span>
                  </span>
                </a>
              </SidebarMenuButton>
            </SidebarMenuItem></SidebarMenu>
          </SidebarHeader>
          <SidebarContent>
            <AdminNavigation label="Platform" items={platform} />
            <SidebarSeparator />
            <AdminNavigation label="Queue" items={queue} />
          </SidebarContent>
          <SidebarRail />
        </Sidebar>
        <SidebarInset className="min-w-0">
          <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4 group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
            <div className="flex items-center gap-2 px-2">
              <SidebarTrigger className="-ml-1" />
              <Separator className="mr-2 data-[orientation=vertical]:h-4" orientation="vertical" />
              <Button type="button" variant="ghost" size="icon-sm" aria-label="Toggle theme" onClick={toggleTheme}>
                {dark ? <Moon className="size-4" /> : <Sun className="size-4" />}
              </Button>
            </div>
          </header>
          <div className="flex flex-1 flex-col gap-4 p-4">{children}</div>
        </SidebarInset>
        <FlashMessages />
      </SidebarProvider>
    </TooltipProvider>
  )
}
