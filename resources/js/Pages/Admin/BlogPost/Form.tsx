import { Head, Link, useForm } from '@inertiajs/react'
import { useState } from 'react'
import AdminLayout from '@/Layouts/AdminLayout'
import RichEditor from '@/Components/RichEditor'
import { Button } from '@/Components/ui/button'
import { Input } from '@/Components/ui/input'
import { Textarea } from '@/Components/ui/textarea'
import { Label } from '@/Components/ui/label'
import { NativeSelect } from '@/Components/ui/native-select'

export type BlogPost = { ID: number; Title: string; Slug: string; Excerpt: string; Body: string; Status: string; CoverImageUrl: string | null; Tags: string | null; PublishedAt: string | null; CreatedAt: string; UpdatedAt: string; PublicationSchedule: { job_id: number; scheduled_at: string } | null }
const input = 'flex h-9 w-full border border-input bg-background px-3 py-1 text-sm shadow-sm outline-none focus-visible:ring-2 focus-visible:ring-ring'
const outline = 'inline-flex items-center justify-center border bg-background px-4 py-2 text-sm font-medium hover:bg-muted'
const primary = 'inline-flex items-center justify-center bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-50'

export default function Form({ item }: { item?: BlogPost }) {
  const editing = !!item
  const [slugEdited, setSlugEdited] = useState(false)
  const form = useForm({ title: item?.Title ?? '', slug: item?.Slug ?? '', excerpt: item?.Excerpt ?? '', body: item?.Body ?? '', status: item?.Status ?? 'draft', coverImageUrl: item?.CoverImageUrl ?? '', tags: item?.Tags ?? '', ...(editing ? {} : { scheduledAt: '' }) })
  const data = form.data as typeof form.data & { scheduledAt?: string }
  const cancel = editing ? `/admin/blog-posts/${item.ID}` : '/admin/blog-posts'
  const error = (field: keyof typeof form.data) => form.errors[field] && <p className="text-sm text-destructive">{form.errors[field]}</p>
  function submit(e: React.FormEvent) {
    e.preventDefault()
    if (editing) { form.put(cancel); return }
    form.clearErrors('scheduledAt')
    if (data.status === 'scheduled') {
      const date = new Date(data.scheduledAt || '')
      const message = !data.scheduledAt ? 'Publication date and time is required' : Number.isNaN(date.getTime()) ? 'Publication date and time is invalid' : date.getTime() <= Date.now() ? 'Publication date and time must be in the future' : ''
      if (message) { form.setError('scheduledAt', message); return }
    }
    form.transform(values => ({ ...values, scheduledAt: values.status === 'scheduled' ? new Date(values.scheduledAt || '').toISOString() : '' }))
    form.post('/admin/blog-posts')
  }
  function setTitle(title: string) {
    const slug = title.toLowerCase().trim().replace(/[^\w\s-]/g, '').replace(/[\s_]+/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '')
    form.setData({ ...form.data, title, ...(!editing || !slugEdited ? { slug } : {}) })
  }
  return <AdminLayout><Head title={editing ? `Edit ${item.Title}` : 'New Blog Post'} /><div className="mx-auto w-full min-w-0 max-w-7xl space-y-6">
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"><div><h1 className="text-2xl font-semibold tracking-tight">{editing ? 'Edit Blog Post' : 'New Blog Post'}</h1><p className="text-muted-foreground">{editing ? item.Title : 'Write a new blog post.'}</p></div><Button asChild variant="outline" className={`${outline} normal-case tracking-normal`}><Link href={cancel}>← {editing ? 'Cancel' : 'Back'}</Link></Button></div>
    <form onSubmit={submit} className="border bg-card"><div className="space-y-8 p-6">
      {editing && Object.keys(form.errors).length > 0 && <div className="border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">Please fix the highlighted fields.</div>}
      <section className="space-y-4"><h2 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Primary details</h2>
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
          <div className="space-y-1"><Label htmlFor="title" className="text-sm font-medium normal-case tracking-normal">Title</Label><Input id="title" className={input} value={data.title} onChange={e => setTitle(e.target.value)} />{error('title')}</div>
          <div className="space-y-1"><Label htmlFor="slug" className="text-sm font-medium normal-case tracking-normal">Slug</Label><Input id="slug" className={input} value={data.slug} onChange={e => { setSlugEdited(true); form.setData('slug', e.target.value) }} />{error('slug')}</div>
        </div>
        <div className="space-y-1"><Label htmlFor="excerpt" className="text-sm font-medium normal-case tracking-normal">Excerpt</Label><Textarea id="excerpt" className={`${input} h-20`} value={data.excerpt} onChange={e => form.setData('excerpt', e.target.value)} />{error('excerpt')}</div>
      </section>
      <section className="space-y-4"><h2 className="text-sm font-semibold uppercase tracking-wide text-muted-foreground">Publishing</h2>
        {item?.PublicationSchedule && <div className="border bg-muted/40 px-4 py-3 text-sm">Scheduled for {new Date(item.PublicationSchedule.scheduled_at).toLocaleString()}. Changing the status to Published or Archived clears this schedule.</div>}
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2"><div className="space-y-1"><Label htmlFor="status" className="text-sm font-medium normal-case tracking-normal">Status</Label><NativeSelect id="status" className="w-full [&_select]:h-9 [&_select]:border-input [&_select]:bg-background [&_select]:px-3" value={data.status} onChange={e => form.setData('status', e.target.value)}><option value="draft">Draft</option><option value="published">Published</option>{!editing && <option value="scheduled">Scheduled</option>}<option value="archived">Archived</option></NativeSelect>{error('status')}</div>
          {!editing && data.status === 'scheduled' && <div className="space-y-1"><Label htmlFor="scheduledAt" className="text-sm font-medium normal-case tracking-normal">Publish at</Label><Input id="scheduledAt" type="datetime-local" className={input} value={data.scheduledAt ?? ''} onChange={e => form.setData('scheduledAt', e.target.value)} />{error('scheduledAt')}</div>}
          <div className="space-y-1 md:col-span-2"><Label htmlFor="coverImageUrl" className="text-sm font-medium normal-case tracking-normal">Cover Image URL</Label><Input id="coverImageUrl" className={input} value={data.coverImageUrl} onChange={e => form.setData('coverImageUrl', e.target.value)} />{error('coverImageUrl')}</div>
          <div className="space-y-1 md:col-span-2"><Label htmlFor="tags" className="text-sm font-medium normal-case tracking-normal">Tags</Label><Input id="tags" className={input} placeholder={editing ? undefined : 'strategy, systems, launch'} value={data.tags} onChange={e => form.setData('tags', e.target.value)} />{error('tags')}</div>
        </div>
      </section>
      <section className="space-y-2"><Label htmlFor="body" className="text-sm font-medium normal-case tracking-normal">Body</Label><RichEditor id="body" value={data.body} onChange={value => form.setData('body', value)} placeholder="Write your blog post..." className="min-h-[32rem]" />{error('body')}</section>
    </div><div className="flex justify-end gap-2 border-t px-6 py-4"><Button asChild variant="outline" className={`${outline} normal-case tracking-normal`}><Link href={cancel}>Cancel</Link></Button><Button type="submit" className={`${primary} normal-case tracking-normal`} disabled={form.processing}>{form.processing ? editing ? 'Updating...' : 'Creating...' : editing ? 'Update Blog Post' : 'Create Blog Post'}</Button></div></form>
  </div></AdminLayout>
}
