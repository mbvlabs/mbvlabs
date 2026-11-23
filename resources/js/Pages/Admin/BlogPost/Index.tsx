import { Head, Link, router } from '@inertiajs/react'
import AdminLayout from '@/Layouts/AdminLayout'
import type { BlogPost } from './Form'
import { Badge } from '@/Components/ui/badge'
import { Button } from '@/Components/ui/button'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/Components/ui/table'

export default function Index({ items }: { items: BlogPost[] }) {
  return <AdminLayout><Head title="Blog Posts" /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div className="flex items-center justify-between"><div><h1 className="text-2xl font-semibold tracking-tight">Blog Posts</h1><p className="text-muted-foreground">Manage your blog posts.</p></div><Button asChild className="px-4 py-2 text-sm font-medium normal-case tracking-normal"><Link href="/admin/blog-posts/new">＋ New Blog Post</Link></Button></div>
    <div className="overflow-x-auto border bg-card"><Table className="w-full text-left text-sm"><TableHeader><TableRow className="border-b">{['Title', 'Slug', 'Status', 'Created', 'Updated', 'Published'].map(x => <TableHead key={x} className="h-10 px-4 text-sm font-medium normal-case tracking-normal text-muted-foreground">{x}</TableHead>)}</TableRow></TableHeader><TableBody>
      {items.map(item => <TableRow key={item.ID} role="link" tabIndex={0} aria-label={`View blog post ${item.Title}`} onClick={() => router.visit(`/admin/blog-posts/${item.ID}`)} onKeyDown={e => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); router.visit(`/admin/blog-posts/${item.ID}`) } }} className="cursor-pointer border-b hover:bg-muted/40 focus-visible:outline-2 focus-visible:outline-ring"><TableCell className="max-w-80 truncate p-4 font-medium">{item.Title}</TableCell><TableCell className="max-w-64 truncate p-4 text-muted-foreground">{item.Slug}</TableCell><TableCell className="p-4"><Badge className="bg-secondary px-2 py-1 text-xs normal-case tracking-normal">{item.Status === 'draft' && item.PublicationSchedule ? 'Scheduled' : item.Status}</Badge></TableCell><TableCell className="p-4 text-muted-foreground">{new Date(item.CreatedAt).toLocaleDateString()}</TableCell><TableCell className="p-4 text-muted-foreground">{new Date(item.UpdatedAt).toLocaleDateString()}</TableCell><TableCell className="p-4 text-muted-foreground">{item.PublishedAt ? new Date(item.PublishedAt).toLocaleDateString() : '—'}</TableCell></TableRow>)}
      {!items.length && <TableRow><TableCell colSpan={6} className="h-24 text-center text-muted-foreground">No blog posts found.</TableCell></TableRow>}
    </TableBody></Table></div>
  </div></AdminLayout>
}
