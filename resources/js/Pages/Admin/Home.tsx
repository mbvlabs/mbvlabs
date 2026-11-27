import { routes } from '@/routes'
import { Card } from "@/Components/ui/card"
import { Button } from "@/Components/ui/button"
import { Badge } from "@/Components/ui/badge"
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "@/Components/ui/table"
import { Head, Link } from '@inertiajs/react'
import { ArrowUpRight, FileText, FolderKanban, Inbox, Newspaper } from 'lucide-react'
import AdminLayout from '@/Layouts/AdminLayout'
import type { AdminHomeProps, AdminInquiryData, AdminWorkData, AdminBlogData, AdminProjectData } from '@/types/payloads'

type Item = Pick<AdminBlogData, 'ID' | 'CreatedAt'> & Partial<
  Pick<AdminInquiryData, 'Name' | 'Email' | 'ProjectType' | 'Status'> &
  Pick<AdminWorkData, 'Title' | 'ClientName' | 'PublishedAt'> &
  Pick<AdminProjectData, 'IsFeatured'>
>

export default function Home({ projectInquiries, works, blogPosts, projects, canManageContentToken }: AdminHomeProps) {
  const sections = [
    { title: 'Latest Inquiries', card: 'Inquiries', icon: Inbox, description: 'Recent project requests.', href: routes.adminProjectInquiryIndex.path(), items: projectInquiries ?? [], columns: ['Name', 'Project', 'Status', 'Created'], values: (i: Item) => [<>{i.Name}<p className="text-xs text-muted-foreground">{i.Email}</p></>, i.ProjectType, status(i.Status), date(i.CreatedAt)], empty: 'No inquiries found.' },
    { title: 'Work', card: 'Work', icon: FileText, description: 'Recent work entries.', href: routes.adminWorkIndex.path(), items: works ?? [], columns: ['Title', 'Client', 'Status', 'Published'], values: (i: Item) => [i.Title, i.ClientName, status(i.Status), date(i.PublishedAt)], empty: 'No work found.' },
    { title: 'Blog Posts', card: 'Blog Posts', icon: Newspaper, description: 'Recent writing.', href: routes.adminBlogPostIndex.path(), items: blogPosts ?? [], columns: ['Title', 'Status', 'Published', 'Created'], values: (i: Item) => [i.Title, status(i.Status), date(i.PublishedAt), date(i.CreatedAt)], empty: 'No blog posts found.' },
    { title: 'Projects', card: 'Projects', icon: FolderKanban, description: 'Recent projects.', href: routes.adminProjectIndex.path(), items: projects ?? [], columns: ['Name', 'Type', 'Status', 'Published'], values: (i: Item) => [i.Name, i.ProjectType, status(i.IsFeatured ? 'Featured' : 'Standard', i.IsFeatured), date(i.PublishedAt)], empty: 'No projects found.' },
  ]
  return <AdminLayout><Head title="Admin" /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div><h1 className="text-2xl font-semibold tracking-tight">Admin</h1><p className="text-muted-foreground">Latest content and incoming project inquiries.</p></div>
    {canManageContentToken && <Card className="flex flex-row flex-wrap items-center justify-between gap-4 border p-5"><div><h2 className="font-semibold">Content API token</h2><p className="text-sm text-muted-foreground">Create or revoke a short-lived content API token.</p></div><Button asChild variant="outline"><Link href={routes.adminContentTokenShow.path()}>Manage token</Link></Button></Card>}
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">{sections.map(s => <Link key={s.href} href={s.href} className="border bg-card p-5 text-card-foreground shadow-sm transition-colors hover:bg-muted/40"><div className="flex items-center justify-between text-sm font-medium text-muted-foreground"><span>{s.card}</span><s.icon className="size-4" aria-hidden="true" /></div><div className="mt-3 text-2xl font-semibold tracking-tight">{s.items.length}</div><div className="mt-1 flex items-center text-xs text-muted-foreground"><ArrowUpRight className="mr-1 size-3" aria-hidden="true" />Latest {s.card === 'Inquiries' ? 'inquiries' : s.card === 'Work' ? 'entries' : s.card === 'Blog Posts' ? 'posts' : 'projects'}</div></Link>)}</div>
    <div className="grid gap-4 xl:grid-cols-2">{sections.map(s => <Card key={s.href} className="gap-0 border py-0 ring-0"><div className="flex items-center justify-between border-b p-5"><div><h2 className="text-base font-semibold leading-none tracking-tight">{s.title}</h2><p className="mt-1 text-sm text-muted-foreground">{s.description}</p></div><Button asChild variant="outline" size="sm" className="border-border"><Link href={s.href}>View all</Link></Button></div><Table className="w-full text-left text-sm"><TableHeader><TableRow>{s.columns.map(c => <TableHead key={c}>{c}</TableHead>)}</TableRow></TableHeader><TableBody>{s.items.map(i => <TableRow key={i.ID} className="hover:bg-muted/40">{s.values(i).map((v, n) => <TableCell key={n} className="px-3 py-3">{n === 0 ? <Link className="font-medium" href={`${s.href}/${i.ID}`}>{v}</Link> : <span className="text-muted-foreground">{v || '-'}</span>}</TableCell>)}</TableRow>)}{!s.items.length && <TableRow><TableCell colSpan={s.columns.length} className="h-20 text-center text-muted-foreground">{s.empty}</TableCell></TableRow>}</TableBody></Table></Card>)}</div>
  </div></AdminLayout>
}

function date(value?: string | null) { return value ? new Date(value).toLocaleDateString() : '-' }
function status(value?: string, featured = false) {
  return value ? <Badge variant={featured ? 'default' : 'secondary'}>{value}</Badge> : null
}
